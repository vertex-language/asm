// x86_64/section.go
package x86_64

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/vertex-language/asm/x86_64/internal/encode"
	"github.com/vertex-language/asm/x86_64/internal/isa"
	"github.com/vertex-language/asm/x86_64/operand"
)

// Section is an append-only byte buffer plus the three facts bytes cannot
// carry. Nothing is symbolic-and-rewritable: an instruction's bytes are
// decided at the call and never revisited.
type Section struct {
	m    *Module
	kind SectionKind

	buf    []byte
	labels map[string]int // every label, bare or promoted; section namespace
	syms   []Symbol       // promoted labels, in definition order
	pend   []pending      // fixups awaiting Finalize
	refs   []Reference    // what survives Finalize
}

// pending is one fixup between emission and Finalize: enough to either patch
// it (same-section label, PC-relative) or turn it into a Reference.
type pending struct {
	off, size, tail int
	pcrel           bool
	use             encode.Use
	kind            RefKind
	sym             string
	label           bool // target was an operand.Label, not a SymRef
	addend          int64
	context         string
}

func (s *Section) Kind() SectionKind { return s.kind }
func (s *Section) Name() string      { return s.kind.String() }

// Bytes is the finished machine code. After Finalize every same-section
// label displacement is patched in; before it, holes are zero.
func (s *Section) Bytes() []byte { return s.buf }

// Refs is the set every downstream format must turn into relocation
// records. Empty until Finalize.
func (s *Section) Refs() []Reference { return s.refs }

// Symbols is every promoted label, sizes closed at Finalize.
func (s *Section) Symbols() []Symbol { return s.syms }

// Offset is the current end of the section: the offset the next byte will
// land at, the value a Label placed now would name. It is exported because
// jump tables and literal pools are yours to build — this package
// deliberately refuses to build them — and building them requires knowing
// where you are.
func (s *Section) Offset() int { return len(s.buf) }

// ---- sticky-error plumbing --------------------------------------------

func (s *Section) fail(context string, sentinel error, notes ...string) {
	s.m.fail(&Error{
		Section: s.Name(), Offset: len(s.buf),
		Context: context, Err: sentinel, Notes: notes,
	})
}

func (s *Section) blocked() bool {
	if s.m.err != nil {
		return true
	}
	if s.m.done {
		s.fail("builder call after Finalize", ErrFinalized)
		return true
	}
	return false
}

// ---- labels -------------------------------------------------------------

// Label defines a label at the current offset. A bare label lives only in
// this section's namespace; any attribute promotes it into Symbols(), where
// names are module-global and a duplicate is ErrDuplicate.
func (s *Section) Label(name string, attrs ...SymbolAttr) {
	if s.blocked() {
		return
	}
	if _, dup := s.labels[name]; dup {
		s.fail(fmt.Sprintf("label %q defined twice in %s", name, s.Name()), ErrDuplicate)
		return
	}
	off := len(s.buf)
	s.labels[name] = off
	if len(attrs) == 0 {
		return
	}
	if owner := s.m.symOf[name]; owner != nil {
		s.fail(fmt.Sprintf("symbol %q already defined in %s", name, owner.Name()), ErrDuplicate)
		return
	}
	sym := Symbol{Name: name, Offset: off}
	for _, a := range attrs {
		a.applyTo(&sym)
	}
	s.syms = append(s.syms, sym)
	s.m.symOf[name] = s
}

// ---- data ---------------------------------------------------------------

func (s *Section) Byte(b byte) {
	if s.blocked() {
		return
	}
	s.buf = append(s.buf, b)
}

// Long is four bytes, little-endian.
func (s *Section) Long(v uint32) {
	if s.blocked() {
		return
	}
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], v)
	s.buf = append(s.buf, b[:]...)
}

// Quad is eight bytes, little-endian.
func (s *Section) Quad(v uint64) {
	if s.blocked() {
		return
	}
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], v)
	s.buf = append(s.buf, b[:]...)
}

// Ascii is the string's bytes, no terminator.
func (s *Section) Ascii(str string) {
	if s.blocked() {
		return
	}
	s.buf = append(s.buf, str...)
}

// Asciz is the string's bytes plus a NUL.
func (s *Section) Asciz(str string) {
	if s.blocked() {
		return
	}
	s.buf = append(s.buf, str...)
	s.buf = append(s.buf, 0)
}

// Zero is n zero bytes.
func (s *Section) Zero(n int) {
	if s.blocked() || n <= 0 {
		return
	}
	s.buf = append(s.buf, make([]byte, n)...)
}

// Data is raw bytes. It is named Data because Bytes is the read side of the
// contract.
func (s *Section) Data(b []byte) {
	if s.blocked() {
		return
	}
	s.buf = append(s.buf, b...)
}

// Align pads to an n-byte boundary: the canonical multi-byte no-ops in a
// code section, zeros in a data section. n must be a power of two.
func (s *Section) Align(n int) {
	if s.blocked() {
		return
	}
	if n <= 0 || n&(n-1) != 0 {
		s.fail(fmt.Sprintf("align %d", n), ErrAlign)
		return
	}
	pad := -len(s.buf) & (n - 1)
	if pad == 0 {
		return
	}
	if s.kind.IsCode() {
		s.buf = append(s.buf, encode.Nops(pad)...)
	} else {
		s.buf = append(s.buf, make([]byte, pad)...)
	}
}

// ---- instruction placement ----------------------------------------------

// place is where every instruction lands, from a typed helper or from Emit.
// The form is already chosen — a helper pinned it, Emit resolved it — so the
// work here is the gate, the class check the any-typed r/m parameters defer
// to the call, the encode, and the bookkeeping of what the encode left open.
func (s *Section) place(f *isa.Form, o encode.Opts, ops ...any) {
	if s.blocked() {
		return
	}

	if !s.m.feats.Has(f.Need) {
		s.fail(
			fmt.Sprintf("%s requires %s, not in the active feature set", f.Op, f.Need),
			ErrFeature,
			"active: "+s.m.feats.String(),
			"note: enable with "+qualify(s.m.feats.Plus(f.Need).GoExpr()),
		)
		return
	}

	args, err := encode.Args(ops...)
	if err != nil {
		s.fail(err.Error(), ErrForm)
		return
	}
	if !accepts(f, args) {
		s.fail(fmt.Sprintf("operands (%s) do not fit %s", argList(args), f), ErrForm)
		return
	}

	b, fixes, err := encode.Encode(f, o, ops...)
	if err != nil {
		s.fail(err.Error(), sentinelFor(err))
		return
	}

	base := len(s.buf)
	s.buf = append(s.buf, b...)
	for _, fx := range fixes {
		p := pending{
			off:  base + fx.Offset,
			size: fx.Size, tail: fx.Tail,
			pcrel: fx.PCRel, use: fx.Use,
			kind:    RefKind(fx.Kind),
			addend:  fx.Addend,
			sym:     fx.Target.SymName(),
			context: f.String(),
		}
		_, p.label = fx.Target.(operand.Label)
		s.pend = append(s.pend, p)
	}
}

// accepts is the pinned-form class check. Immediate slots are exempt from
// re-matching: the helper pinned the field width, and whether the value fits
// it is encode's range check, not a re-resolution toward a narrower form.
// The literal 1 of the shift forms is passed as Imm(1) and dropped by
// encode's bind, since the opcode already names it.
func accepts(f *isa.Form, args []isa.Arg) bool {
	slots := f.Explicit()
	if len(slots) != len(args) {
		return false
	}
	for i, sl := range slots {
		switch {
		case sl.Class == isa.One:
			if args[i].Kind != isa.KindImm {
				return false
			}
		case sl.Class.IsImm():
			if args[i].Kind != isa.KindImm {
				return false
			}
		default:
			if !sl.Class.Match(args[i]) {
				return false
			}
		}
	}
	return true
}

func argList(args []isa.Arg) string {
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = a.String()
	}
	return strings.Join(parts, ", ")
}

// sentinelFor maps encode/'s vocabulary into this package's.
func sentinelFor(err error) error {
	var imm *encode.ImmediateError
	if errors.As(err, &imm) {
		return ErrRange
	}
	var cnt *encode.CountError
	var op *encode.OperandError
	if errors.As(err, &cnt) || errors.As(err, &op) {
		return ErrForm
	}
	return ErrEncoding
}

// qualify turns feature.Set's unqualified GoExpr into one that compiles at a
// call site: V2.Plus(AVX512VL) -> x86_64.V2.Plus(x86_64.AVX512VL).
func qualify(expr string) string {
	const q = "x86_64."
	expr = strings.ReplaceAll(expr, "()", "\x00")
	expr = q + expr
	expr = strings.ReplaceAll(expr, "(", "("+q)
	expr = strings.ReplaceAll(expr, ", ", ", "+q)
	return strings.ReplaceAll(expr, "\x00", "()")
}

// ---- Finalize -----------------------------------------------------------

func (s *Section) finalize() *Error {
	for _, p := range s.pend {
		// A Label reference to a label in this section, PC-relative:
		// patched and dropped. A SymRef never folds — naming a kind is
		// asking for a relocation, and the split is the naming convention's:
		// JmpLabel resolves here, JmpRef survives into Refs().
		if off, here := s.labels[p.sym]; here && p.label {
			if p.pcrel {
				disp := int64(off) + p.addend - int64(p.off+p.size+p.tail)
				if !dispFits(disp, p.size) {
					return &Error{
						Section: s.Name(), Offset: p.off,
						Context: fmt.Sprintf("%s: %q is %d bytes away, out of range for a %d-byte field",
							p.context, p.sym, disp, p.size),
						Err: ErrRange,
					}
				}
				patchLE(s.buf[p.off:p.off+p.size], disp)
				continue
			}
			// An absolute field holding a same-section address still needs a
			// relocation — no address exists at this layer — and a relocation
			// needs a symbol. A bare label is not one.
			if s.m.symOf[p.sym] != s {
				return &Error{
					Section: s.Name(), Offset: p.off,
					Context: fmt.Sprintf("%s: absolute reference to bare label %q; give it an attribute to make it a symbol",
						p.context, p.sym),
					Err: ErrUndefined,
				}
			}
			s.refs = append(s.refs, s.reference(p))
			continue
		}

		// Cross-section or external. A Label resolves against promoted
		// symbols only — bare labels are per-section namespaces. A SymRef
		// may also resolve to a declared Extern.
		ok := s.m.defined(p.sym)
		if !ok && !p.label {
			ok = s.m.resolvable(p.sym)
		}
		if !ok {
			return &Error{
				Section: s.Name(), Offset: p.off,
				Context: fmt.Sprintf("%s: reference to %q, which is neither defined nor declared Extern",
					p.context, p.sym),
				Err: ErrUndefined,
			}
		}
		s.refs = append(s.refs, s.reference(p))
	}
	s.pend = nil

	// Close every symbol's size at the next symbol or section end.
	for i := range s.syms {
		end := len(s.buf)
		for j := range s.syms {
			if o := s.syms[j].Offset; o > s.syms[i].Offset && o < end {
				end = o
			}
		}
		s.syms[i].Size = end - s.syms[i].Offset
	}
	return nil
}

func (s *Section) reference(p pending) Reference {
	k := p.kind
	if k == RefNone {
		k = inferKind(p)
	}
	return Reference{
		Offset: p.off, Size: p.size, PCRel: p.pcrel,
		Adjust: int32(-p.tail),
		Sym:    p.sym, Kind: k, Addend: p.addend,
	}
}

// inferKind is the encoder's choice when the caller named no kind: a call
// site gets RefPLT, a rip-relative load gets RefPC32, an absolute field gets
// the RefAbs of its width. A Label target that crossed a section boundary is
// internal by construction, so its branch is plain RefPC32, not the PLT.
func inferKind(p pending) RefKind {
	switch p.use {
	case encode.UseBranch:
		if p.label {
			return RefPC32
		}
		return RefPLT
	case encode.UsePCRel:
		return RefPC32
	}
	switch p.size {
	case 8:
		return RefAbs64
	case 4:
		return RefAbs32
	case 2:
		return RefAbs16
	}
	return RefAbs8
}

func dispFits(v int64, size int) bool {
	switch size {
	case 1:
		return v >= -128 && v <= 127
	case 2:
		return v >= -32768 && v <= 32767
	case 4:
		return v >= -2147483648 && v <= 2147483647
	}
	return size == 8
}

func patchLE(b []byte, v int64) {
	for i := range b {
		b[i] = byte(v >> (8 * uint(i)))
	}
}