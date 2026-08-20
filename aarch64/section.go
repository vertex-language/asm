package aarch64

import (
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/vertex-language/asm/aarch64/internal/encode"
	"github.com/vertex-language/asm/aarch64/operand"
)

// Section is an append-only byte buffer plus the three facts bytes cannot
// carry. Nothing here is symbolic-and-rewritable.
type Section struct {
	m    *Module
	kind SectionKind
	buf  []byte

	labels  map[string]int // every label in this section, bare or promoted
	symbols []Symbol
	fixups  []encode.Fixup // pending until Finalize
	refs    []Reference    // what survives Finalize
}

func (s *Section) Kind() SectionKind  { return s.kind }
func (s *Section) Name() string       { return s.kind.String() }
func (s *Section) Bytes() []byte      { return s.buf }
func (s *Section) Refs() []Reference  { return s.refs }
func (s *Section) Symbols() []Symbol  { return s.symbols }

// ready is the sticky-error and finalized guard every builder call runs.
func (s *Section) ready(context string) bool {
	if s.m.err != nil {
		return false
	}
	if s.m.finalized {
		s.failAt(len(s.buf), ErrFinalized, nil, context)
		return false
	}
	return true
}

func (s *Section) fail(sentinel, cause error, context string, notes ...string) {
	s.failAt(len(s.buf), sentinel, cause, context, notes...)
}

func (s *Section) failAt(off int, sentinel, cause error, context string, notes ...string) {
	s.m.setErr(&Error{
		Section: s.Name(), Offset: off, Context: context, Notes: notes,
		sentinel: sentinel, cause: cause,
	})
}

// Label defines a name at the current offset. Bare labels live only in this
// section's namespace; any attribute promotes the label into Symbols(), whose
// namespace is module-global.
func (s *Section) Label(name string, attrs ...SymbolAttr) {
	if !s.ready("label " + name) {
		return
	}
	if _, dup := s.labels[name]; dup {
		s.fail(ErrDuplicate, nil, "label "+name, "defined twice in "+s.Name())
		return
	}
	s.labels[name] = len(s.buf)
	if len(attrs) == 0 {
		return
	}
	if s.m.symbols[name] != nil {
		s.fail(ErrDuplicate, nil, "symbol "+name, "already defined in "+s.m.symbols[name].Name())
		return
	}
	sym := Symbol{Name: name, Offset: len(s.buf)}
	for _, a := range attrs {
		a.apply(&sym)
	}
	s.m.symbols[name] = s
	s.symbols = append(s.symbols, sym)
}

// ---- Data ----

func (s *Section) Byte(v uint8) {
	if s.ready("byte") {
		s.buf = append(s.buf, v)
	}
}

func (s *Section) Long(v uint32) {
	if s.ready("long") {
		var b [4]byte
		binary.LittleEndian.PutUint32(b[:], v)
		s.buf = append(s.buf, b[:]...)
	}
}

func (s *Section) Quad(v uint64) {
	if s.ready("quad") {
		var b [8]byte
		binary.LittleEndian.PutUint64(b[:], v)
		s.buf = append(s.buf, b[:]...)
	}
}

func (s *Section) Ascii(str string) {
	if s.ready("ascii") {
		s.buf = append(s.buf, str...)
	}
}

func (s *Section) Asciz(str string) {
	if s.ready("asciz") {
		s.buf = append(s.buf, str...)
		s.buf = append(s.buf, 0)
	}
}

func (s *Section) Zero(n int) {
	if s.ready("zero") {
		s.buf = append(s.buf, make([]byte, n)...)
	}
}

func (s *Section) Data(b []byte) {
	if s.ready("data") {
		s.buf = append(s.buf, b...)
	}
}

// Align pads a code section with d503201f and a data section with zeros.
// n must be a power of two, and in a code section a multiple of four; a text
// alignment that would strand a partial instruction is ErrAlign, not rounded.
func (s *Section) Align(n int) {
	if !s.ready("align") {
		return
	}
	if n <= 0 || n&(n-1) != 0 {
		s.fail(ErrAlign, nil, fmt.Sprintf("align %d", n), "alignment must be a power of two")
		return
	}
	pad := -len(s.buf) & (n - 1)
	if pad == 0 {
		return
	}
	if s.kind == Text {
		if n%4 != 0 || pad%4 != 0 {
			s.fail(ErrAlign, nil, fmt.Sprintf("align %d", n),
				"padding of a code section must be whole instructions; this alignment strands a partial word")
			return
		}
		s.buf = encode.Pad(s.buf, pad)
		return
	}
	s.buf = append(s.buf, make([]byte, pad)...)
}

// ---- Finalize machinery ----

// resolve runs at Finalize: label targets fold into the bytes, symbol targets
// become References, and everything else is refused by name.
func (s *Section) resolve() {
	for _, fx := range s.fixups {
		switch t := fx.Target.(type) {
		case operand.Label:
			s.patchLabel(fx, string(t))
		case operand.SymRef:
			if s.m.symbols[t.Name] == nil && !s.m.externs[t.Name] {
				s.failAt(fx.Offset, ErrUndefined, nil, "reference to "+t.Name,
					"define the symbol in this module or declare it with Extern")
				return
			}
			s.refs = append(s.refs, Reference{
				Offset: fx.Offset, Sym: t.Name, Role: fx.Role, Kind: fx.Kind,
				Addend: fx.Addend, Access: fx.Access, Branch: fx.Branch,
			})
		default:
			s.failAt(fx.Offset, ErrUndefined, nil, "reference with no target")
		}
		if s.m.err != nil {
			return
		}
	}
	s.fixups = nil
}

// patchLabel folds a same-section, direct, pc-relative reference into the
// word it belongs to. Only RoleDirect folds: the page of an address depends
// on where the section finally loads, which nothing at this layer assigns.
func (s *Section) patchLabel(fx encode.Fixup, name string) {
	target, ok := s.labels[name]
	if !ok {
		s.failAt(fx.Offset, ErrUndefined, nil, "label "+name,
			"labels are per-section; a cross-section or external target needs a symbol and a Ref")
		return
	}
	if fx.Role != operand.RoleDirect {
		s.failAt(fx.Offset, ErrUndefined, nil, "label "+name,
			"a "+fx.Role.String()+" reference needs a relocation, and a relocation needs a symbol; promote the label with an attribute")
		return
	}
	delta := int64(target) - int64(fx.Offset) + fx.Addend
	if fx.Scale > 0 && delta%(1<<fx.Scale) != 0 {
		s.failAt(fx.Offset, ErrRange, nil, "branch to "+name,
			fmt.Sprintf("displacement %d is not a multiple of %d", delta, 1<<fx.Scale))
		return
	}
	v := delta >> fx.Scale
	if !fitsSigned(v, fx.Bits) {
		s.failAt(fx.Offset, ErrRange, nil, "branch to "+name,
			fmt.Sprintf("displacement %d does not fit %d signed bits", delta, fx.Bits))
		return
	}
	word := binary.LittleEndian.Uint32(s.buf[fx.Offset:])
	word = fx.Field.Put(word, uint64(v)&maskBits(fx.Bits))
	binary.LittleEndian.PutUint32(s.buf[fx.Offset:], word)
}

// closeSizes closes every symbol's size at the next symbol or section end.
func (s *Section) closeSizes() {
	if len(s.symbols) == 0 {
		return
	}
	order := make([]int, len(s.symbols))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(a, b int) bool {
		return s.symbols[order[a]].Offset < s.symbols[order[b]].Offset
	})
	for i, idx := range order {
		end := len(s.buf)
		if i+1 < len(order) {
			end = s.symbols[order[i+1]].Offset
		}
		s.symbols[idx].Size = end - s.symbols[idx].Offset
	}
}

func fitsSigned(v int64, bits uint8) bool {
	if bits == 0 || bits > 63 {
		return false
	}
	lo := -(int64(1) << (bits - 1))
	hi := (int64(1) << (bits - 1)) - 1
	return v >= lo && v <= hi
}

func maskBits(bits uint8) uint64 {
	if bits >= 64 {
		return ^uint64(0)
	}
	return (uint64(1) << bits) - 1
}