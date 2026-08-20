package i386

// Re-exports, so a lowering imports one package and writes i386.EAX,
// i386.Mem32, i386.Ref, i386.Global. The subpackages remain the definitions.

import (
	"github.com/vertex-language/asm/i386/feature"
	"github.com/vertex-language/asm/i386/operand"
	"github.com/vertex-language/asm/i386/reg"
)

// Operand is anything an instruction can take.
type Operand = operand.Operand

// Registers.
const (
	EAX, ECX, EDX, EBX = reg.EAX, reg.ECX, reg.EDX, reg.EBX
	ESP, EBP, ESI, EDI = reg.ESP, reg.EBP, reg.ESI, reg.EDI

	AX, CX, DX, BX = reg.AX, reg.CX, reg.DX, reg.BX
	SP, BP, SI, DI = reg.SP, reg.BP, reg.SI, reg.DI

	AL, CL, DL, BL = reg.AL, reg.CL, reg.DL, reg.BL
	AH, CH, DH, BH = reg.AH, reg.CH, reg.DH, reg.BH

	ES, CS, SS, DS, FS, GS = reg.ES, reg.CS, reg.SS, reg.DS, reg.FS, reg.GS
)

const (
	ST0, ST1, ST2, ST3, ST4, ST5, ST6, ST7 = reg.ST0, reg.ST1, reg.ST2, reg.ST3, reg.ST4, reg.ST5, reg.ST6, reg.ST7
	MM0, MM1, MM2, MM3, MM4, MM5, MM6, MM7 = reg.MM0, reg.MM1, reg.MM2, reg.MM3, reg.MM4, reg.MM5, reg.MM6, reg.MM7

	XMM0, XMM1, XMM2, XMM3 = reg.XMM0, reg.XMM1, reg.XMM2, reg.XMM3
	XMM4, XMM5, XMM6, XMM7 = reg.XMM4, reg.XMM5, reg.XMM6, reg.XMM7
	YMM0, YMM1, YMM2, YMM3 = reg.YMM0, reg.YMM1, reg.YMM2, reg.YMM3
	YMM4, YMM5, YMM6, YMM7 = reg.YMM4, reg.YMM5, reg.YMM6, reg.YMM7
	ZMM0, ZMM1, ZMM2, ZMM3 = reg.ZMM0, reg.ZMM1, reg.ZMM2, reg.ZMM3
	ZMM4, ZMM5, ZMM6, ZMM7 = reg.ZMM4, reg.ZMM5, reg.ZMM6, reg.ZMM7
	K0, K1, K2, K3, K4, K5, K6, K7 = reg.K0, reg.K1, reg.K2, reg.K3, reg.K4, reg.K5, reg.K6, reg.K7

	CR0, CR1, CR2, CR3, CR4, CR5, CR6, CR7 = reg.CR0, reg.CR1, reg.CR2, reg.CR3, reg.CR4, reg.CR5, reg.CR6, reg.CR7
	DR0, DR1, DR2, DR3, DR4, DR5, DR6, DR7 = reg.DR0, reg.DR1, reg.DR2, reg.DR3, reg.DR4, reg.DR5, reg.DR6, reg.DR7
)

// Memory operand constructors, by access width. Three per width, one
// question each: Mem is based, Abs is symbolic (a relocation), Addr is a
// direct address (no relocation).
var (
	Mem8, Mem16, Mem32, Mem64     = operand.Mem8, operand.Mem16, operand.Mem32, operand.Mem64
	Mem80, Mem128, Mem256, Mem512 = operand.Mem80, operand.Mem128, operand.Mem256, operand.Mem512

	Abs8, Abs16, Abs32, Abs64     = operand.Abs8, operand.Abs16, operand.Abs32, operand.Abs64
	Abs80, Abs128, Abs256, Abs512 = operand.Abs80, operand.Abs128, operand.Abs256, operand.Abs512

	Addr8, Addr16, Addr32, Addr64     = operand.Addr8, operand.Addr16, operand.Addr32, operand.Addr64
	Addr80, Addr128, Addr256, Addr512 = operand.Addr80, operand.Addr128, operand.Addr256, operand.Addr512
)

// Imm builds an immediate operand for Emit. Typed helpers take plain
// integers; this exists because a mnemonic-as-data path has no types.
func Imm(n int64) operand.Imm { return operand.NewImm(n) }

// Ref names a symbol and the link semantics that resolve it.
func Ref(name string, kind RefKind) operand.SymRef {
	return operand.Ref(name, operand.RelocKind(kind))
}

// NewLabel builds a same-section label operand for Emit. Typed helpers take
// the name as a string.
func NewLabel(name string) operand.Label { return operand.NewLabel(name) }

// FeatureSet is the module's gate: a base level plus extensions.
type FeatureSet = feature.Set

// NewFeatureSet returns the set for a base level with no extensions.
func NewFeatureSet(l feature.Level) feature.Set { return feature.New(l) }

// Base levels.
const (
	I386 = feature.I386
	I486 = feature.I486
	I586 = feature.I586
	I686 = feature.I686
)

// Extensions, canonical order.
const (
	MMX, FXSR, SSE, SSE2, SSE3, SSSE3       = feature.MMX, feature.FXSR, feature.SSE, feature.SSE2, feature.SSE3, feature.SSSE3
	SSE41, SSE42, POPCNT, AES, PCLMUL       = feature.SSE41, feature.SSE42, feature.POPCNT, feature.AES, feature.PCLMUL
	XSAVE, AVX, F16C, FMA, AVX2             = feature.XSAVE, feature.AVX, feature.F16C, feature.FMA, feature.AVX2
	BMI1, BMI2, LZCNT, MOVBE, ADX           = feature.BMI1, feature.BMI2, feature.LZCNT, feature.MOVBE, feature.ADX
	RDRAND, RDSEED, SHA                     = feature.RDRAND, feature.RDSEED, feature.SHA
	AVX512F, AVX512CD, AVX512VL             = feature.AVX512F, feature.AVX512CD, feature.AVX512VL
	AVX512BW, AVX512DQ                      = feature.AVX512BW, feature.AVX512DQ
)