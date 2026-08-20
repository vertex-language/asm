// Package encode turns a resolved form and operand values into bytes.
//
// It is a pure function over resolved forms: it takes an *isa.Form and a
// list of operands and returns bytes and fixups. Nothing survives the call.
// It never sees text, never chooses a different instruction, and never
// reorders anything — the form was chosen before it was called, and this
// package only lays out the fields that form declares.
package encode

import (
	"errors"
	"fmt"

	"github.com/vertex-language/asm/i386/internal/isa"
	"github.com/vertex-language/asm/i386/operand"
	"github.com/vertex-language/asm/i386/reg"
)

// ErrEncode is the sentinel for an operand combination this form cannot
// hold structurally. A value that does not fit a field is a *RangeError,
// which is a different failure: the operands are the right kinds and the
// value is the problem, and the root surfaces it under a different
// sentinel.
var ErrEncode = errors.New("encode")

// RangeError is a value that does not fit the field its form pins.
type RangeError struct {
	Form   string
	Value  int64
	Lo, Hi int64
	Width  int
}

func (e *RangeError) Error() string {
	return fmt.Sprintf("%d does not fit the %d-byte immediate of %s (range %d..%d)",
		e.Value, e.Width, e.Form, e.Lo, e.Hi)
}

// FixupKind distinguishes the two things that get patched later.
type FixupKind uint8

const (
	// FixupLabel is a reference to a name in the same section. It resolves
	// at Finalize as a direct patch and produces no relocation record.
	FixupLabel FixupKind = iota

	// FixupReloc is a reference that leaves the section. It produces a
	// relocation record carrying Kind.
	FixupReloc
)

// Fixup is a field this instruction left for the assembler to fill.
//
// Offset and Size locate the field within the instruction. Adjust is the
// field-position correction: a PC-relative field is resolved against the
// end of the instruction, so the value written is corrected by the bytes
// that follow the field. It is computed here because this is the code that
// placed the field and knows how many bytes follow it.
type Fixup struct {
	Kind   FixupKind
	Offset int
	Size   int
	PCRel  bool
	Adjust int32

	Name   string
	Reloc  operand.RelocKind
	Addend int32
}

// Inst is one encoded instruction.
type Inst struct {
	Bytes  []byte
	Fixups []Fixup
}

// Len is the encoded length in bytes. Emit encodes every candidate form and
// takes the shortest; there is no separate size estimator to disagree with
// the encoder.
func (i Inst) Len() int { return len(i.Bytes) }

// Encode lays out one form with the given operands.
//
// The form is assumed already resolved and already permitted by the feature
// set: the caller answered both questions, and re-asking them here would be
// a second place for the answer to live. Values are judged here — this is
// the code that knows the field widths.
func Encode(f *isa.Form, ops []operand.Operand) (Inst, error) {
	if len(ops) != len(f.Ops) {
		return Inst{}, fmt.Errorf("%w: %s takes %d operands, got %d",
			ErrEncode, f.Signature(), len(f.Ops), len(ops))
	}

	var (
		out    []byte
		fixups []Fixup

		regField  uint8
		haveReg   bool
		rmOperand operand.Operand
		haveRM    bool

		opcodeAdd uint8
		haveAdd   bool

		imms []immField
		rels []relField
	)

	// Pass one: sort operands into the fields they occupy. Nothing is
	// emitted yet, because the prefixes depend on the memory operand and
	// the ModRM byte depends on both halves.
	for i, o := range ops {
		switch f.Ops[i].Slot {
		case isa.SlotReg:
			n, err := regNum(o)
			if err != nil {
				return Inst{}, err
			}
			regField, haveReg = n, true

		case isa.SlotRM:
			rmOperand, haveRM = o, true

		case isa.SlotOpcode:
			n, err := regNum(o)
			if err != nil {
				return Inst{}, err
			}
			opcodeAdd, haveAdd = n, true

		case isa.SlotImm:
			v, ok := o.(operand.Imm)
			if !ok {
				return Inst{}, fmt.Errorf("%w: %s expects an immediate", ErrEncode, f.Signature())
			}
			lo, hi, ok := isa.ImmRange(f.Ops[i].Class)
			if !ok {
				return Inst{}, fmt.Errorf("%w: %s is not an immediate class", ErrEncode, f.Ops[i].Class)
			}
			w := immWidth(f.Ops[i].Class)
			if v.Int64() < lo || v.Int64() > hi {
				return Inst{}, &RangeError{
					Form: f.Signature(), Value: v.Int64(), Lo: lo, Hi: hi, Width: w,
				}
			}
			imms = append(imms, immField{value: v.Int64(), width: w})

		case isa.SlotRel:
			w := 4
			if f.Ops[i].Class == isa.Rel8 {
				w = 1
			}
			rels = append(rels, relField{op: o, width: w})

		case isa.SlotFixed:
			// Named in the syntax, encoded nowhere.
		}
	}

	// A /digit occupies ModRM.reg, so a form cannot have both.
	if f.Ext >= 0 {
		regField, haveReg = uint8(f.Ext), true
	}

	// Prefixes. The memory operand is validated exactly here — its Err is
	// the sticky construction error — and modrm below assumes a valid
	// operand. Group 2 (segment override) precedes group 3 (operand size);
	// the processor accepts any order but only one spelling comes out.
	if haveRM {
		if m, ok := rmOperand.(operand.Memory); ok {
			if err := m.Err(); err != nil {
				return Inst{}, err
			}
			if s, ok := m.Segment(); ok {
				out = append(out, segPrefix(s))
			}
		}
	}
	if f.OpSize16 {
		out = append(out, 0x66)
	}

	// Opcode.
	opcode := append([]byte(nil), f.Opcode...)
	if haveAdd {
		opcode[len(opcode)-1] += opcodeAdd
	}
	out = append(out, opcode...)

	// ModRM, SIB and displacement.
	if haveRM {
		mrm, err := modrm(rmOperand, regField, haveReg)
		if err != nil {
			return Inst{}, err
		}
		out = append(out, mrm.bytes...)
		if mrm.fixup != nil {
			fx := *mrm.fixup
			fx.Offset += len(out) - len(mrm.bytes) + mrm.fixupAt
			fixups = append(fixups, fx)
		}
	} else if haveReg && f.Ext >= 0 {
		return Inst{}, fmt.Errorf("%w: %s has /%d but no r/m operand", ErrEncode, f.Signature(), f.Ext)
	}

	// Immediates.
	for _, im := range imms {
		out = appendLE(out, uint64(im.value), im.width)
	}

	// Branch displacements.
	for _, rl := range rels {
		off := len(out)
		fx := Fixup{Offset: off, Size: rl.width, PCRel: true}
		switch v := rl.op.(type) {
		case operand.Label:
			fx.Kind = FixupLabel
			fx.Name = v.Name()
		case operand.SymRef:
			fx.Kind = FixupReloc
			fx.Name = v.Name()
			fx.Reloc = v.Kind()
			fx.Addend = v.Addend()
		default:
			return Inst{}, fmt.Errorf("%w: %s expects a label or symbol", ErrEncode, f.Signature())
		}
		out = appendLE(out, 0, rl.width)
		fixups = append(fixups, fx)
	}

	// The field-position correction, computed now that the length is known.
	// A PC-relative field is resolved against the end of the instruction,
	// so the correction is the field width plus the bytes that follow it.
	for i := range fixups {
		if fixups[i].PCRel {
			fixups[i].Adjust = -int32(len(out) - fixups[i].Offset - fixups[i].Size)
			fixups[i].Adjust -= int32(fixups[i].Size)
		}
	}

	return Inst{Bytes: out, Fixups: fixups}, nil
}

type immField struct {
	value int64
	width int
}

type relField struct {
	op    operand.Operand
	width int
}

func immWidth(c isa.Class) int {
	switch c {
	case isa.Imm8, isa.Imm8S:
		return 1
	case isa.Imm16:
		return 2
	}
	return 4
}

func appendLE(b []byte, v uint64, width int) []byte {
	for i := 0; i < width; i++ {
		b = append(b, byte(v>>(8*i)))
	}
	return b
}

// regNum is the 0-7 encoding number of any register operand.
func regNum(o operand.Operand) (uint8, error) {
	if r, ok := o.(reg.Reg); ok {
		return r.Num(), nil
	}
	return 0, fmt.Errorf("%w: %T is not a register", ErrEncode, o)
}

// segPrefix is the override byte for a segment. The order of these bytes is
// not the order of the segment encoding numbers, which is why this is a
// switch and not arithmetic on Sreg.Num().
func segPrefix(s reg.Sreg) byte {
	switch s {
	case reg.ES:
		return 0x26
	case reg.CS:
		return 0x2e
	case reg.SS:
		return 0x36
	case reg.DS:
		return 0x3e
	case reg.FS:
		return 0x64
	case reg.GS:
		return 0x65
	}
	return 0
}