package i386

import (
	"fmt"
	"strings"

	"github.com/vertex-language/asm/i386/internal/encode"
	"github.com/vertex-language/asm/i386/internal/isa"
	"github.com/vertex-language/asm/i386/operand"
)

// Emit is the runtime-resolved escape hatch: shortest legal encoding among
// the forms that match, ties broken by table order. It picks an *encoding*
// of the instruction you named, never a different instruction. If you know
// the instruction at compile time, use the typed helper.
func (s *Section) Emit(mnemonic string, ops ...operand.Operand) {
	if !s.ok() {
		return
	}
	match, gated := isa.Resolve(mnemonic, s.m.features, ops)
	if len(match) == 0 {
		if len(gated) > 0 {
			gates := make([]string, 0, len(gated))
			seen := map[string]bool{}
			for _, f := range gated {
				if g := f.Gate(); g != "" && !seen[g] {
					seen[g] = true
					gates = append(gates, g)
				}
			}
			s.fail(ErrFeature, "%s %s is gated behind %s (module has %s)",
				strings.ToUpper(mnemonic), opsString(ops),
				strings.Join(gates, ", "), s.m.features)
			return
		}
		s.fail(ErrForm, "no form of %q accepts %s", mnemonic, opsString(ops))
		return
	}

	// Encode every candidate and take the shortest; Resolve returned them in
	// table order, and a strict < keeps the earlier form on a tie.
	var (
		best      encode.Inst
		bestForm  *isa.Form
		firstErr  error
		firstForm *isa.Form
	)
	for _, f := range match {
		in, err := encode.Encode(f, ops)
		if err != nil {
			if firstErr == nil {
				firstErr, firstForm = err, f
			}
			continue
		}
		if bestForm == nil || in.Len() < best.Len() {
			best, bestForm = in, f
		}
	}
	if bestForm == nil {
		// Every candidate refused. Resolve matched on values, so a range
		// failure is rare here — but a sticky operand-construction error
		// flows through this path, and failEncode classifies either way,
		// carrying the encoder's error as the cause.
		s.failEncode(firstForm.Signature(), firstErr)
		return
	}
	s.place(best)
}

func opsString(ops []operand.Operand) string {
	if len(ops) == 0 {
		return "(no operands)"
	}
	parts := make([]string, len(ops))
	for i, o := range ops {
		parts[i] = fmt.Sprintf("%v", o)
	}
	return strings.Join(parts, ", ")
}