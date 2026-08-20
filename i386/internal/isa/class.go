// Package isa is the i386 instruction table: every form, its operand
// classes, its opcode bytes, and the level or extension that gates it.
//
// This is the one table. The typed helpers in the arch root are written
// against it, Emit resolves against it, and any introspection enumerates
// it — the filter and the gate are the same values from the same rows, so
// a form that can be named can be encoded and nothing else can.
//
// isa does not encode. It says which form a mnemonic and a list of
// operands resolve to, and what fields that form has; turning a form plus
// operand values into bytes is encode's job. Nothing here holds an
// instruction.
package isa

import (
	"github.com/vertex-language/asm/i386/operand"
	"github.com/vertex-language/asm/i386/reg"
)

// Class is an operand class: the kind of operand a form's slot accepts.
//
// The names are the SDM's, as they appear in an instruction's Operand/
// Instruction column — r/m32, imm8, rel32, CL. A helper's name is built
// from these, so they are also the vocabulary of the typed API.
type Class uint8

const (
	ClassNone Class = iota

	// Registers.
	R8
	R16
	R32
	Sreg
	St
	Mm
	Xmm
	Ymm
	Zmm
	Cr
	Dr

	// Register or memory, by access width.
	RM8
	RM16
	RM32
	RM64
	RM128

	// Memory with no access width of its own: LEA's operand, and the far
	// pointer loads. The address is the operand.
	M

	// Immediates. Imm8S is a sign-extended byte — the imm8 of ADD r/m32,
	// imm8, which is a different form from ADD r/m32, imm32 and four bytes
	// shorter. Keeping them apart is the whole reason the helper names are
	// what they are.
	Imm8
	Imm8S
	Imm16
	Imm32

	// Branch displacements.
	Rel8
	Rel32

	// Fixed operands. These appear in the instruction's syntax but occupy
	// no encoding field, which is what makes ADD EAX, imm32 a distinct
	// six-byte form from ADD r/m32, imm32.
	AL
	AX
	EAX
	CL
	DX
	One // the literal 1 of SHL r/m32, 1
)

var classNames = map[Class]string{
	R8: "r8", R16: "r16", R32: "r32",
	Sreg: "Sreg", St: "ST(i)", Mm: "mm", Xmm: "xmm", Ymm: "ymm", Zmm: "zmm",
	Cr: "CR", Dr: "DR",
	RM8: "r/m8", RM16: "r/m16", RM32: "r/m32", RM64: "r/m64", RM128: "r/m128",
	M:     "m",
	Imm8:  "imm8", Imm8S: "imm8", Imm16: "imm16", Imm32: "imm32",
	Rel8: "rel8", Rel32: "rel32",
	AL: "AL", AX: "AX", EAX: "EAX", CL: "CL", DX: "DX", One: "1",
}

func (c Class) String() string {
	if n, ok := classNames[c]; ok {
		return n
	}
	return "?"
}

// helperNames is the class as it appears in a typed helper's name:
// MovR32Imm32, AddRM32Imm8S, AddEAXImm32. Imm8S and Imm8 differ as forms
// but share an SDM spelling, so the sign-extended form carries the S.
var helperNames = map[Class]string{
	R8: "R8", R16: "R16", R32: "R32",
	Sreg: "Sreg", St: "St", Mm: "Mm", Xmm: "Xmm", Ymm: "Ymm", Zmm: "Zmm",
	Cr: "Cr", Dr: "Dr",
	RM8: "RM8", RM16: "RM16", RM32: "RM32", RM64: "RM64", RM128: "RM128",
	M:     "M",
	Imm8:  "Imm8", Imm8S: "Imm8S", Imm16: "Imm16", Imm32: "Imm32",
	Rel8: "Rel8", Rel32: "Rel32",
	AL: "AL", AX: "AX", EAX: "EAX", CL: "CL", DX: "DX", One: "One",
}

// Accepts reports whether o is the right kind of operand for this class —
// the type question only. For the sized immediate classes it says nothing
// about the value; Fits answers that, and the two are asked separately
// because the two callers need different composites: resolution needs both,
// a pinned helper needs only this one.
//
// The fixed classes (AL, CL, DX, One, ...) check identity here rather than
// type, because "which register" is what defines the form — a form that
// names AL has no field to put another register in. One is identity on the
// value for the same reason: the literal 1 is the operand, not a range.
func (c Class) Accepts(o operand.Operand) bool {
	switch c {
	case R8:
		_, ok := o.(reg.R8)
		return ok
	case R16:
		_, ok := o.(reg.R16)
		return ok
	case R32:
		_, ok := o.(reg.R32)
		return ok
	case Sreg:
		_, ok := o.(reg.Sreg)
		return ok
	case St:
		_, ok := o.(reg.St)
		return ok
	case Mm:
		_, ok := o.(reg.Mm)
		return ok
	case Xmm:
		_, ok := o.(reg.Xmm)
		return ok
	case Ymm:
		_, ok := o.(reg.Ymm)
		return ok
	case Zmm:
		_, ok := o.(reg.Zmm)
		return ok
	case Cr:
		_, ok := o.(reg.Cr)
		return ok
	case Dr:
		_, ok := o.(reg.Dr)
		return ok

	case RM8:
		_, ok := o.(operand.RM8)
		return ok
	case RM16:
		_, ok := o.(operand.RM16)
		return ok
	case RM32:
		_, ok := o.(operand.RM32)
		return ok
	case RM64:
		_, ok := o.(operand.RM64)
		return ok
	case RM128:
		_, ok := o.(operand.RM128)
		return ok

	case M:
		_, ok := o.(operand.Memory)
		return ok

	case Imm8, Imm8S, Imm16, Imm32:
		_, ok := o.(operand.Imm)
		return ok

	case Rel8, Rel32:
		switch o.(type) {
		case operand.Label, operand.SymRef:
			return true
		}
		return false

	case AL:
		return o == operand.Operand(reg.AL)
	case AX:
		return o == operand.Operand(reg.AX)
	case EAX:
		return o == operand.Operand(reg.EAX)
	case CL:
		return o == operand.Operand(reg.CL)
	case DX:
		return o == operand.Operand(reg.DX)
	case One:
		v, ok := o.(operand.Imm)
		return ok && v.Int64() == 1
	}
	return false
}

// Fits reports whether o's value fits this class's field. Only the sized
// immediate classes constrain a value; every other class fits by
// construction once Accepts holds, including One, whose value check is
// identity and lives in Accepts.
func (c Class) Fits(o operand.Operand) bool {
	lo, hi, ok := ImmRange(c)
	if !ok {
		return true
	}
	v, isImm := o.(operand.Imm)
	return isImm && v.Int64() >= lo && v.Int64() <= hi
}

// Matches is Accepts && Fits — what resolution uses. Emit must not pick
// the imm8 form for a value only the imm32 form can hold, so value
// sensitivity belongs in matching: it is what lets ADD r/m32, imm8 win
// over ADD r/m32, imm32 for a small constant, four bytes shorter. A pinned
// helper, by contrast, checks Accepts alone and hands the value to the
// encoder, which is what routes a too-big constant to ErrRange — the
// caller named the form, so "no such form" would be the wrong diagnosis.
func (c Class) Matches(o operand.Operand) bool { return c.Accepts(o) && c.Fits(o) }

// ImmRange is the value range of a sized immediate class, shared by Fits
// and the encoder's range check so the two can never disagree. ok is false
// for every other class, including One, which is identity, not a range.
//
// Imm8 spans -128..255 because both signed and unsigned bytes are written
// against the same one-byte field; Imm8S is the sign-extending byte — the
// group-1 imm8 the processor widens to 32 bits — and spans -128..127,
// which is exactly why it is a distinct class and a distinct helper name.
func ImmRange(c Class) (lo, hi int64, ok bool) {
	switch c {
	case Imm8:
		return -128, 255, true
	case Imm8S:
		return -128, 127, true
	case Imm16:
		return -32768, 65535, true
	case Imm32:
		return -2147483648, 4294967295, true
	}
	return 0, 0, false
}

// Slot is where a form puts an operand in the encoding.
type Slot uint8

const (
	SlotNone   Slot = iota
	SlotReg         // ModRM.reg
	SlotRM          // ModRM.rm, with SIB and displacement as needed
	SlotOpcode      // added to the last opcode byte: +rb, +rw, +rd
	SlotImm         // immediate bytes following
	SlotRel         // a displacement relative to the end of the instruction
	SlotFixed       // named in the syntax, encoded nowhere
)

// Op is one operand of a form: what it accepts and where it goes.
type Op struct {
	Class Class
	Slot  Slot
}