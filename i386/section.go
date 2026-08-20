package i386

import (
	"errors"
	"fmt"

	"github.com/vertex-language/asm/i386/internal/encode"
	"github.com/vertex-language/asm/i386/internal/isa"
	"github.com/vertex-language/asm/i386/operand"
)

// Section is an append-only byte buffer plus a reference list. Nothing in it
// is symbolic-and-rewritable; no pass can run over it.
type Section struct {
	m    *Module
	kind SectionKind
	buf  []byte

	labels      map[string]int // every label, bare or attributed
	labelFixups []labelFixup   // same-section holes, patched at Finalize
	refs        []Reference    // holes that survive into Refs()
	syms        []Symbol       // attributed labels, in offset order
}

// labelFixup is an encode.Fixup translated into section coordinates.
type labelFixup struct {
	offset int
	size   int
	pcrel  bool
	adjust int32
	addend int32
	name   string
}

// Kind is the section's load-time kind.
func (s *Section) Kind() SectionKind { return s.kind }

// Name is the conventional spelling of the kind: ".text", ".rodata", ...
func (s *Section) Name() string { return s.kind.String() }

// Offset is the current end of the section: the offset the next byte will
// land at, and the value a Label placed now would name. It is exported
// because jump tables and literal pools are the caller's to build — this
// package deliberately refuses to build them — and building them requires
// knowing where you are.
func (s *Section) Offset() int { return len(s.buf) }

// Bytes is the finished machine code. After Finalize, same-section label
// references are patched in; the remaining holes are zero-filled and listed
// by Refs.
func (s *Section) Bytes() []byte {
	out := make([]byte, len(s.buf))
	copy(out, s.buf)
	return out
}

// Refs is every hole a linker must fill.
func (s *Section) Refs() []Reference {
	out := make([]Reference, len(s.refs))
	copy(out, s.refs)
	return out
}

// Symbols is every attributed label, sizes closed at Finalize.
func (s *Section) Symbols() []Symbol {
	out := make([]Symbol, len(s.syms))
	copy(out, s.syms)
	return out
}

func (s *Section) ok() bool {
	if s.m.finalized {
		s.m.failAt(s.Name(), len(s.buf), ErrFinalized, nil, nil, "builder call after Finalize")
		return false
	}
	return s.m.err == nil
}

// fail records a positioned error with no cause and no notes — the common
// case for failures this package diagnoses itself.
func (s *Section) fail(sentinel error, format string, args ...any) {
	s.failWith(sentinel, nil, nil, format, args...)
}

// failWith records a positioned error carrying an underlying cause and
// notes. The cause joins the chain for errors.Is; anything a caller might
// need from its internal type is restated in the notes as text.
func (s *Section) failWith(sentinel, cause error, notes []string, format string, args ...any) {
	s.m.failAt(s.Name(), len(s.buf), sentinel, cause, notes, format, args...)
}

// failEncode classifies an encoder error under the right sentinel. A
// *encode.RangeError is a value that does not fit the field the named form
// pins — ErrRange, with the field width and range in the notes. Everything
// else — a sticky operand-construction error, an internal mismatch — is
// ErrForm, with the encoder's error as the cause.
func (s *Section) failEncode(sig string, err error) {
	var re *encode.RangeError
	if errors.As(err, &re) {
		s.failWith(ErrRange, err, []string{
			fmt.Sprintf("the immediate field of %s is %d bytes; the range is %d..%d",
				re.Form, re.Width, re.Lo, re.Hi),
		}, "%d does not fit %s", re.Value, sig)
		return
	}
	s.failWith(ErrForm, err, nil, "%v", err)
}

// Label names the current offset. Bare, it is a branch target in this
// section's namespace only. Any attribute promotes it into Symbols().
func (s *Section) Label(name string, attrs ...SymbolAttr) {
	if !s.ok() {
		return
	}
	if _, dup := s.labels[name]; dup {
		s.fail(ErrDuplicate, "label %q defined twice in %s", name, s.Name())
		return
	}
	s.labels[name] = len(s.buf)
	if len(attrs) == 0 {
		return
	}
	sym := Symbol{Name: name, Offset: len(s.buf)}
	for _, a := range attrs {
		a.applySymbol(&sym)
	}
	s.syms = append(s.syms, sym)
}

// Align pads to an n-byte boundary: multi-byte nops in a Text section, gated
// by the module's feature set, and zeros everywhere else.
func (s *Section) Align(n int) {
	if !s.ok() {
		return
	}
	if n <= 1 {
		return
	}
	if n&(n-1) != 0 {
		s.fail(ErrAlign, "alignment %d is not a power of two", n)
		return
	}
	pad := (n - len(s.buf)%n) % n
	if pad == 0 {
		return
	}
	if s.kind == Text {
		s.buf = append(s.buf, encode.Nops(pad, s.m.features)...)
		return
	}
	s.buf = append(s.buf, make([]byte, pad)...)
}

// Data appends raw bytes. (The README's builder spelling `Bytes(blob)`
// collides with the reader `Bytes()`; the reader keeps the name because it
// is the downstream contract.)
func (s *Section) Data(b []byte) {
	if !s.ok() {
		return
	}
	s.buf = append(s.buf, b...)
}

// Byte appends one byte.
func (s *Section) Byte(v byte) {
	if !s.ok() {
		return
	}
	s.buf = append(s.buf, v)
}

// Long appends a little-endian 4-byte value.
func (s *Section) Long(v uint32) {
	if !s.ok() {
		return
	}
	s.buf = append(s.buf, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
}

// Quad appends a little-endian 8-byte value.
func (s *Section) Quad(v uint64) {
	if !s.ok() {
		return
	}
	for i := 0; i < 8; i++ {
		s.buf = append(s.buf, byte(v>>(8*i)))
	}
}

// Ascii appends a string with no terminator.
func (s *Section) Ascii(str string) {
	if !s.ok() {
		return
	}
	s.buf = append(s.buf, str...)
}

// Asciz appends a string and a NUL.
func (s *Section) Asciz(str string) {
	if !s.ok() {
		return
	}
	s.buf = append(s.buf, str...)
	s.buf = append(s.buf, 0)
}

// Zero appends n zero bytes.
func (s *Section) Zero(n int) {
	if !s.ok() {
		return
	}
	if n <= 0 {
		return
	}
	s.buf = append(s.buf, make([]byte, n)...)
}

// inst is what every typed helper funnels into: gate, type-check, encode,
// place. The form is pinned by the caller.
//
// The operand check is AcceptsTypes, not Matches: a helper named its form,
// so the only question left about a value is whether it fits the field that
// form declares — the encoder's range check, surfaced as ErrRange by
// failEncode. Asking Matches here would misreport a too-big constant as
// "operands do not fit", sending the caller hunting for a type mistake that
// is not there.
func (s *Section) inst(f *isa.Form, ops ...operand.Operand) {
	if !s.ok() {
		return
	}
	if !f.Enabled(s.m.features) {
		s.fail(ErrFeature, "%s is gated behind %s (module has %s)",
			f.Signature(), f.Gate(), s.m.features)
		return
	}
	if !f.AcceptsTypes(ops) {
		s.fail(ErrForm, "operands %s do not fit %s", opsString(ops), f.Signature())
		return
	}
	in, err := encode.Encode(f, ops)
	if err != nil {
		s.failEncode(f.Signature(), err)
		return
	}
	s.place(in)
}

// place appends an encoded instruction and translates its fixups into
// section coordinates. Adjust arrives already computed by encode/; this is
// transcription, not decision.
func (s *Section) place(in encode.Inst) {
	base := len(s.buf)
	s.buf = append(s.buf, in.Bytes...)
	for _, fx := range in.Fixups {
		switch fx.Kind {
		case encode.FixupLabel:
			s.labelFixups = append(s.labelFixups, labelFixup{
				offset: base + fx.Offset,
				size:   fx.Size,
				pcrel:  fx.PCRel,
				adjust: fx.Adjust,
				addend: fx.Addend,
				name:   fx.Name,
			})
		case encode.FixupReloc:
			s.refs = append(s.refs, Reference{
				Offset: base + fx.Offset,
				Size:   fx.Size,
				PCRel:  fx.PCRel,
				Adjust: fx.Adjust,
				Sym:    fx.Name,
				Kind:   RefKind(fx.Reloc),
				Addend: fx.Addend,
			})
		}
	}
}

// patchLabels resolves every same-section label hole. Both addresses are in
// this section, so the section's final address cancels and the patch is
// computable now. Reports success; failure is recorded on the module.
func (s *Section) patchLabels() bool {
	for _, fx := range s.labelFixups {
		target, ok := s.labels[fx.name]
		if !ok {
			s.m.failAt(s.Name(), fx.offset, ErrUndefined, nil, nil,
				"label %q is not defined in %s", fx.name, s.Name())
			return false
		}
		if !fx.pcrel {
			// Labels only reach here through SlotRel, which is always
			// PC-relative; anything else is a bug in this package.
			s.m.failAt(s.Name(), fx.offset, ErrForm, nil, nil,
				"internal: non-PC-relative label fixup for %q", fx.name)
			return false
		}
		v := int64(target) - int64(fx.offset) + int64(fx.adjust) + int64(fx.addend)
		var lo, hi int64
		switch fx.size {
		case 1:
			lo, hi = -128, 127
		case 2:
			lo, hi = -32768, 32767
		default:
			lo, hi = -2147483648, 2147483647
		}
		if v < lo || v > hi {
			s.m.failAt(s.Name(), fx.offset, ErrRange, nil,
				[]string{fmt.Sprintf("a %d-byte displacement reaches %d..%d", fx.size, lo, hi)},
				"branch to %q is %d bytes away, out of range for a %d-byte displacement",
				fx.name, v, fx.size)
			return false
		}
		for i := 0; i < fx.size; i++ {
			s.buf[fx.offset+i] = byte(uint64(v) >> (8 * i))
		}
	}
	s.labelFixups = nil
	return true
}

// closeSizes closes each symbol at the next symbol or section end. Symbols
// were appended at increasing offsets, so this is a single walk.
func (s *Section) closeSizes() {
	for i := range s.syms {
		end := len(s.buf)
		if i+1 < len(s.syms) {
			end = s.syms[i+1].Offset
		}
		s.syms[i].Size = end - s.syms[i].Offset
	}
}