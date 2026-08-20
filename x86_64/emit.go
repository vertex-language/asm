// x86_64/emit.go
package x86_64

import (
	"errors"

	"github.com/vertex-language/asm/x86_64/internal/encode"
	"github.com/vertex-language/asm/x86_64/internal/isa"
	"github.com/vertex-language/asm/x86_64/operand"
)

// Emit is the escape hatch: runtime form resolution, shortest legal encoding
// among matching forms, ties broken by table order. It exists for
// table-driven emission where the mnemonic is data. It never picks a
// different instruction than the one named — no form of this mnemonic taking
// these operands is ErrForm, not a relaxation into one that does.
//
// Note that Emit with a label will happily pick a rel8 jump because it is
// shortest, and pay for it at Finalize if the target is far.
func (s *Section) Emit(mnemonic string, ops ...any) {
	if s.blocked() {
		return
	}
	ops = normalize(ops)

	args, err := encode.Args(ops...)
	if err != nil {
		s.fail(err.Error(), ErrForm)
		return
	}

	f, err := isa.Resolve(s.m.feats, mnemonic, args...)
	if err != nil {
		var g *isa.GateError
		if errors.As(err, &g) {
			s.fail(err.Error(), ErrFeature,
				"note: enable with "+qualify(g.Active.Plus(g.Need).GoExpr()))
			return
		}
		s.fail(err.Error(), ErrForm)
		return
	}
	s.place(f, encode.Opts{}, ops...)
}

// normalize lets Emit take bare Go integers where an operand.Imm is meant,
// since a mnemonic-as-data caller usually has values, not operand types.
func normalize(ops []any) []any {
	out := make([]any, len(ops))
	for i, o := range ops {
		switch v := o.(type) {
		case int:
			out[i] = operand.Imm(v)
		case int8:
			out[i] = operand.Imm(v)
		case int16:
			out[i] = operand.Imm(v)
		case int32:
			out[i] = operand.Imm(v)
		case int64:
			out[i] = operand.Imm(v)
		case uint8:
			out[i] = operand.Imm(v)
		case uint16:
			out[i] = operand.Imm(v)
		case uint32:
			out[i] = operand.Imm(v)
		case uint64:
			out[i] = operand.Uimm(v)
		default:
			out[i] = o
		}
	}
	return out
}