package i386

// Typed helpers, one method per form in the isa table, expanded by hand from
// isa.All(). Binding is by HelperName, so appending rows to the table breaks
// nothing here; a removed or renamed form panics at the first call, by name.
// New tranches (table_sse.go, ...) mean appending methods to this file.

import (
	"github.com/vertex-language/asm/i386/internal/isa"
	"github.com/vertex-language/asm/i386/operand"
	"github.com/vertex-language/asm/i386/reg"
)

var helperForms = func() map[string]*isa.Form {
	m := make(map[string]*isa.Form, len(isa.All()))
	for _, f := range isa.All() {
		if _, dup := m[f.HelperName()]; dup {
			panic("i386: duplicate helper name " + f.HelperName())
		}
		m[f.HelperName()] = f
	}
	return m
}()

func form(name string) *isa.Form {
	f, ok := helperForms[name]
	if !ok {
		panic("i386: no form " + name + " in the isa table; inst.go is stale")
	}
	return f
}

func imm(v int64) operand.Imm    { return operand.NewImm(v) }
func lbl(name string) operand.Label { return operand.NewLabel(name) }

// ── ALU group: 0x00–0x3f pattern plus group 1 ───────────────────────────────

func (s *Section) AddRM8R8(a operand.RM8, b reg.R8)     { s.inst(form("AddRM8R8"), a, b) }
func (s *Section) AddRM32R32(a operand.RM32, b reg.R32) { s.inst(form("AddRM32R32"), a, b) }
func (s *Section) AddR8RM8(a reg.R8, b operand.RM8)     { s.inst(form("AddR8RM8"), a, b) }
func (s *Section) AddR32RM32(a reg.R32, b operand.RM32) { s.inst(form("AddR32RM32"), a, b) }
func (s *Section) AddALImm8(v int64)                    { s.inst(form("AddALImm8"), reg.AL, imm(v)) }
func (s *Section) AddEAXImm32(v int64)                  { s.inst(form("AddEAXImm32"), reg.EAX, imm(v)) }
func (s *Section) AddRM8Imm8(a operand.RM8, v int64)    { s.inst(form("AddRM8Imm8"), a, imm(v)) }
func (s *Section) AddRM32Imm8S(a operand.RM32, v int64) { s.inst(form("AddRM32Imm8S"), a, imm(v)) }
func (s *Section) AddRM32Imm32(a operand.RM32, v int64) { s.inst(form("AddRM32Imm32"), a, imm(v)) }

func (s *Section) OrRM8R8(a operand.RM8, b reg.R8)     { s.inst(form("OrRM8R8"), a, b) }
func (s *Section) OrRM32R32(a operand.RM32, b reg.R32) { s.inst(form("OrRM32R32"), a, b) }
func (s *Section) OrR8RM8(a reg.R8, b operand.RM8)     { s.inst(form("OrR8RM8"), a, b) }
func (s *Section) OrR32RM32(a reg.R32, b operand.RM32) { s.inst(form("OrR32RM32"), a, b) }
func (s *Section) OrALImm8(v int64)                    { s.inst(form("OrALImm8"), reg.AL, imm(v)) }
func (s *Section) OrEAXImm32(v int64)                  { s.inst(form("OrEAXImm32"), reg.EAX, imm(v)) }
func (s *Section) OrRM8Imm8(a operand.RM8, v int64)    { s.inst(form("OrRM8Imm8"), a, imm(v)) }
func (s *Section) OrRM32Imm8S(a operand.RM32, v int64) { s.inst(form("OrRM32Imm8S"), a, imm(v)) }
func (s *Section) OrRM32Imm32(a operand.RM32, v int64) { s.inst(form("OrRM32Imm32"), a, imm(v)) }

func (s *Section) AdcRM8R8(a operand.RM8, b reg.R8)     { s.inst(form("AdcRM8R8"), a, b) }
func (s *Section) AdcRM32R32(a operand.RM32, b reg.R32) { s.inst(form("AdcRM32R32"), a, b) }
func (s *Section) AdcR8RM8(a reg.R8, b operand.RM8)     { s.inst(form("AdcR8RM8"), a, b) }
func (s *Section) AdcR32RM32(a reg.R32, b operand.RM32) { s.inst(form("AdcR32RM32"), a, b) }
func (s *Section) AdcALImm8(v int64)                    { s.inst(form("AdcALImm8"), reg.AL, imm(v)) }
func (s *Section) AdcEAXImm32(v int64)                  { s.inst(form("AdcEAXImm32"), reg.EAX, imm(v)) }
func (s *Section) AdcRM8Imm8(a operand.RM8, v int64)    { s.inst(form("AdcRM8Imm8"), a, imm(v)) }
func (s *Section) AdcRM32Imm8S(a operand.RM32, v int64) { s.inst(form("AdcRM32Imm8S"), a, imm(v)) }
func (s *Section) AdcRM32Imm32(a operand.RM32, v int64) { s.inst(form("AdcRM32Imm32"), a, imm(v)) }

func (s *Section) SbbRM8R8(a operand.RM8, b reg.R8)     { s.inst(form("SbbRM8R8"), a, b) }
func (s *Section) SbbRM32R32(a operand.RM32, b reg.R32) { s.inst(form("SbbRM32R32"), a, b) }
func (s *Section) SbbR8RM8(a reg.R8, b operand.RM8)     { s.inst(form("SbbR8RM8"), a, b) }
func (s *Section) SbbR32RM32(a reg.R32, b operand.RM32) { s.inst(form("SbbR32RM32"), a, b) }
func (s *Section) SbbALImm8(v int64)                    { s.inst(form("SbbALImm8"), reg.AL, imm(v)) }
func (s *Section) SbbEAXImm32(v int64)                  { s.inst(form("SbbEAXImm32"), reg.EAX, imm(v)) }
func (s *Section) SbbRM8Imm8(a operand.RM8, v int64)    { s.inst(form("SbbRM8Imm8"), a, imm(v)) }
func (s *Section) SbbRM32Imm8S(a operand.RM32, v int64) { s.inst(form("SbbRM32Imm8S"), a, imm(v)) }
func (s *Section) SbbRM32Imm32(a operand.RM32, v int64) { s.inst(form("SbbRM32Imm32"), a, imm(v)) }

func (s *Section) AndRM8R8(a operand.RM8, b reg.R8)     { s.inst(form("AndRM8R8"), a, b) }
func (s *Section) AndRM32R32(a operand.RM32, b reg.R32) { s.inst(form("AndRM32R32"), a, b) }
func (s *Section) AndR8RM8(a reg.R8, b operand.RM8)     { s.inst(form("AndR8RM8"), a, b) }
func (s *Section) AndR32RM32(a reg.R32, b operand.RM32) { s.inst(form("AndR32RM32"), a, b) }
func (s *Section) AndALImm8(v int64)                    { s.inst(form("AndALImm8"), reg.AL, imm(v)) }
func (s *Section) AndEAXImm32(v int64)                  { s.inst(form("AndEAXImm32"), reg.EAX, imm(v)) }
func (s *Section) AndRM8Imm8(a operand.RM8, v int64)    { s.inst(form("AndRM8Imm8"), a, imm(v)) }
func (s *Section) AndRM32Imm8S(a operand.RM32, v int64) { s.inst(form("AndRM32Imm8S"), a, imm(v)) }
func (s *Section) AndRM32Imm32(a operand.RM32, v int64) { s.inst(form("AndRM32Imm32"), a, imm(v)) }

func (s *Section) SubRM8R8(a operand.RM8, b reg.R8)     { s.inst(form("SubRM8R8"), a, b) }
func (s *Section) SubRM32R32(a operand.RM32, b reg.R32) { s.inst(form("SubRM32R32"), a, b) }
func (s *Section) SubR8RM8(a reg.R8, b operand.RM8)     { s.inst(form("SubR8RM8"), a, b) }
func (s *Section) SubR32RM32(a reg.R32, b operand.RM32) { s.inst(form("SubR32RM32"), a, b) }
func (s *Section) SubALImm8(v int64)                    { s.inst(form("SubALImm8"), reg.AL, imm(v)) }
func (s *Section) SubEAXImm32(v int64)                  { s.inst(form("SubEAXImm32"), reg.EAX, imm(v)) }
func (s *Section) SubRM8Imm8(a operand.RM8, v int64)    { s.inst(form("SubRM8Imm8"), a, imm(v)) }
func (s *Section) SubRM32Imm8S(a operand.RM32, v int64) { s.inst(form("SubRM32Imm8S"), a, imm(v)) }
func (s *Section) SubRM32Imm32(a operand.RM32, v int64) { s.inst(form("SubRM32Imm32"), a, imm(v)) }

func (s *Section) XorRM8R8(a operand.RM8, b reg.R8)     { s.inst(form("XorRM8R8"), a, b) }
func (s *Section) XorRM32R32(a operand.RM32, b reg.R32) { s.inst(form("XorRM32R32"), a, b) }
func (s *Section) XorR8RM8(a reg.R8, b operand.RM8)     { s.inst(form("XorR8RM8"), a, b) }
func (s *Section) XorR32RM32(a reg.R32, b operand.RM32) { s.inst(form("XorR32RM32"), a, b) }
func (s *Section) XorALImm8(v int64)                    { s.inst(form("XorALImm8"), reg.AL, imm(v)) }
func (s *Section) XorEAXImm32(v int64)                  { s.inst(form("XorEAXImm32"), reg.EAX, imm(v)) }
func (s *Section) XorRM8Imm8(a operand.RM8, v int64)    { s.inst(form("XorRM8Imm8"), a, imm(v)) }
func (s *Section) XorRM32Imm8S(a operand.RM32, v int64) { s.inst(form("XorRM32Imm8S"), a, imm(v)) }
func (s *Section) XorRM32Imm32(a operand.RM32, v int64) { s.inst(form("XorRM32Imm32"), a, imm(v)) }

func (s *Section) CmpRM8R8(a operand.RM8, b reg.R8)     { s.inst(form("CmpRM8R8"), a, b) }
func (s *Section) CmpRM32R32(a operand.RM32, b reg.R32) { s.inst(form("CmpRM32R32"), a, b) }
func (s *Section) CmpR8RM8(a reg.R8, b operand.RM8)     { s.inst(form("CmpR8RM8"), a, b) }
func (s *Section) CmpR32RM32(a reg.R32, b operand.RM32) { s.inst(form("CmpR32RM32"), a, b) }
func (s *Section) CmpALImm8(v int64)                    { s.inst(form("CmpALImm8"), reg.AL, imm(v)) }
func (s *Section) CmpEAXImm32(v int64)                  { s.inst(form("CmpEAXImm32"), reg.EAX, imm(v)) }
func (s *Section) CmpRM8Imm8(a operand.RM8, v int64)    { s.inst(form("CmpRM8Imm8"), a, imm(v)) }
func (s *Section) CmpRM32Imm8S(a operand.RM32, v int64) { s.inst(form("CmpRM32Imm8S"), a, imm(v)) }
func (s *Section) CmpRM32Imm32(a operand.RM32, v int64) { s.inst(form("CmpRM32Imm32"), a, imm(v)) }

// ── MOV family, LEA, extend, exchange ────────────────────────────────────────

func (s *Section) MovRM8R8(a operand.RM8, b reg.R8)     { s.inst(form("MovRM8R8"), a, b) }
func (s *Section) MovRM32R32(a operand.RM32, b reg.R32) { s.inst(form("MovRM32R32"), a, b) }
func (s *Section) MovR8RM8(a reg.R8, b operand.RM8)     { s.inst(form("MovR8RM8"), a, b) }
func (s *Section) MovR32RM32(a reg.R32, b operand.RM32) { s.inst(form("MovR32RM32"), a, b) }
func (s *Section) MovR8Imm8(a reg.R8, v int64)          { s.inst(form("MovR8Imm8"), a, imm(v)) }
func (s *Section) MovR32Imm32(a reg.R32, v int64)       { s.inst(form("MovR32Imm32"), a, imm(v)) }
func (s *Section) MovRM8Imm8(a operand.RM8, v int64)    { s.inst(form("MovRM8Imm8"), a, imm(v)) }
func (s *Section) MovRM32Imm32(a operand.RM32, v int64) { s.inst(form("MovRM32Imm32"), a, imm(v)) }
func (s *Section) MovRM32Sreg(a operand.RM32, b reg.Sreg) { s.inst(form("MovRM32Sreg"), a, b) }
func (s *Section) MovSregRM32(a reg.Sreg, b operand.RM32) { s.inst(form("MovSregRM32"), a, b) }

func (s *Section) LeaR32M(a reg.R32, b operand.Memory) { s.inst(form("LeaR32M"), a, b) }

func (s *Section) MovzxR32RM8(a reg.R32, b operand.RM8)   { s.inst(form("MovzxR32RM8"), a, b) }
func (s *Section) MovzxR32RM16(a reg.R32, b operand.RM16) { s.inst(form("MovzxR32RM16"), a, b) }
func (s *Section) MovsxR32RM8(a reg.R32, b operand.RM8)   { s.inst(form("MovsxR32RM8"), a, b) }
func (s *Section) MovsxR32RM16(a reg.R32, b operand.RM16) { s.inst(form("MovsxR32RM16"), a, b) }

func (s *Section) XchgEAXR32(a reg.R32)                  { s.inst(form("XchgEAXR32"), reg.EAX, a) }
func (s *Section) XchgRM8R8(a operand.RM8, b reg.R8)     { s.inst(form("XchgRM8R8"), a, b) }
func (s *Section) XchgRM32R32(a operand.RM32, b reg.R32) { s.inst(form("XchgRM32R32"), a, b) }

func (s *Section) BswapR32(a reg.R32)                        { s.inst(form("BswapR32"), a) }
func (s *Section) CmpxchgRM8R8(a operand.RM8, b reg.R8)      { s.inst(form("CmpxchgRM8R8"), a, b) }
func (s *Section) CmpxchgRM32R32(a operand.RM32, b reg.R32)  { s.inst(form("CmpxchgRM32R32"), a, b) }
func (s *Section) XaddRM8R8(a operand.RM8, b reg.R8)         { s.inst(form("XaddRM8R8"), a, b) }
func (s *Section) XaddRM32R32(a operand.RM32, b reg.R32)     { s.inst(form("XaddRM32R32"), a, b) }
func (s *Section) Cmpxchg8bRM64(a operand.RM64)              { s.inst(form("Cmpxchg8bRM64"), a) }

// ── Stack ────────────────────────────────────────────────────────────────────

func (s *Section) PushR32(a reg.R32)      { s.inst(form("PushR32"), a) }
func (s *Section) PushRM32(a operand.RM32) { s.inst(form("PushRM32"), a) }
func (s *Section) PushImm8S(v int64)      { s.inst(form("PushImm8S"), imm(v)) }
func (s *Section) PushImm32(v int64)      { s.inst(form("PushImm32"), imm(v)) }
func (s *Section) PopR32(a reg.R32)       { s.inst(form("PopR32"), a) }
func (s *Section) PopRM32(a operand.RM32) { s.inst(form("PopRM32"), a) }

func (s *Section) Pusha() { s.inst(form("Pusha")) }
func (s *Section) Popa()  { s.inst(form("Popa")) }
func (s *Section) Pushf() { s.inst(form("Pushf")) }
func (s *Section) Popf()  { s.inst(form("Popf")) }

func (s *Section) EnterImm16Imm8(a, b int64) { s.inst(form("EnterImm16Imm8"), imm(a), imm(b)) }
func (s *Section) Leave()                    { s.inst(form("Leave")) }

// ── Arithmetic groups 3, 4, 5; TEST; sign extension ─────────────────────────

func (s *Section) TestRM8Imm8(a operand.RM8, v int64)    { s.inst(form("TestRM8Imm8"), a, imm(v)) }
func (s *Section) TestRM32Imm32(a operand.RM32, v int64) { s.inst(form("TestRM32Imm32"), a, imm(v)) }
func (s *Section) NotRM8(a operand.RM8)                  { s.inst(form("NotRM8"), a) }
func (s *Section) NotRM32(a operand.RM32)                { s.inst(form("NotRM32"), a) }
func (s *Section) NegRM8(a operand.RM8)                  { s.inst(form("NegRM8"), a) }
func (s *Section) NegRM32(a operand.RM32)                { s.inst(form("NegRM32"), a) }
func (s *Section) MulRM8(a operand.RM8)                  { s.inst(form("MulRM8"), a) }
func (s *Section) MulRM32(a operand.RM32)                { s.inst(form("MulRM32"), a) }
func (s *Section) ImulRM8(a operand.RM8)                 { s.inst(form("ImulRM8"), a) }
func (s *Section) ImulRM32(a operand.RM32)               { s.inst(form("ImulRM32"), a) }
func (s *Section) DivRM8(a operand.RM8)                  { s.inst(form("DivRM8"), a) }
func (s *Section) DivRM32(a operand.RM32)                { s.inst(form("DivRM32"), a) }
func (s *Section) IdivRM8(a operand.RM8)                 { s.inst(form("IdivRM8"), a) }
func (s *Section) IdivRM32(a operand.RM32)               { s.inst(form("IdivRM32"), a) }

func (s *Section) ImulR32RM32(a reg.R32, b operand.RM32) { s.inst(form("ImulR32RM32"), a, b) }
func (s *Section) ImulR32RM32Imm8S(a reg.R32, b operand.RM32, v int64) {
	s.inst(form("ImulR32RM32Imm8S"), a, b, imm(v))
}
func (s *Section) ImulR32RM32Imm32(a reg.R32, b operand.RM32, v int64) {
	s.inst(form("ImulR32RM32Imm32"), a, b, imm(v))
}

func (s *Section) IncR32(a reg.R32)       { s.inst(form("IncR32"), a) }
func (s *Section) IncRM8(a operand.RM8)   { s.inst(form("IncRM8"), a) }
func (s *Section) IncRM32(a operand.RM32) { s.inst(form("IncRM32"), a) }
func (s *Section) DecR32(a reg.R32)       { s.inst(form("DecR32"), a) }
func (s *Section) DecRM8(a operand.RM8)   { s.inst(form("DecRM8"), a) }
func (s *Section) DecRM32(a operand.RM32) { s.inst(form("DecRM32"), a) }

func (s *Section) TestRM8R8(a operand.RM8, b reg.R8)     { s.inst(form("TestRM8R8"), a, b) }
func (s *Section) TestRM32R32(a operand.RM32, b reg.R32) { s.inst(form("TestRM32R32"), a, b) }
func (s *Section) TestALImm8(v int64)                    { s.inst(form("TestALImm8"), reg.AL, imm(v)) }
func (s *Section) TestEAXImm32(v int64)                  { s.inst(form("TestEAXImm32"), reg.EAX, imm(v)) }

func (s *Section) Cwde() { s.inst(form("Cwde")) }
func (s *Section) Cdq()  { s.inst(form("Cdq")) }

// ── Shifts and rotates ───────────────────────────────────────────────────────

func (s *Section) RolRM8One(a operand.RM8)             { s.inst(form("RolRM8One"), a, imm(1)) }
func (s *Section) RolRM32One(a operand.RM32)           { s.inst(form("RolRM32One"), a, imm(1)) }
func (s *Section) RolRM8CL(a operand.RM8)              { s.inst(form("RolRM8CL"), a, reg.CL) }
func (s *Section) RolRM32CL(a operand.RM32)            { s.inst(form("RolRM32CL"), a, reg.CL) }
func (s *Section) RolRM8Imm8(a operand.RM8, v int64)   { s.inst(form("RolRM8Imm8"), a, imm(v)) }
func (s *Section) RolRM32Imm8(a operand.RM32, v int64) { s.inst(form("RolRM32Imm8"), a, imm(v)) }

func (s *Section) RorRM8One(a operand.RM8)             { s.inst(form("RorRM8One"), a, imm(1)) }
func (s *Section) RorRM32One(a operand.RM32)           { s.inst(form("RorRM32One"), a, imm(1)) }
func (s *Section) RorRM8CL(a operand.RM8)              { s.inst(form("RorRM8CL"), a, reg.CL) }
func (s *Section) RorRM32CL(a operand.RM32)            { s.inst(form("RorRM32CL"), a, reg.CL) }
func (s *Section) RorRM8Imm8(a operand.RM8, v int64)   { s.inst(form("RorRM8Imm8"), a, imm(v)) }
func (s *Section) RorRM32Imm8(a operand.RM32, v int64) { s.inst(form("RorRM32Imm8"), a, imm(v)) }

func (s *Section) RclRM8One(a operand.RM8)             { s.inst(form("RclRM8One"), a, imm(1)) }
func (s *Section) RclRM32One(a operand.RM32)           { s.inst(form("RclRM32One"), a, imm(1)) }
func (s *Section) RclRM8CL(a operand.RM8)              { s.inst(form("RclRM8CL"), a, reg.CL) }
func (s *Section) RclRM32CL(a operand.RM32)            { s.inst(form("RclRM32CL"), a, reg.CL) }
func (s *Section) RclRM8Imm8(a operand.RM8, v int64)   { s.inst(form("RclRM8Imm8"), a, imm(v)) }
func (s *Section) RclRM32Imm8(a operand.RM32, v int64) { s.inst(form("RclRM32Imm8"), a, imm(v)) }

func (s *Section) RcrRM8One(a operand.RM8)             { s.inst(form("RcrRM8One"), a, imm(1)) }
func (s *Section) RcrRM32One(a operand.RM32)           { s.inst(form("RcrRM32One"), a, imm(1)) }
func (s *Section) RcrRM8CL(a operand.RM8)              { s.inst(form("RcrRM8CL"), a, reg.CL) }
func (s *Section) RcrRM32CL(a operand.RM32)            { s.inst(form("RcrRM32CL"), a, reg.CL) }
func (s *Section) RcrRM8Imm8(a operand.RM8, v int64)   { s.inst(form("RcrRM8Imm8"), a, imm(v)) }
func (s *Section) RcrRM32Imm8(a operand.RM32, v int64) { s.inst(form("RcrRM32Imm8"), a, imm(v)) }

func (s *Section) ShlRM8One(a operand.RM8)             { s.inst(form("ShlRM8One"), a, imm(1)) }
func (s *Section) ShlRM32One(a operand.RM32)           { s.inst(form("ShlRM32One"), a, imm(1)) }
func (s *Section) ShlRM8CL(a operand.RM8)              { s.inst(form("ShlRM8CL"), a, reg.CL) }
func (s *Section) ShlRM32CL(a operand.RM32)            { s.inst(form("ShlRM32CL"), a, reg.CL) }
func (s *Section) ShlRM8Imm8(a operand.RM8, v int64)   { s.inst(form("ShlRM8Imm8"), a, imm(v)) }
func (s *Section) ShlRM32Imm8(a operand.RM32, v int64) { s.inst(form("ShlRM32Imm8"), a, imm(v)) }

func (s *Section) ShrRM8One(a operand.RM8)             { s.inst(form("ShrRM8One"), a, imm(1)) }
func (s *Section) ShrRM32One(a operand.RM32)           { s.inst(form("ShrRM32One"), a, imm(1)) }
func (s *Section) ShrRM8CL(a operand.RM8)              { s.inst(form("ShrRM8CL"), a, reg.CL) }
func (s *Section) ShrRM32CL(a operand.RM32)            { s.inst(form("ShrRM32CL"), a, reg.CL) }
func (s *Section) ShrRM8Imm8(a operand.RM8, v int64)   { s.inst(form("ShrRM8Imm8"), a, imm(v)) }
func (s *Section) ShrRM32Imm8(a operand.RM32, v int64) { s.inst(form("ShrRM32Imm8"), a, imm(v)) }

func (s *Section) SarRM8One(a operand.RM8)             { s.inst(form("SarRM8One"), a, imm(1)) }
func (s *Section) SarRM32One(a operand.RM32)           { s.inst(form("SarRM32One"), a, imm(1)) }
func (s *Section) SarRM8CL(a operand.RM8)              { s.inst(form("SarRM8CL"), a, reg.CL) }
func (s *Section) SarRM32CL(a operand.RM32)            { s.inst(form("SarRM32CL"), a, reg.CL) }
func (s *Section) SarRM8Imm8(a operand.RM8, v int64)   { s.inst(form("SarRM8Imm8"), a, imm(v)) }
func (s *Section) SarRM32Imm8(a operand.RM32, v int64) { s.inst(form("SarRM32Imm8"), a, imm(v)) }

// SAL is a documented alias of SHL; the encoder emits SHL bytes either way.
func (s *Section) SalRM8One(a operand.RM8)             { s.inst(form("SalRM8One"), a, imm(1)) }
func (s *Section) SalRM32One(a operand.RM32)           { s.inst(form("SalRM32One"), a, imm(1)) }
func (s *Section) SalRM8CL(a operand.RM8)              { s.inst(form("SalRM8CL"), a, reg.CL) }
func (s *Section) SalRM32CL(a operand.RM32)            { s.inst(form("SalRM32CL"), a, reg.CL) }
func (s *Section) SalRM8Imm8(a operand.RM8, v int64)   { s.inst(form("SalRM8Imm8"), a, imm(v)) }
func (s *Section) SalRM32Imm8(a operand.RM32, v int64) { s.inst(form("SalRM32Imm8"), a, imm(v)) }

func (s *Section) ShldRM32R32Imm8(a operand.RM32, b reg.R32, v int64) {
	s.inst(form("ShldRM32R32Imm8"), a, b, imm(v))
}
func (s *Section) ShldRM32R32CL(a operand.RM32, b reg.R32) { s.inst(form("ShldRM32R32CL"), a, b, reg.CL) }
func (s *Section) ShrdRM32R32Imm8(a operand.RM32, b reg.R32, v int64) {
	s.inst(form("ShrdRM32R32Imm8"), a, b, imm(v))
}
func (s *Section) ShrdRM32R32CL(a operand.RM32, b reg.R32) { s.inst(form("ShrdRM32R32CL"), a, b, reg.CL) }

// ── Branches. The Label form targets this section; the Ref form leaves it. ──
// Short pins rel8 and can fail at Finalize with ErrRange; the plain name pins
// rel32. There is no relaxation between them.

func (s *Section) JmpShortLabel(name string)   { s.inst(form("JmpRel8"), lbl(name)) }
func (s *Section) JmpShortRef(r operand.SymRef) { s.inst(form("JmpRel8"), r) }
func (s *Section) JmpLabel(name string)        { s.inst(form("JmpRel32"), lbl(name)) }
func (s *Section) JmpRef(r operand.SymRef)     { s.inst(form("JmpRel32"), r) }
func (s *Section) JmpRM32(a operand.RM32)      { s.inst(form("JmpRM32"), a) }

func (s *Section) CallLabel(name string)    { s.inst(form("CallRel32"), lbl(name)) }
func (s *Section) CallRef(r operand.SymRef) { s.inst(form("CallRel32"), r) }
func (s *Section) CallRM32(a operand.RM32)  { s.inst(form("CallRM32"), a) }

func (s *Section) Ret()               { s.inst(form("Ret")) }
func (s *Section) RetImm16(v int64)   { s.inst(form("RetImm16"), imm(v)) }

func (s *Section) LoopLabel(name string)     { s.inst(form("LoopRel8"), lbl(name)) }
func (s *Section) LoopeLabel(name string)    { s.inst(form("LoopeRel8"), lbl(name)) }
func (s *Section) LoopneLabel(name string)   { s.inst(form("LoopneRel8"), lbl(name)) }
func (s *Section) JecxzLabel(name string)    { s.inst(form("JecxzRel8"), lbl(name)) }

// ── Jcc / SETcc / CMOVcc — all 30 documented condition spellings ────────────

func (s *Section) JoShortLabel(name string)  { s.inst(form("JoRel8"), lbl(name)) }
func (s *Section) JoLabel(name string)       { s.inst(form("JoRel32"), lbl(name)) }
func (s *Section) JoRef(r operand.SymRef)    { s.inst(form("JoRel32"), r) }
func (s *Section) JnoShortLabel(name string) { s.inst(form("JnoRel8"), lbl(name)) }
func (s *Section) JnoLabel(name string)      { s.inst(form("JnoRel32"), lbl(name)) }
func (s *Section) JnoRef(r operand.SymRef)   { s.inst(form("JnoRel32"), r) }
func (s *Section) JbShortLabel(name string)  { s.inst(form("JbRel8"), lbl(name)) }
func (s *Section) JbLabel(name string)       { s.inst(form("JbRel32"), lbl(name)) }
func (s *Section) JbRef(r operand.SymRef)    { s.inst(form("JbRel32"), r) }
func (s *Section) JcShortLabel(name string)  { s.inst(form("JcRel8"), lbl(name)) }
func (s *Section) JcLabel(name string)       { s.inst(form("JcRel32"), lbl(name)) }
func (s *Section) JcRef(r operand.SymRef)    { s.inst(form("JcRel32"), r) }
func (s *Section) JnaeShortLabel(name string) { s.inst(form("JnaeRel8"), lbl(name)) }
func (s *Section) JnaeLabel(name string)      { s.inst(form("JnaeRel32"), lbl(name)) }
func (s *Section) JnaeRef(r operand.SymRef)   { s.inst(form("JnaeRel32"), r) }
func (s *Section) JaeShortLabel(name string) { s.inst(form("JaeRel8"), lbl(name)) }
func (s *Section) JaeLabel(name string)      { s.inst(form("JaeRel32"), lbl(name)) }
func (s *Section) JaeRef(r operand.SymRef)   { s.inst(form("JaeRel32"), r) }
func (s *Section) JnbShortLabel(name string) { s.inst(form("JnbRel8"), lbl(name)) }
func (s *Section) JnbLabel(name string)      { s.inst(form("JnbRel32"), lbl(name)) }
func (s *Section) JnbRef(r operand.SymRef)   { s.inst(form("JnbRel32"), r) }
func (s *Section) JncShortLabel(name string) { s.inst(form("JncRel8"), lbl(name)) }
func (s *Section) JncLabel(name string)      { s.inst(form("JncRel32"), lbl(name)) }
func (s *Section) JncRef(r operand.SymRef)   { s.inst(form("JncRel32"), r) }
func (s *Section) JeShortLabel(name string)  { s.inst(form("JeRel8"), lbl(name)) }
func (s *Section) JeLabel(name string)       { s.inst(form("JeRel32"), lbl(name)) }
func (s *Section) JeRef(r operand.SymRef)    { s.inst(form("JeRel32"), r) }
func (s *Section) JzShortLabel(name string)  { s.inst(form("JzRel8"), lbl(name)) }
func (s *Section) JzLabel(name string)       { s.inst(form("JzRel32"), lbl(name)) }
func (s *Section) JzRef(r operand.SymRef)    { s.inst(form("JzRel32"), r) }
func (s *Section) JneShortLabel(name string) { s.inst(form("JneRel8"), lbl(name)) }
func (s *Section) JneLabel(name string)      { s.inst(form("JneRel32"), lbl(name)) }
func (s *Section) JneRef(r operand.SymRef)   { s.inst(form("JneRel32"), r) }
func (s *Section) JnzShortLabel(name string) { s.inst(form("JnzRel8"), lbl(name)) }
func (s *Section) JnzLabel(name string)      { s.inst(form("JnzRel32"), lbl(name)) }
func (s *Section) JnzRef(r operand.SymRef)   { s.inst(form("JnzRel32"), r) }
func (s *Section) JbeShortLabel(name string) { s.inst(form("JbeRel8"), lbl(name)) }
func (s *Section) JbeLabel(name string)      { s.inst(form("JbeRel32"), lbl(name)) }
func (s *Section) JbeRef(r operand.SymRef)   { s.inst(form("JbeRel32"), r) }
func (s *Section) JnaShortLabel(name string) { s.inst(form("JnaRel8"), lbl(name)) }
func (s *Section) JnaLabel(name string)      { s.inst(form("JnaRel32"), lbl(name)) }
func (s *Section) JnaRef(r operand.SymRef)   { s.inst(form("JnaRel32"), r) }
func (s *Section) JaShortLabel(name string)  { s.inst(form("JaRel8"), lbl(name)) }
func (s *Section) JaLabel(name string)       { s.inst(form("JaRel32"), lbl(name)) }
func (s *Section) JaRef(r operand.SymRef)    { s.inst(form("JaRel32"), r) }
func (s *Section) JnbeShortLabel(name string) { s.inst(form("JnbeRel8"), lbl(name)) }
func (s *Section) JnbeLabel(name string)      { s.inst(form("JnbeRel32"), lbl(name)) }
func (s *Section) JnbeRef(r operand.SymRef)   { s.inst(form("JnbeRel32"), r) }
func (s *Section) JsShortLabel(name string)  { s.inst(form("JsRel8"), lbl(name)) }
func (s *Section) JsLabel(name string)       { s.inst(form("JsRel32"), lbl(name)) }
func (s *Section) JsRef(r operand.SymRef)    { s.inst(form("JsRel32"), r) }
func (s *Section) JnsShortLabel(name string) { s.inst(form("JnsRel8"), lbl(name)) }
func (s *Section) JnsLabel(name string)      { s.inst(form("JnsRel32"), lbl(name)) }
func (s *Section) JnsRef(r operand.SymRef)   { s.inst(form("JnsRel32"), r) }
func (s *Section) JpShortLabel(name string)  { s.inst(form("JpRel8"), lbl(name)) }
func (s *Section) JpLabel(name string)       { s.inst(form("JpRel32"), lbl(name)) }
func (s *Section) JpRef(r operand.SymRef)    { s.inst(form("JpRel32"), r) }
func (s *Section) JpeShortLabel(name string) { s.inst(form("JpeRel8"), lbl(name)) }
func (s *Section) JpeLabel(name string)      { s.inst(form("JpeRel32"), lbl(name)) }
func (s *Section) JpeRef(r operand.SymRef)   { s.inst(form("JpeRel32"), r) }
func (s *Section) JnpShortLabel(name string) { s.inst(form("JnpRel8"), lbl(name)) }
func (s *Section) JnpLabel(name string)      { s.inst(form("JnpRel32"), lbl(name)) }
func (s *Section) JnpRef(r operand.SymRef)   { s.inst(form("JnpRel32"), r) }
func (s *Section) JpoShortLabel(name string) { s.inst(form("JpoRel8"), lbl(name)) }
func (s *Section) JpoLabel(name string)      { s.inst(form("JpoRel32"), lbl(name)) }
func (s *Section) JpoRef(r operand.SymRef)   { s.inst(form("JpoRel32"), r) }
func (s *Section) JlShortLabel(name string)  { s.inst(form("JlRel8"), lbl(name)) }
func (s *Section) JlLabel(name string)       { s.inst(form("JlRel32"), lbl(name)) }
func (s *Section) JlRef(r operand.SymRef)    { s.inst(form("JlRel32"), r) }
func (s *Section) JngeShortLabel(name string) { s.inst(form("JngeRel8"), lbl(name)) }
func (s *Section) JngeLabel(name string)      { s.inst(form("JngeRel32"), lbl(name)) }
func (s *Section) JngeRef(r operand.SymRef)   { s.inst(form("JngeRel32"), r) }
func (s *Section) JgeShortLabel(name string) { s.inst(form("JgeRel8"), lbl(name)) }
func (s *Section) JgeLabel(name string)      { s.inst(form("JgeRel32"), lbl(name)) }
func (s *Section) JgeRef(r operand.SymRef)   { s.inst(form("JgeRel32"), r) }
func (s *Section) JnlShortLabel(name string) { s.inst(form("JnlRel8"), lbl(name)) }
func (s *Section) JnlLabel(name string)      { s.inst(form("JnlRel32"), lbl(name)) }
func (s *Section) JnlRef(r operand.SymRef)   { s.inst(form("JnlRel32"), r) }
func (s *Section) JleShortLabel(name string) { s.inst(form("JleRel8"), lbl(name)) }
func (s *Section) JleLabel(name string)      { s.inst(form("JleRel32"), lbl(name)) }
func (s *Section) JleRef(r operand.SymRef)   { s.inst(form("JleRel32"), r) }
func (s *Section) JngShortLabel(name string) { s.inst(form("JngRel8"), lbl(name)) }
func (s *Section) JngLabel(name string)      { s.inst(form("JngRel32"), lbl(name)) }
func (s *Section) JngRef(r operand.SymRef)   { s.inst(form("JngRel32"), r) }
func (s *Section) JgShortLabel(name string)  { s.inst(form("JgRel8"), lbl(name)) }
func (s *Section) JgLabel(name string)       { s.inst(form("JgRel32"), lbl(name)) }
func (s *Section) JgRef(r operand.SymRef)    { s.inst(form("JgRel32"), r) }
func (s *Section) JnleShortLabel(name string) { s.inst(form("JnleRel8"), lbl(name)) }
func (s *Section) JnleLabel(name string)      { s.inst(form("JnleRel32"), lbl(name)) }
func (s *Section) JnleRef(r operand.SymRef)   { s.inst(form("JnleRel32"), r) }

func (s *Section) SetoRM8(a operand.RM8)   { s.inst(form("SetoRM8"), a) }
func (s *Section) SetnoRM8(a operand.RM8)  { s.inst(form("SetnoRM8"), a) }
func (s *Section) SetbRM8(a operand.RM8)   { s.inst(form("SetbRM8"), a) }
func (s *Section) SetcRM8(a operand.RM8)   { s.inst(form("SetcRM8"), a) }
func (s *Section) SetnaeRM8(a operand.RM8) { s.inst(form("SetnaeRM8"), a) }
func (s *Section) SetaeRM8(a operand.RM8)  { s.inst(form("SetaeRM8"), a) }
func (s *Section) SetnbRM8(a operand.RM8)  { s.inst(form("SetnbRM8"), a) }
func (s *Section) SetncRM8(a operand.RM8)  { s.inst(form("SetncRM8"), a) }
func (s *Section) SeteRM8(a operand.RM8)   { s.inst(form("SeteRM8"), a) }
func (s *Section) SetzRM8(a operand.RM8)   { s.inst(form("SetzRM8"), a) }
func (s *Section) SetneRM8(a operand.RM8)  { s.inst(form("SetneRM8"), a) }
func (s *Section) SetnzRM8(a operand.RM8)  { s.inst(form("SetnzRM8"), a) }
func (s *Section) SetbeRM8(a operand.RM8)  { s.inst(form("SetbeRM8"), a) }
func (s *Section) SetnaRM8(a operand.RM8)  { s.inst(form("SetnaRM8"), a) }
func (s *Section) SetaRM8(a operand.RM8)   { s.inst(form("SetaRM8"), a) }
func (s *Section) SetnbeRM8(a operand.RM8) { s.inst(form("SetnbeRM8"), a) }
func (s *Section) SetsRM8(a operand.RM8)   { s.inst(form("SetsRM8"), a) }
func (s *Section) SetnsRM8(a operand.RM8)  { s.inst(form("SetnsRM8"), a) }
func (s *Section) SetpRM8(a operand.RM8)   { s.inst(form("SetpRM8"), a) }
func (s *Section) SetpeRM8(a operand.RM8)  { s.inst(form("SetpeRM8"), a) }
func (s *Section) SetnpRM8(a operand.RM8)  { s.inst(form("SetnpRM8"), a) }
func (s *Section) SetpoRM8(a operand.RM8)  { s.inst(form("SetpoRM8"), a) }
func (s *Section) SetlRM8(a operand.RM8)   { s.inst(form("SetlRM8"), a) }
func (s *Section) SetngeRM8(a operand.RM8) { s.inst(form("SetngeRM8"), a) }
func (s *Section) SetgeRM8(a operand.RM8)  { s.inst(form("SetgeRM8"), a) }
func (s *Section) SetnlRM8(a operand.RM8)  { s.inst(form("SetnlRM8"), a) }
func (s *Section) SetleRM8(a operand.RM8)  { s.inst(form("SetleRM8"), a) }
func (s *Section) SetngRM8(a operand.RM8)  { s.inst(form("SetngRM8"), a) }
func (s *Section) SetgRM8(a operand.RM8)   { s.inst(form("SetgRM8"), a) }
func (s *Section) SetnleRM8(a operand.RM8) { s.inst(form("SetnleRM8"), a) }

func (s *Section) CmovoR32RM32(a reg.R32, b operand.RM32)   { s.inst(form("CmovoR32RM32"), a, b) }
func (s *Section) CmovnoR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovnoR32RM32"), a, b) }
func (s *Section) CmovbR32RM32(a reg.R32, b operand.RM32)   { s.inst(form("CmovbR32RM32"), a, b) }
func (s *Section) CmovcR32RM32(a reg.R32, b operand.RM32)   { s.inst(form("CmovcR32RM32"), a, b) }
func (s *Section) CmovnaeR32RM32(a reg.R32, b operand.RM32) { s.inst(form("CmovnaeR32RM32"), a, b) }
func (s *Section) CmovaeR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovaeR32RM32"), a, b) }
func (s *Section) CmovnbR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovnbR32RM32"), a, b) }
func (s *Section) CmovncR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovncR32RM32"), a, b) }
func (s *Section) CmoveR32RM32(a reg.R32, b operand.RM32)   { s.inst(form("CmoveR32RM32"), a, b) }
func (s *Section) CmovzR32RM32(a reg.R32, b operand.RM32)   { s.inst(form("CmovzR32RM32"), a, b) }
func (s *Section) CmovneR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovneR32RM32"), a, b) }
func (s *Section) CmovnzR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovnzR32RM32"), a, b) }
func (s *Section) CmovbeR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovbeR32RM32"), a, b) }
func (s *Section) CmovnaR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovnaR32RM32"), a, b) }
func (s *Section) CmovaR32RM32(a reg.R32, b operand.RM32)   { s.inst(form("CmovaR32RM32"), a, b) }
func (s *Section) CmovnbeR32RM32(a reg.R32, b operand.RM32) { s.inst(form("CmovnbeR32RM32"), a, b) }
func (s *Section) CmovsR32RM32(a reg.R32, b operand.RM32)   { s.inst(form("CmovsR32RM32"), a, b) }
func (s *Section) CmovnsR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovnsR32RM32"), a, b) }
func (s *Section) CmovpR32RM32(a reg.R32, b operand.RM32)   { s.inst(form("CmovpR32RM32"), a, b) }
func (s *Section) CmovpeR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovpeR32RM32"), a, b) }
func (s *Section) CmovnpR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovnpR32RM32"), a, b) }
func (s *Section) CmovpoR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovpoR32RM32"), a, b) }
func (s *Section) CmovlR32RM32(a reg.R32, b operand.RM32)   { s.inst(form("CmovlR32RM32"), a, b) }
func (s *Section) CmovngeR32RM32(a reg.R32, b operand.RM32) { s.inst(form("CmovngeR32RM32"), a, b) }
func (s *Section) CmovgeR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovgeR32RM32"), a, b) }
func (s *Section) CmovnlR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovnlR32RM32"), a, b) }
func (s *Section) CmovleR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovleR32RM32"), a, b) }
func (s *Section) CmovngR32RM32(a reg.R32, b operand.RM32)  { s.inst(form("CmovngR32RM32"), a, b) }
func (s *Section) CmovgR32RM32(a reg.R32, b operand.RM32)   { s.inst(form("CmovgR32RM32"), a, b) }
func (s *Section) CmovnleR32RM32(a reg.R32, b operand.RM32) { s.inst(form("CmovnleR32RM32"), a, b) }

// ── Misc, bit ops, I/O ───────────────────────────────────────────────────────

func (s *Section) Nop()             { s.inst(form("Nop")) }
func (s *Section) Int3()            { s.inst(form("Int3")) }
func (s *Section) IntImm8(v int64)  { s.inst(form("IntImm8"), imm(v)) }
func (s *Section) Into()            { s.inst(form("Into")) }
func (s *Section) Iret()            { s.inst(form("Iret")) }
func (s *Section) Hlt()             { s.inst(form("Hlt")) }
func (s *Section) Cld()             { s.inst(form("Cld")) }
func (s *Section) Std()             { s.inst(form("Std")) }
func (s *Section) Clc()             { s.inst(form("Clc")) }
func (s *Section) Stc()             { s.inst(form("Stc")) }
func (s *Section) Cmc()             { s.inst(form("Cmc")) }
func (s *Section) Lahf()            { s.inst(form("Lahf")) }
func (s *Section) Sahf()            { s.inst(form("Sahf")) }

func (s *Section) BtRM32R32(a operand.RM32, b reg.R32)  { s.inst(form("BtRM32R32"), a, b) }
func (s *Section) BtsRM32R32(a operand.RM32, b reg.R32) { s.inst(form("BtsRM32R32"), a, b) }
func (s *Section) BtrRM32R32(a operand.RM32, b reg.R32) { s.inst(form("BtrRM32R32"), a, b) }
func (s *Section) BtcRM32R32(a operand.RM32, b reg.R32) { s.inst(form("BtcRM32R32"), a, b) }
func (s *Section) BsfR32RM32(a reg.R32, b operand.RM32) { s.inst(form("BsfR32RM32"), a, b) }
func (s *Section) BsrR32RM32(a reg.R32, b operand.RM32) { s.inst(form("BsrR32RM32"), a, b) }

func (s *Section) InALImm8(v int64)  { s.inst(form("InALImm8"), reg.AL, imm(v)) }
func (s *Section) InALDX()           { s.inst(form("InALDX"), reg.AL, reg.DX) }
func (s *Section) OutImm8AL(v int64) { s.inst(form("OutImm8AL"), imm(v), reg.AL) }
func (s *Section) OutDXAL()          { s.inst(form("OutDXAL"), reg.DX, reg.AL) }

// ── System ───────────────────────────────────────────────────────────────────

func (s *Section) MovR32Cr(a reg.R32, b reg.Cr) { s.inst(form("MovR32Cr"), a, b) }
func (s *Section) MovCrR32(a reg.Cr, b reg.R32) { s.inst(form("MovCrR32"), a, b) }
func (s *Section) MovR32Dr(a reg.R32, b reg.Dr) { s.inst(form("MovR32Dr"), a, b) }
func (s *Section) MovDrR32(a reg.Dr, b reg.R32) { s.inst(form("MovDrR32"), a, b) }

func (s *Section) Invd()                    { s.inst(form("Invd")) }
func (s *Section) Wbinvd()                  { s.inst(form("Wbinvd")) }
func (s *Section) InvlpgM(a operand.Memory) { s.inst(form("InvlpgM"), a) }
func (s *Section) Cpuid()                   { s.inst(form("Cpuid")) }
func (s *Section) Rdtsc()                   { s.inst(form("Rdtsc")) }
func (s *Section) Rdmsr()                   { s.inst(form("Rdmsr")) }
func (s *Section) Wrmsr()                   { s.inst(form("Wrmsr")) }
func (s *Section) Rdpmc()                   { s.inst(form("Rdpmc")) }