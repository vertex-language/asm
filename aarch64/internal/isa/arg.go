package isa

import "github.com/vertex-language/asm/aarch64/reg"

// AddrForm is which of the addressing forms a memory operand uses. Resolve
// needs it because LDR has a different form for each, and they are different
// encodings rather than different immediate ranges.
type AddrForm uint8

const (
	AddrNone      AddrForm = iota
	AddrBase               // [Xn]
	AddrOffset             // [Xn, #imm] — scaled unsigned
	AddrUnscaled           // [Xn, #imm] — unscaled signed, LDUR
	AddrRegOffset          // [Xn, Xm{, LSL #s}] / [Xn, Wm, SXTW #s]
	AddrPreIndex           // [Xn, #imm]!
	AddrPostIndex          // [Xn], #imm
	AddrLiteral            // a PC-relative literal
)

// Arg is one operand as Resolve sees it: enough to decide which form accepts
// it, and nothing more. The concrete operand value stays with the caller;
// encode/ is what lowers it.
type Arg struct {
	Class Class

	// Num is the register number, for the classes that name a register. It
	// matters because register 31 is two registers and the class alone does
	// not say which one an argument is.
	Num uint16

	// Imm is an immediate's value, for slots whose acceptance depends on it.
	Imm int64

	// Arr and Elem describe a vector operand.
	Arr  reg.Arrangement
	Elem reg.Elem

	// Addr is the addressing form of a memory operand.
	Addr AddrForm

	// AccessBits is the width a memory operand is accessed at, when the caller
	// stated it (Mem64 rather than Mem). Zero means the form decides.
	AccessBits uint16
}

// ArgOf builds an Arg from a register.
func ArgOf(r reg.Reg) Arg {
	a := Arg{Num: r.Num()}
	switch v := r.(type) {
	case reg.X:
		a.Class = ClassX
	case reg.W:
		a.Class = ClassW
	case reg.Xsp:
		a.Class = ClassXsp
	case reg.Wsp:
		a.Class = ClassWsp
	case reg.V:
		a.Class = ClassV
	case reg.Q:
		a.Class = ClassQ
	case reg.D:
		a.Class = ClassD
	case reg.S:
		a.Class = ClassS
	case reg.H:
		a.Class = ClassH
	case reg.B:
		a.Class = ClassB
	case reg.Vec:
		a.Class, a.Arr = ClassVArr, v.A
	case reg.VLane:
		a.Class, a.Elem = ClassVLane, v.E
	case reg.Z:
		a.Class = ClassZ
	case reg.P:
		a.Class = ClassP
	case reg.Sys:
		a.Class = ClassSys
	}
	return a
}

// ImmArg builds an immediate argument.
func ImmArg(v int64) Arg { return Arg{Class: ClassImm, Imm: v} }

// MemArg builds a memory argument.
func MemArg(form AddrForm, accessBits uint16) Arg {
	return Arg{Class: memClass(accessBits), Addr: form, AccessBits: accessBits}
}

// LabelArg builds a branch or address target.
func LabelArg() Arg { return Arg{Class: ClassLabel} }

// Match reports whether a slot's class accepts an argument.
//
// Two asymmetries are deliberate.
//
// An Xsp slot accepts a numbered X: "add x0, x1, #1" is legal, and a caller
// building it holds x1 as a plain X with no reason to have converted it. It
// does not accept XZR, which is a different register that happens to share an
// encoding.
//
// An X slot does not accept SP. Register 31 in such a slot is the zero
// register, and a caller who wrote SP meant the other one.
//
// A third acceptance is symmetric with those two and equally deliberate: an
// immediate slot accepts an address reference, and a label slot accepts a
// bare immediate. Both directions are the same fact seen from opposite ends —
// a field can hold a number the caller computed or a number a linker will —
// and each direction has its own justification.
//
// ClassLabel taking ClassImm is a branch to a displacement the caller
// computed and takes responsibility for.
//
// ClassImm taking ClassLabel is the :lo12: half of an address landing in a
// data-processing immediate: add x0, x0, :lo12:msg is the second instruction
// of the ADRP/ADD pair, and its field is an ordinary Imm12 whose value is not
// a number yet. Resolve only decides which form the operand lands in; whether
// the reference's role is one that may land there is encode/'s check, which
// has the role in hand and refuses everything but the page-offset ones.
//
// This widening cannot create ambiguity. Two forms of one mnemonic differing
// only by imm-versus-label in the same slot would now both match one argument
// list — and checkTable's signature collision check is what guarantees no such
// pair exists, panicking at init rather than leaving Resolve a choice to make.
func (c Class) Match(a Arg) bool {
	switch c {
	case ClassXsp:
		return a.Class == ClassXsp ||
			(a.Class == ClassX && a.Num != 31)
	case ClassWsp:
		return a.Class == ClassWsp ||
			(a.Class == ClassW && a.Num != 31)
	case ClassX, ClassW:
		return a.Class == c
	case ClassPg:
		return a.Class == ClassP && a.Num <= 7
	case ClassMem8, ClassMem16, ClassMem32, ClassMem64, ClassMem128:
		// A caller that stated a width must match it; one that did not takes
		// the form's.
		if a.AccessBits != 0 && a.AccessBits != c.AccessBits() {
			return false
		}
		return a.Class.Mem()
	case ClassImm:
		// A symbolic value whose number arrives at link time may fill an
		// immediate field; see above.
		return a.Class == ClassImm || a.Class == ClassLabel
	case ClassLabel:
		// A branch target may arrive as a label or as a bare immediate
		// displacement.
		return a.Class == ClassLabel || a.Class == ClassImm
	}
	return a.Class == c
}