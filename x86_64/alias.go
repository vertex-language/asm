// x86_64/alias.go
//
// Re-exports, so a lowering imports one package and every term it needs —
// registers, operand constructors, feature spellings — is spelled x86_64.X.
// The feature constants must all be here: a gating diagnostic's note line
// prints x86_64.V2.Plus(x86_64.AVX512VL), and that has to compile.
package x86_64

import (
	"github.com/vertex-language/asm/x86_64/feature"
	"github.com/vertex-language/asm/x86_64/operand"
	"github.com/vertex-language/asm/x86_64/reg"
)

// ---- registers -------------------------------------------------------

const (
	RAX, EAX, AX, AL = reg.RAX, reg.EAX, reg.AX, reg.AL
	RCX, ECX, CX, CL = reg.RCX, reg.ECX, reg.CX, reg.CL
	RDX, EDX, DX, DL = reg.RDX, reg.EDX, reg.DX, reg.DL
	RBX, EBX, BX, BL = reg.RBX, reg.EBX, reg.BX, reg.BL
	RSP, ESP, SP     = reg.RSP, reg.ESP, reg.SP
	RBP, EBP, BP     = reg.RBP, reg.EBP, reg.BP
	RSI, ESI, SI     = reg.RSI, reg.ESI, reg.SI
	RDI, EDI, DI     = reg.RDI, reg.EDI, reg.DI

	R8, R8D, R8W, R8B     = reg.R8, reg.R8D, reg.R8W, reg.R8B
	R9, R9D, R9W, R9B     = reg.R9, reg.R9D, reg.R9W, reg.R9B
	R10, R10D, R10W, R10B = reg.R10, reg.R10D, reg.R10W, reg.R10B
	R11, R11D, R11W, R11B = reg.R11, reg.R11D, reg.R11W, reg.R11B
	R12, R12D, R12W, R12B = reg.R12, reg.R12D, reg.R12W, reg.R12B
	R13, R13D, R13W, R13B = reg.R13, reg.R13D, reg.R13W, reg.R13B
	R14, R14D, R14W, R14B = reg.R14, reg.R14D, reg.R14W, reg.R14B
	R15, R15D, R15W, R15B = reg.R15, reg.R15D, reg.R15W, reg.R15B

	SPL, BPL, SIL, DIL = reg.SPL, reg.BPL, reg.SIL, reg.DIL
	AH, CH, DH, BH     = reg.AH, reg.CH, reg.DH, reg.BH
)

const (
	XMM0, YMM0, ZMM0    = reg.XMM0, reg.YMM0, reg.ZMM0
	XMM1, YMM1, ZMM1    = reg.XMM1, reg.YMM1, reg.ZMM1
	XMM2, YMM2, ZMM2    = reg.XMM2, reg.YMM2, reg.ZMM2
	XMM3, YMM3, ZMM3    = reg.XMM3, reg.YMM3, reg.ZMM3
	XMM4, YMM4, ZMM4    = reg.XMM4, reg.YMM4, reg.ZMM4
	XMM5, YMM5, ZMM5    = reg.XMM5, reg.YMM5, reg.ZMM5
	XMM6, YMM6, ZMM6    = reg.XMM6, reg.YMM6, reg.ZMM6
	XMM7, YMM7, ZMM7    = reg.XMM7, reg.YMM7, reg.ZMM7
	XMM8, YMM8, ZMM8    = reg.XMM8, reg.YMM8, reg.ZMM8
	XMM9, YMM9, ZMM9    = reg.XMM9, reg.YMM9, reg.ZMM9
	XMM10, YMM10, ZMM10 = reg.XMM10, reg.YMM10, reg.ZMM10
	XMM11, YMM11, ZMM11 = reg.XMM11, reg.YMM11, reg.ZMM11
	XMM12, YMM12, ZMM12 = reg.XMM12, reg.YMM12, reg.ZMM12
	XMM13, YMM13, ZMM13 = reg.XMM13, reg.YMM13, reg.ZMM13
	XMM14, YMM14, ZMM14 = reg.XMM14, reg.YMM14, reg.ZMM14
	XMM15, YMM15, ZMM15 = reg.XMM15, reg.YMM15, reg.ZMM15
	XMM16, YMM16, ZMM16 = reg.XMM16, reg.YMM16, reg.ZMM16
	XMM17, YMM17, ZMM17 = reg.XMM17, reg.YMM17, reg.ZMM17
	XMM18, YMM18, ZMM18 = reg.XMM18, reg.YMM18, reg.ZMM18
	XMM19, YMM19, ZMM19 = reg.XMM19, reg.YMM19, reg.ZMM19
	XMM20, YMM20, ZMM20 = reg.XMM20, reg.YMM20, reg.ZMM20
	XMM21, YMM21, ZMM21 = reg.XMM21, reg.YMM21, reg.ZMM21
	XMM22, YMM22, ZMM22 = reg.XMM22, reg.YMM22, reg.ZMM22
	XMM23, YMM23, ZMM23 = reg.XMM23, reg.YMM23, reg.ZMM23
	XMM24, YMM24, ZMM24 = reg.XMM24, reg.YMM24, reg.ZMM24
	XMM25, YMM25, ZMM25 = reg.XMM25, reg.YMM25, reg.ZMM25
	XMM26, YMM26, ZMM26 = reg.XMM26, reg.YMM26, reg.ZMM26
	XMM27, YMM27, ZMM27 = reg.XMM27, reg.YMM27, reg.ZMM27
	XMM28, YMM28, ZMM28 = reg.XMM28, reg.YMM28, reg.ZMM28
	XMM29, YMM29, ZMM29 = reg.XMM29, reg.YMM29, reg.ZMM29
	XMM30, YMM30, ZMM30 = reg.XMM30, reg.YMM30, reg.ZMM30
	XMM31, YMM31, ZMM31 = reg.XMM31, reg.YMM31, reg.ZMM31
)

const (
	K0, K1, K2, K3 = reg.K0, reg.K1, reg.K2, reg.K3
	K4, K5, K6, K7 = reg.K4, reg.K5, reg.K6, reg.K7

	TMM0, TMM1, TMM2, TMM3 = reg.TMM0, reg.TMM1, reg.TMM2, reg.TMM3
	TMM4, TMM5, TMM6, TMM7 = reg.TMM4, reg.TMM5, reg.TMM6, reg.TMM7

	ST0, ST1, ST2, ST3 = reg.ST0, reg.ST1, reg.ST2, reg.ST3
	ST4, ST5, ST6, ST7 = reg.ST4, reg.ST5, reg.ST6, reg.ST7

	MM0, MM1, MM2, MM3 = reg.MM0, reg.MM1, reg.MM2, reg.MM3
	MM4, MM5, MM6, MM7 = reg.MM4, reg.MM5, reg.MM6, reg.MM7

	ES, CS, SS, DS, FS, GS = reg.ES, reg.CS, reg.SS, reg.DS, reg.FS, reg.GS
)

// ---- operand constructors ---------------------------------------------

// The memory constructors, one per access width, plus the width-agnostic
// forms. Mem64(RBX).Disp(8) is the operand in `mov rax, [rbx+8]`.
var (
	Mem8   = operand.Mem8
	Mem16  = operand.Mem16
	Mem32  = operand.Mem32
	Mem64  = operand.Mem64
	Mem128 = operand.Mem128
	Mem256 = operand.Mem256
	Mem512 = operand.Mem512

	// Abs is an absolute reference, disp32 sign-extended; AbsSym the same
	// through a symbol; RIPRel is how PIC reaches static data here — one
	// instruction, one reference, no thunk.
	Abs        = operand.Abs
	AbsSym     = operand.AbsSym
	RIPRel     = operand.RIPRel
	RIPRelDisp = operand.RIPRelDisp

	Imm  = func(v int64) operand.Imm { return operand.Imm(v) }
	Uimm = operand.Uimm
)

// ---- features -----------------------------------------------------------

// The levels and every extension, re-exported so the note line of a gating
// diagnostic — x86_64.V2.Plus(x86_64.AVX512VL) — names something that
// compiles at the caller.
const (
	V1 = feature.V1
	V2 = feature.V2
	V3 = feature.V3
	V4 = feature.V4

	MMX, SSE, SSE2, SSE3    = feature.MMX, feature.SSE, feature.SSE2, feature.SSE3
	SSSE3, SSE41, SSE42     = feature.SSSE3, feature.SSE41, feature.SSE42
	AVX, AVX2               = feature.AVX, feature.AVX2
	POPCNT, CMPXCHG16B      = feature.POPCNT, feature.CMPXCHG16B
	LAHFSAHF, LZCNT, MOVBE  = feature.LAHFSAHF, feature.LZCNT, feature.MOVBE
	BMI1, BMI2, ADX         = feature.BMI1, feature.BMI2, feature.ADX
	F16C, FMA, FSGSBASE     = feature.F16C, feature.FMA, feature.FSGSBASE
	RDRAND, RDSEED          = feature.RDRAND, feature.RDSEED
	AES, VAES               = feature.AES, feature.VAES
	PCLMULQDQ, VPCLMULQDQ   = feature.PCLMULQDQ, feature.VPCLMULQDQ
	SHA                     = feature.SHA

	AVX512F, AVX512CD, AVX512BW   = feature.AVX512F, feature.AVX512CD, feature.AVX512BW
	AVX512DQ, AVX512VL            = feature.AVX512DQ, feature.AVX512VL
	AVX512IFMA, AVX512VBMI        = feature.AVX512IFMA, feature.AVX512VBMI
	AVX512VBMI2, AVX512VNNI       = feature.AVX512VBMI2, feature.AVX512VNNI
	AVX512BITALG, AVX512VPOPCNTDQ = feature.AVX512BITALG, feature.AVX512VPOPCNTDQ
	AVX512BF16, AVX512FP16        = feature.AVX512BF16, feature.AVX512FP16

	AMXTILE, AMXINT8, AMXBF16 = feature.AMXTILE, feature.AMXINT8, feature.AMXBF16
)

// FeatureSet is the value WithFeatures takes; Level and Feature spell it.
type FeatureSet = feature.Set

// Baseline is V1: SSE2 and nothing above it.
var Baseline = feature.Baseline

// Empty is the set with nothing above Base — not Baseline. It exists for
// callers building a set from scratch, and it is what a level-less
// ParseFeatures spelling starts from.
var Empty = feature.Empty

// ParseFeatures accepts the spellings the world writes — "x86-64-v3",
// "x86-64-v4-avx512vl", "sse2+aes" — and Set.String() round-trips through it.
var ParseFeatures = feature.ParseFeatures