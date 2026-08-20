package aarch64

import (
	"github.com/vertex-language/asm/aarch64/operand"
	"github.com/vertex-language/asm/aarch64/reg"
)

// RegSP64 is an operand slot that reads register 31 as the stack pointer: it
// accepts SP or a numbered reg.X. Go cannot close that union in a method
// signature, so — exactly like x86_64's RM types — it is a documented any,
// and handing it XZR is refused at the call with the diagnostic naming which
// register 31 the slot reads. ImmOrRef is an immediate slot that also takes
// the :lo12: half of an address (a PageOff/GotPageOff AddrRef).
type (
	RegSP64  = any
	RegSP32  = any
	ImmOrRef = any
	TargetOp = any // a SymRef, or an AddrRef wrapping one
)

// ---- Arithmetic, immediate ----

func (s *Section) AddImm64(rd, rn RegSP64, imm ImmOrRef, shift ...operand.ShiftOp) {
	s.inst("AddImm64", withOpt([]any{rd, rn, imm}, shift)...)
}

func (s *Section) SubImm64(rd, rn RegSP64, imm int64, shift ...operand.ShiftOp) {
	s.inst("SubImm64", withOpt([]any{rd, rn, imm}, shift)...)
}

// ---- Arithmetic and logical, register ----

func (s *Section) AddShifted64(rd, rn, rm reg.X, shift ...operand.ShiftOp) {
	s.inst("AddShifted64", withOpt([]any{rd, rn, rm}, shift)...)
}

func (s *Section) OrrImm64(rd RegSP64, rn reg.X, imm uint64) {
	s.inst("OrrImm64", rd, rn, imm)
}

// ---- Move wide ----

func (s *Section) MovzImm64(rd reg.X, imm uint64, shift ...operand.ShiftOp) {
	s.inst("MovzImm64", withOpt([]any{rd, imm}, shift)...)
}

func (s *Section) MovkImm64(rd reg.X, imm uint64, shift ...operand.ShiftOp) {
	s.inst("MovkImm64", withOpt([]any{rd, imm}, shift)...)
}

// ---- PC-relative address ----

func (s *Section) Adr(rd reg.X, target TargetOp)  { s.inst("Adr", target2(rd, target)...) }
func (s *Section) Adrp(rd reg.X, target TargetOp) { s.inst("Adrp", target2(rd, target)...) }

// ---- Loads and stores ----

func (s *Section) LdrImm64(rt reg.X, m operand.Mem) { s.inst("LdrImm64", rt, m) }
func (s *Section) StrImm64(rt reg.X, m operand.Mem) { s.inst("StrImm64", rt, m) }
func (s *Section) Stp64(rt, rt2 reg.X, m operand.Mem)    { s.inst("Stp64", rt, rt2, m) }
func (s *Section) StpPre64(rt, rt2 reg.X, m operand.Mem) { s.inst("StpPre64", rt, rt2, m) }
func (s *Section) Ldp64(rt, rt2 reg.X, m operand.Mem)    { s.inst("Ldp64", rt, rt2, m) }
func (s *Section) LdpPost64(rt, rt2 reg.X, m operand.Mem){ s.inst("LdpPost64", rt, rt2, m) }

// ---- Conditional ----

func (s *Section) Csel64(rd, rn, rm reg.X, c operand.Cond) { s.inst("Csel64", rd, rn, rm, c) }
func (s *Section) Cset64(rd reg.X, c operand.Cond)         { s.inst("Cset64", rd, c) }

// ---- Branches: targets split by where they resolve ----
//
// The Label helpers are same-section, patched at Finalize, no relocation.
// The Ref helpers survive into Refs().

func (s *Section) BLabel(label string)  { s.inst("B", operand.Label(label)) }
func (s *Section) BRef(t TargetOp)      { s.inst("B", t) }
func (s *Section) BlLabel(label string) { s.inst("Bl", operand.Label(label)) }
func (s *Section) BlRef(t TargetOp)     { s.inst("Bl", t) }

func (s *Section) BCondLabel(c operand.Cond, label string) {
	s.inst("BCond", c, operand.Label(label))
}

func (s *Section) CbzLabel64(rt reg.X, label string)  { s.inst("Cbz64", rt, operand.Label(label)) }
func (s *Section) CbnzLabel64(rt reg.X, label string) { s.inst("Cbnz64", rt, operand.Label(label)) }

func (s *Section) Br(rn reg.X)  { s.inst("Br", rn) }
func (s *Section) Blr(rn reg.X) { s.inst("Blr", rn) }

// Ret with its operand omitted is X30 — the encoding's own default.
func (s *Section) Ret(rn ...reg.X) {
	if len(rn) > 0 {
		s.inst("Ret", rn[0])
		return
	}
	s.inst("Ret")
}

func (s *Section) Nop() { s.inst("Nop") }

// ---- small shared shims ----

func withOpt(ops []any, shift []operand.ShiftOp) []any {
	for _, sh := range shift {
		ops = append(ops, sh)
	}
	return ops
}

func target2(rd reg.X, t TargetOp) []any { return []any{rd, t} }