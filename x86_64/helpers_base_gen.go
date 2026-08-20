// x86_64/helpers_base_gen.go
//
// Code generated from the isa tables by helpergen. DO NOT EDIT.
//
// One method per declared form, named by Form.GoName: the mnemonic, then
// each explicit slot's class. A helper pins its form — MovR64Imm64 is the
// ten-byte imm64 and nothing quietly relaxes it. Fixed operands are in the
// name, not the parameters: AddRAXImm32 leaves no field to put another
// register in, and passes reg.RAX itself so the class check holds.
package x86_64

import (
	"github.com/vertex-language/asm/x86_64/operand"
	"github.com/vertex-language/asm/x86_64/reg"
)

// ---- mov ------------------------------------------------------------

func (s *Section) MovRM8R8(dst RM8, src reg.Reg8)     { s.place(form("MovRM8R8"), noOpts, dst, src) }
func (s *Section) MovRM32R32(dst RM32, src reg.Reg32) { s.place(form("MovRM32R32"), noOpts, dst, src) }
func (s *Section) MovRM64R64(dst RM64, src reg.Reg64) { s.place(form("MovRM64R64"), noOpts, dst, src) }
func (s *Section) MovR8RM8(dst reg.Reg8, src RM8)     { s.place(form("MovR8RM8"), noOpts, dst, src) }
func (s *Section) MovR32RM32(dst reg.Reg32, src RM32) { s.place(form("MovR32RM32"), noOpts, dst, src) }
func (s *Section) MovR64RM64(dst reg.Reg64, src RM64) { s.place(form("MovR64RM64"), noOpts, dst, src) }

func (s *Section) MovR8Imm8(dst reg.Reg8, v int64)   { s.place(form("MovR8Imm8"), noOpts, dst, operand.Imm(v)) }
func (s *Section) MovR32Imm32(dst reg.Reg32, v int64) { s.place(form("MovR32Imm32"), noOpts, dst, operand.Imm(v)) }
func (s *Section) MovR64Imm64(dst reg.Reg64, v int64) { s.place(form("MovR64Imm64"), noOpts, dst, operand.Imm(v)) }
func (s *Section) MovRM8Imm8(dst RM8, v int64)       { s.place(form("MovRM8Imm8"), noOpts, dst, operand.Imm(v)) }
func (s *Section) MovRM32Imm32(dst RM32, v int64)    { s.place(form("MovRM32Imm32"), noOpts, dst, operand.Imm(v)) }
func (s *Section) MovRM64Imm32(dst RM64, v int64)    { s.place(form("MovRM64Imm32"), noOpts, dst, operand.Imm(v)) }

// ---- lea, widening moves ---------------------------------------------

func (s *Section) LeaR32M(dst reg.Reg32, m Memory) { s.place(form("LeaR32M"), noOpts, dst, m) }
func (s *Section) LeaR64M(dst reg.Reg64, m Memory) { s.place(form("LeaR64M"), noOpts, dst, m) }

func (s *Section) MovzxR32RM8(dst reg.Reg32, src RM8)   { s.place(form("MovzxR32RM8"), noOpts, dst, src) }
func (s *Section) MovzxR64RM8(dst reg.Reg64, src RM8)   { s.place(form("MovzxR64RM8"), noOpts, dst, src) }
func (s *Section) MovzxR32RM16(dst reg.Reg32, src RM16) { s.place(form("MovzxR32RM16"), noOpts, dst, src) }
func (s *Section) MovsxR64RM8(dst reg.Reg64, src RM8)   { s.place(form("MovsxR64RM8"), noOpts, dst, src) }
func (s *Section) MovsxdR64RM32(dst reg.Reg64, src RM32) { s.place(form("MovsxdR64RM32"), noOpts, dst, src) }

// ---- the ALU group ------------------------------------------------------

func (s *Section) AddRM64R64(dst RM64, src reg.Reg64) { s.place(form("AddRM64R64"), noOpts, dst, src) }
func (s *Section) AddR64RM64(dst reg.Reg64, src RM64) { s.place(form("AddR64RM64"), noOpts, dst, src) }
func (s *Section) AddR32RM32(dst reg.Reg32, src RM32) { s.place(form("AddR32RM32"), noOpts, dst, src) }
func (s *Section) AddRAXImm32(v int64)                { s.place(form("AddRAXImm32"), noOpts, reg.RAX, operand.Imm(v)) }
func (s *Section) AddRM64Imm8(dst RM64, v int64)      { s.place(form("AddRM64Imm8"), noOpts, dst, operand.Imm(v)) }
func (s *Section) AddRM64Imm32(dst RM64, v int64)     { s.place(form("AddRM64Imm32"), noOpts, dst, operand.Imm(v)) }

func (s *Section) SubRM64R64(dst RM64, src reg.Reg64) { s.place(form("SubRM64R64"), noOpts, dst, src) }
func (s *Section) SubR64RM64(dst reg.Reg64, src RM64) { s.place(form("SubR64RM64"), noOpts, dst, src) }
func (s *Section) SubRM64Imm8(dst RM64, v int64)      { s.place(form("SubRM64Imm8"), noOpts, dst, operand.Imm(v)) }
func (s *Section) SubRM64Imm32(dst RM64, v int64)     { s.place(form("SubRM64Imm32"), noOpts, dst, operand.Imm(v)) }

func (s *Section) XorR32RM32(dst reg.Reg32, src RM32) { s.place(form("XorR32RM32"), noOpts, dst, src) }
func (s *Section) XorR64RM64(dst reg.Reg64, src RM64) { s.place(form("XorR64RM64"), noOpts, dst, src) }
func (s *Section) AndRM64R64(dst RM64, src reg.Reg64) { s.place(form("AndRM64R64"), noOpts, dst, src) }
func (s *Section) OrRM64R64(dst RM64, src reg.Reg64)  { s.place(form("OrRM64R64"), noOpts, dst, src) }

func (s *Section) CmpRM64R64(dst RM64, src reg.Reg64) { s.place(form("CmpRM64R64"), noOpts, dst, src) }
func (s *Section) CmpR64RM64(dst reg.Reg64, src RM64) { s.place(form("CmpR64RM64"), noOpts, dst, src) }
func (s *Section) CmpRM64Imm8(dst RM64, v int64)      { s.place(form("CmpRM64Imm8"), noOpts, dst, operand.Imm(v)) }
func (s *Section) CmpRM64Imm32(dst RM64, v int64)     { s.place(form("CmpRM64Imm32"), noOpts, dst, operand.Imm(v)) }

// ---- shifts ---------------------------------------------------------

func (s *Section) ShlRM64One(dst RM64)           { s.place(form("ShlRM64One"), noOpts, dst, operand.Imm(1)) }
func (s *Section) ShlRM64CL(dst RM64)            { s.place(form("ShlRM64CL"), noOpts, dst, reg.CL) }
func (s *Section) ShlRM64Imm8(dst RM64, v int64) { s.place(form("ShlRM64Imm8"), noOpts, dst, operand.Imm(v)) }
func (s *Section) ShrRM64Imm8(dst RM64, v int64) { s.place(form("ShrRM64Imm8"), noOpts, dst, operand.Imm(v)) }
func (s *Section) SarRM64Imm8(dst RM64, v int64) { s.place(form("SarRM64Imm8"), noOpts, dst, operand.Imm(v)) }

// ---- test, unary, imul ------------------------------------------------

func (s *Section) TestRM64R64(a RM64, b reg.Reg64) { s.place(form("TestRM64R64"), noOpts, a, b) }
func (s *Section) TestRAXImm32(v int64)            { s.place(form("TestRAXImm32"), noOpts, reg.RAX, operand.Imm(v)) }

func (s *Section) IncRM64(dst RM64) { s.place(form("IncRM64"), noOpts, dst) }
func (s *Section) DecRM64(dst RM64) { s.place(form("DecRM64"), noOpts, dst) }
func (s *Section) NegRM64(dst RM64) { s.place(form("NegRM64"), noOpts, dst) }
func (s *Section) NotRM64(dst RM64) { s.place(form("NotRM64"), noOpts, dst) }
func (s *Section) MulRM64(src RM64)  { s.place(form("MulRM64"), noOpts, src) }
func (s *Section) DivRM64(src RM64)  { s.place(form("DivRM64"), noOpts, src) }
func (s *Section) IdivRM64(src RM64) { s.place(form("IdivRM64"), noOpts, src) }

func (s *Section) ImulR64RM64(dst reg.Reg64, src RM64) { s.place(form("ImulR64RM64"), noOpts, dst, src) }
func (s *Section) ImulR64RM64Imm32(dst reg.Reg64, src RM64, v int64) {
	s.place(form("ImulR64RM64Imm32"), noOpts, dst, src, operand.Imm(v))
}

// ---- stack, control ----------------------------------------------------

func (s *Section) PushR64(r reg.Reg64) { s.place(form("PushR64"), noOpts, r) }
func (s *Section) PopR64(r reg.Reg64)  { s.place(form("PopR64"), noOpts, r) }
func (s *Section) PushRM64(v RM64)     { s.place(form("PushRM64"), noOpts, v) }
func (s *Section) PopRM64(v RM64)      { s.place(form("PopRM64"), noOpts, v) }
func (s *Section) PushImm8(v int64)    { s.place(form("PushImm8"), noOpts, operand.Imm(v)) }
func (s *Section) PushImm32(v int64)   { s.place(form("PushImm32"), noOpts, operand.Imm(v)) }
func (s *Section) Leave()              { s.place(form("Leave"), noOpts) }

func (s *Section) CallRM64(t RM64)   { s.place(form("CallRM64"), noOpts, t) }
func (s *Section) JmpRM64(t RM64)    { s.place(form("JmpRM64"), noOpts, t) }
func (s *Section) Ret()              { s.place(form("Ret"), noOpts) }
func (s *Section) RetImm16(v int64)  { s.place(form("RetImm16"), noOpts, operand.Imm(v)) }

func (s *Section) Nop()     { s.place(form("Nop"), noOpts) }
func (s *Section) Int3()    { s.place(form("Int3"), noOpts) }
func (s *Section) Ud2()     { s.place(form("Ud2"), noOpts) }
func (s *Section) Hlt()     { s.place(form("Hlt"), noOpts) }
func (s *Section) Syscall() { s.place(form("Syscall"), noOpts) }
func (s *Section) Cpuid()   { s.place(form("Cpuid"), noOpts) }
func (s *Section) Rdtsc()   { s.place(form("Rdtsc"), noOpts) }
func (s *Section) Pause()   { s.place(form("Pause"), noOpts) }
func (s *Section) Lfence()  { s.place(form("Lfence"), noOpts) }
func (s *Section) Mfence()  { s.place(form("Mfence"), noOpts) }
func (s *Section) Sfence()  { s.place(form("Sfence"), noOpts) }
func (s *Section) Cwde()    { s.place(form("Cwde"), noOpts) }
func (s *Section) Cdqe()    { s.place(form("Cdqe"), noOpts) }
func (s *Section) Cdq()     { s.place(form("Cdq"), noOpts) }
func (s *Section) Cqo()     { s.place(form("Cqo"), noOpts) }

// ---- conditional branches ------------------------------------------------
//
// Canonical condition spellings, as the table declares them, plus the alias
// spellings the world writes: jz and jnz bind the same rows as je and jne.

func (s *Section) JoLabel(l string)  { s.place(form("JoRel32"), noOpts, operand.Label(l)) }
func (s *Section) JnoLabel(l string) { s.place(form("JnoRel32"), noOpts, operand.Label(l)) }
func (s *Section) JbLabel(l string)  { s.place(form("JbRel32"), noOpts, operand.Label(l)) }
func (s *Section) JaeLabel(l string) { s.place(form("JaeRel32"), noOpts, operand.Label(l)) }
func (s *Section) JeLabel(l string)  { s.place(form("JeRel32"), noOpts, operand.Label(l)) }
func (s *Section) JneLabel(l string) { s.place(form("JneRel32"), noOpts, operand.Label(l)) }
func (s *Section) JbeLabel(l string) { s.place(form("JbeRel32"), noOpts, operand.Label(l)) }
func (s *Section) JaLabel(l string)  { s.place(form("JaRel32"), noOpts, operand.Label(l)) }
func (s *Section) JsLabel(l string)  { s.place(form("JsRel32"), noOpts, operand.Label(l)) }
func (s *Section) JnsLabel(l string) { s.place(form("JnsRel32"), noOpts, operand.Label(l)) }
func (s *Section) JpLabel(l string)  { s.place(form("JpRel32"), noOpts, operand.Label(l)) }
func (s *Section) JnpLabel(l string) { s.place(form("JnpRel32"), noOpts, operand.Label(l)) }
func (s *Section) JlLabel(l string)  { s.place(form("JlRel32"), noOpts, operand.Label(l)) }
func (s *Section) JgeLabel(l string) { s.place(form("JgeRel32"), noOpts, operand.Label(l)) }
func (s *Section) JleLabel(l string) { s.place(form("JleRel32"), noOpts, operand.Label(l)) }
func (s *Section) JgLabel(l string)  { s.place(form("JgRel32"), noOpts, operand.Label(l)) }

func (s *Section) JeShortLabel(l string)  { s.place(form("JeRel8"), noOpts, operand.Label(l)) }
func (s *Section) JneShortLabel(l string) { s.place(form("JneRel8"), noOpts, operand.Label(l)) }
func (s *Section) JbShortLabel(l string)  { s.place(form("JbRel8"), noOpts, operand.Label(l)) }
func (s *Section) JaeShortLabel(l string) { s.place(form("JaeRel8"), noOpts, operand.Label(l)) }
func (s *Section) JlShortLabel(l string)  { s.place(form("JlRel8"), noOpts, operand.Label(l)) }
func (s *Section) JgShortLabel(l string)  { s.place(form("JgRel8"), noOpts, operand.Label(l)) }

// Alias spellings: the same encodings under the names the world writes.
func (s *Section) JzLabel(l string)       { s.JeLabel(l) }
func (s *Section) JnzLabel(l string)      { s.JneLabel(l) }
func (s *Section) JzShortLabel(l string)  { s.JeShortLabel(l) }
func (s *Section) JnzShortLabel(l string) { s.JneShortLabel(l) }

func (s *Section) SeteRM8(dst RM8)  { s.place(form("SeteRM8"), noOpts, dst) }
func (s *Section) SetneRM8(dst RM8) { s.place(form("SetneRM8"), noOpts, dst) }

func (s *Section) CmoveR64RM64(dst reg.Reg64, src RM64)  { s.place(form("CmoveR64RM64"), noOpts, dst, src) }
func (s *Section) CmovneR64RM64(dst reg.Reg64, src RM64) { s.place(form("CmovneR64RM64"), noOpts, dst, src) }

// ---- bit ops (some gated) ------------------------------------------------

func (s *Section) BsfR64RM64(dst reg.Reg64, src RM64)    { s.place(form("BsfR64RM64"), noOpts, dst, src) }
func (s *Section) BsrR64RM64(dst reg.Reg64, src RM64)    { s.place(form("BsrR64RM64"), noOpts, dst, src) }
func (s *Section) BswapR64(r reg.Reg64)                  { s.place(form("BswapR64"), noOpts, r) }
func (s *Section) PopcntR64RM64(dst reg.Reg64, src RM64) { s.place(form("PopcntR64RM64"), noOpts, dst, src) }
func (s *Section) LzcntR64RM64(dst reg.Reg64, src RM64)  { s.place(form("LzcntR64RM64"), noOpts, dst, src) }
func (s *Section) TzcntR64RM64(dst reg.Reg64, src RM64)  { s.place(form("TzcntR64RM64"), noOpts, dst, src) }

// ---- SSE ------------------------------------------------------------

func (s *Section) MovapsXmmXmmM128(dst reg.Xmm, src XmmM128) { s.place(form("MovapsXmmXmmM128"), noOpts, dst, src) }
func (s *Section) MovapsXmmM128Xmm(dst XmmM128, src reg.Xmm) { s.place(form("MovapsXmmM128Xmm"), noOpts, dst, src) }
func (s *Section) MovdqaXmmXmmM128(dst reg.Xmm, src XmmM128) { s.place(form("MovdqaXmmXmmM128"), noOpts, dst, src) }
func (s *Section) MovdquXmmXmmM128(dst reg.Xmm, src XmmM128) { s.place(form("MovdquXmmXmmM128"), noOpts, dst, src) }
func (s *Section) MovqXmmRM64(dst reg.Xmm, src RM64)         { s.place(form("MovqXmmRM64"), noOpts, dst, src) }
func (s *Section) MovqRM64Xmm(dst RM64, src reg.Xmm)         { s.place(form("MovqRM64Xmm"), noOpts, dst, src) }
func (s *Section) AddpsXmmXmmM128(dst reg.Xmm, src XmmM128)  { s.place(form("AddpsXmmXmmM128"), noOpts, dst, src) }
func (s *Section) AddsdXmmXmmM64(dst reg.Xmm, src XmmM64)    { s.place(form("AddsdXmmXmmM64"), noOpts, dst, src) }
func (s *Section) PxorXmmXmmM128(dst reg.Xmm, src XmmM128)   { s.place(form("PxorXmmXmmM128"), noOpts, dst, src) }
func (s *Section) PadddXmmXmmM128(dst reg.Xmm, src XmmM128)  { s.place(form("PadddXmmXmmM128"), noOpts, dst, src) }
func (s *Section) UcomisdXmmXmmM64(a reg.Xmm, b XmmM64)      { s.place(form("UcomisdXmmXmmM64"), noOpts, a, b) }

// ---- AVX ------------------------------------------------------------

func (s *Section) VmovupsYmmYmmM256(dst reg.Ymm, src YmmM256) { s.place(form("VmovupsYmmYmmM256"), noOpts, dst, src) }
func (s *Section) VmovapsYmmYmmM256(dst reg.Ymm, src YmmM256) { s.place(form("VmovapsYmmYmmM256"), noOpts, dst, src) }
func (s *Section) Vzeroupper()                                { s.place(form("Vzeroupper"), noOpts) }
func (s *Section) VaddpsYmmYmmYmmM256(dst reg.Ymm, a reg.Ymm, b YmmM256) {
	s.place(form("VaddpsYmmYmmYmmM256"), noOpts, dst, a, b)
}
func (s *Section) VpxorYmmYmmYmmM256(dst reg.Ymm, a reg.Ymm, b YmmM256) {
	s.place(form("VpxorYmmYmmYmmM256"), noOpts, dst, a, b)
}
func (s *Section) VpadddYmmYmmYmmM256(dst reg.Ymm, a reg.Ymm, b YmmM256) {
	s.place(form("VpadddYmmYmmYmmM256"), noOpts, dst, a, b)
}

// ---- EVEX: masks as operands, modifiers as options ------------------------

func (s *Section) VaddpsZmmKZmmZmmM512(dst reg.Zmm, k reg.K, a reg.Zmm, b ZmmM512, opts ...Opt) {
	s.place(form("VaddpsZmmKZmmZmmM512"), optset(opts), dst, k, a, b)
}
func (s *Section) VpadddZmmKZmmZmmM512(dst reg.Zmm, k reg.K, a reg.Zmm, b ZmmM512, opts ...Opt) {
	s.place(form("VpadddZmmKZmmZmmM512"), optset(opts), dst, k, a, b)
}
func (s *Section) Vmovdqa32ZmmKZmmM512(dst reg.Zmm, k reg.K, src ZmmM512, opts ...Opt) {
	s.place(form("Vmovdqa32ZmmKZmmM512"), optset(opts), dst, k, src)
}
func (s *Section) KmovwKKM64(dst reg.K, src KM64) { s.place(form("KmovwKKM64"), noOpts, dst, src) }