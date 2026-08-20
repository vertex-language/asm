// x86_64/inst.go
//
// The machinery under the typed helpers. The helpers themselves are
// generated (helpers_*_gen.go) from the same table everything else reads;
// each binds to its form by GoName lookup, so appending rows breaks nothing
// and a removed or renamed form panics at the first call, naming the missing
// form, rather than silently binding to the wrong row.
package x86_64

import (
	"github.com/vertex-language/asm/x86_64/internal/encode"
	"github.com/vertex-language/asm/x86_64/internal/isa"
	"github.com/vertex-language/asm/x86_64/operand"
)

var formByGoName = func() map[string]*isa.Form {
	m := make(map[string]*isa.Form, isa.Count())
	for _, f := range isa.All() {
		n := f.GoName()
		if _, dup := m[n]; !dup {
			m[n] = f // the earlier row wins, as in Resolve's tie-break
		}
	}
	return m
}()

func form(goName string) *isa.Form {
	f := formByGoName[goName]
	if f == nil {
		panic("x86_64: helper bound to missing form " + goName + "; the isa table moved")
	}
	return f
}

var noOpts encode.Opts

// The r/m parameter types. Go has no union value type, and nothing below the
// root may implement a marker interface declared here, so an r/m slot — the
// one slot that takes a register or memory — is typed as a documented any.
// The class check runs at the call and a mismatch is a sticky, positioned
// ErrForm. Register-only, memory-only, immediate and label parameters remain
// compile-checked, which is most of them.
type (
	RM8  = any // reg.Reg8  or operand.M8
	RM16 = any // reg.Reg16 or operand.M16
	RM32 = any // reg.Reg32 or operand.M32
	RM64 = any // reg.Reg64 or operand.M64

	MmM64   = any // reg.Mm  or operand.M64
	XmmM32  = any // reg.Xmm or operand.M32
	XmmM64  = any // reg.Xmm or operand.M64
	XmmM128 = any // reg.Xmm or operand.M128
	YmmM256 = any // reg.Ymm or operand.M256
	ZmmM512 = any // reg.Zmm or operand.M512
	KM64    = any // reg.K   or operand.M64

	// Memory is an m-of-no-particular-width slot: lea, prefetch, tile
	// loads. operand.Mem or any of the width-fixed wrappers.
	Memory = any
)

// Opt is an EVEX modifier that is one bit with no register behind it:
// zeroing, broadcast, and rounding ride as options, while a writemask names
// a register and so is a reg.K parameter on the masked helper.
type Opt func(*encode.Opts)

// Zeroing selects zeroing-masking, {z}. Legal only alongside a nonzero mask.
func Zeroing() Opt { return func(o *encode.Opts) { o.Zero = true } }

// Broadcast selects embedded broadcast, {1toN}, over the memory operand.
func Broadcast() Opt { return func(o *encode.Opts) { o.Broadcast = true } }

// SAE suppresses all exceptions without naming a rounding mode.
func SAE() Opt { return func(o *encode.Opts) { o.SAE = true } }

// The embedded rounding modes. Only encodable at 512 bits, register
// operands only; encode/ refuses the rest by name.
func RoundNearest() Opt { return func(o *encode.Opts) { o.Round = encode.RoundNearest } }
func RoundDown() Opt    { return func(o *encode.Opts) { o.Round = encode.RoundDown } }
func RoundUp() Opt      { return func(o *encode.Opts) { o.Round = encode.RoundUp } }
func RoundZero() Opt    { return func(o *encode.Opts) { o.Round = encode.RoundZero } }

func optset(opts []Opt) encode.Opts {
	var o encode.Opts
	for _, f := range opts {
		f(&o)
	}
	return o
}

// ---- branch and call helpers --------------------------------------------
//
// Branch targets split by where they resolve. The Label helpers are
// same-section, patched at Finalize, no relocation; the Ref helpers leave
// the module and survive into Refs(). Short pins rel8; the plain names pin
// rel32; there is no relaxation between them — a short branch to a far
// target is ErrRange at Finalize.

// JmpLabel is JMP rel32 to a same-section label.
func (s *Section) JmpLabel(name string) {
	s.place(form("JmpRel32"), noOpts, operand.Label(name))
}

// JmpShortLabel is JMP rel8 to a same-section label.
func (s *Section) JmpShortLabel(name string) {
	s.place(form("JmpRel8"), noOpts, operand.Label(name))
}

// CallLabel is CALL rel32 to a function compiled into this module: a
// bare-label call — no reference, no relocation.
func (s *Section) CallLabel(name string) {
	s.place(form("CallRel32"), noOpts, operand.Label(name))
}

// JmpRef is JMP rel32 through a reference that leaves the module.
func (s *Section) JmpRef(t operand.SymRef) {
	s.place(form("JmpRel32"), noOpts, t)
}

// CallRef is CALL rel32 through a reference: e8, a 4-byte hole, and an
// entry in Refs(). With RefNone the kind resolves to RefPLT.
func (s *Section) CallRef(t operand.SymRef) {
	s.place(form("CallRel32"), noOpts, t)
}