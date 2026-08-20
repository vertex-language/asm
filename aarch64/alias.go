package aarch64

import (
	"github.com/vertex-language/asm/aarch64/feature"
	"github.com/vertex-language/asm/aarch64/operand"
	"github.com/vertex-language/asm/aarch64/reg"
)

// ---- registers ----

const (
	X0, X1, X2, X3, X4, X5, X6, X7         = reg.X0, reg.X1, reg.X2, reg.X3, reg.X4, reg.X5, reg.X6, reg.X7
	X8, X9, X10, X11, X12, X13, X14, X15   = reg.X8, reg.X9, reg.X10, reg.X11, reg.X12, reg.X13, reg.X14, reg.X15
	X16, X17, X18, X19, X20, X21, X22, X23 = reg.X16, reg.X17, reg.X18, reg.X19, reg.X20, reg.X21, reg.X22, reg.X23
	X24, X25, X26, X27, X28, X29, X30, XZR = reg.X24, reg.X25, reg.X26, reg.X27, reg.X28, reg.X29, reg.X30, reg.XZR
)

const (
	W0, W1, W2, W3, W4, W5, W6, W7         = reg.W0, reg.W1, reg.W2, reg.W3, reg.W4, reg.W5, reg.W6, reg.W7
	W8, W9, W10, W11, W12, W13, W14, W15   = reg.W8, reg.W9, reg.W10, reg.W11, reg.W12, reg.W13, reg.W14, reg.W15
	W16, W17, W18, W19, W20, W21, W22, W23 = reg.W16, reg.W17, reg.W18, reg.W19, reg.W20, reg.W21, reg.W22, reg.W23
	W24, W25, W26, W27, W28, W29, W30, WZR = reg.W24, reg.W25, reg.W26, reg.W27, reg.W28, reg.W29, reg.W30, reg.WZR
)

const (
	SP  = reg.SP
	WSP = reg.WSP
	FP  = reg.FP
	LR  = reg.LR
	IP0 = reg.IP0
	IP1 = reg.IP1
)

const (
	V0, V1, V2, V3, V4, V5, V6, V7         = reg.V0, reg.V1, reg.V2, reg.V3, reg.V4, reg.V5, reg.V6, reg.V7
	V8, V9, V10, V11, V12, V13, V14, V15   = reg.V8, reg.V9, reg.V10, reg.V11, reg.V12, reg.V13, reg.V14, reg.V15
	V16, V17, V18, V19, V20, V21, V22, V23 = reg.V16, reg.V17, reg.V18, reg.V19, reg.V20, reg.V21, reg.V22, reg.V23
	V24, V25, V26, V27, V28, V29, V30, V31 = reg.V24, reg.V25, reg.V26, reg.V27, reg.V28, reg.V29, reg.V30, reg.V31
)

// Q0..Q31, D0..D31, S0..S31, H0..H31, B0..B31, Z0..Z31, P0..P15 follow the
// identical four-row pattern; generate or paste them the same way. FFR:
var FFR = reg.FFR

// ---- operand vocabulary ----

type (
	Imm     = operand.Imm
	Mem     = operand.Mem
	Label   = operand.Label
	SymRef  = operand.SymRef
	AddrRef = operand.AddrRef
	Cond    = operand.Cond
	Shift   = operand.Shift
	ShiftOp = operand.ShiftOp
	Extend  = operand.Extend
	Barrier = operand.Barrier
)

// Ref builds a symbol reference; a named kind is a request that blocks folding.
func Ref(name string, kind ...RelocKind) SymRef { return operand.Sym(name, kind...) }

var (
	Page       = operand.Page
	PageOff    = operand.PageOff
	GotPage    = operand.GotPage
	GotPageOff = operand.GotPageOff

	Shifted  = operand.Shifted
	Extended = operand.Extended
	NoShift  = operand.NoShift
)

type baseReg interface{ reg.X | reg.Xsp }

func MemOf[T baseReg](base T) Mem   { return operand.MemOf(base) }
func Mem8[T baseReg](base T) Mem    { return operand.Mem8(base) }
func Mem16[T baseReg](base T) Mem   { return operand.Mem16(base) }
func Mem32[T baseReg](base T) Mem   { return operand.Mem32(base) }
func Mem64[T baseReg](base T) Mem   { return operand.Mem64(base) }
func Mem128[T baseReg](base T) Mem  { return operand.Mem128(base) }

const (
	EQ, NE, CS, CC, MI, PL, VS, VC = operand.EQ, operand.NE, operand.CS, operand.CC, operand.MI, operand.PL, operand.VS, operand.VC
	HI, LS, GE, LT, GT, LE, AL, NV = operand.HI, operand.LS, operand.GE, operand.LT, operand.GT, operand.LE, operand.AL, operand.NV
	HS, LO                         = operand.HS, operand.LO
)

const (
	LSL, LSR, ASR, ROR = operand.LSL, operand.LSR, operand.ASR, operand.ROR
)

const (
	UXTB, UXTH, UXTW, UXTX = operand.UXTB, operand.UXTH, operand.UXTW, operand.UXTX
	SXTB, SXTH, SXTW, SXTX = operand.SXTB, operand.SXTH, operand.SXTW, operand.SXTX
	ExtLSL                 = operand.ExtLSL
)

const (
	SY, ISH, ISHLD, ISHST = operand.SY, operand.ISH, operand.ISHLD, operand.ISHST
	OSH, NSH, LD, ST      = operand.OSH, operand.NSH, operand.LD, operand.ST
)

// ---- features: every level and extension is re-exported at the root so the
// gating diagnostic's note line names something that compiles at the call site.

type Set = feature.Set

var (
	Baseline         = feature.Baseline
	NewFeatureSet    = feature.NewSet
	ParseFeatures    = feature.ParseFeatures
	MustParseFeatures = feature.MustParseFeatures
)

const (
	Armv8A, Armv8_1A, Armv8_2A, Armv8_3A, Armv8_4A = feature.Armv8A, feature.Armv8_1A, feature.Armv8_2A, feature.Armv8_3A, feature.Armv8_4A
	Armv8_5A, Armv8_6A, Armv8_7A, Armv8_8A, Armv8_9A = feature.Armv8_5A, feature.Armv8_6A, feature.Armv8_7A, feature.Armv8_8A, feature.Armv8_9A
	Armv9A, Armv9_1A, Armv9_2A, Armv9_3A, Armv9_4A, Armv9_5A = feature.Armv9A, feature.Armv9_1A, feature.Armv9_2A, feature.Armv9_3A, feature.Armv9_4A, feature.Armv9_5A
)

const (
	FeatFP, SIMD, FP16, BF16, I8MM, DotProd = feature.FP, feature.SIMD, feature.FP16, feature.BF16, feature.I8MM, feature.DotProd
	AES, SHA2, SHA3, SM4                    = feature.AES, feature.SHA2, feature.SHA3, feature.SM4
	LSE, RCPC, MOPS, CRC, RDMA              = feature.LSE, feature.RCPC, feature.MOPS, feature.CRC, feature.RDMA
	PAuth, BTI, MemTag, SB, RNG, TME        = feature.PAuth, feature.BTI, feature.MemTag, feature.SB, feature.RNG, feature.TME
	SVE, SVE2, SVE2BitPerm, SME, SME2       = feature.SVE, feature.SVE2, feature.SVE2BitPerm, feature.SME, feature.SME2
	// ...the remaining features re-export the same way, one const per
	// feature.* name; the GoExpr note lines depend on all of them existing.
)