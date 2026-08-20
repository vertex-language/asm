// Package operand is the i386 operand vocabulary: immediates, memory
// operands, labels, and symbol references.
//
// It imports reg and nothing else. Everything above — the table, the
// encoder, the root — speaks these types. The reference-kind values are
// defined here because this is the lowest package that must carry them;
// the root aliases the type, so there is one list and no conversion
// anywhere in the tree.
package operand

import "github.com/vertex-language/asm/i386/reg"

// Operand is anything an instruction can take. It is reg.Value — the sealed
// interface — under the name the instruction surface uses. The seal lives in
// reg because reg is the lower package; see reg/value.go. Nothing outside
// reg and this package can implement it, which is what keeps Emit's variadic
// from accepting an arbitrary type.
type Operand = reg.Value

// The r/m interfaces: a register or a memory operand of one access width.
// reg's marker methods and the M* types below are the only implementations.
type (
	RM8   interface {
		Operand
		RM8()
	}
	RM16 interface {
		Operand
		RM16()
	}
	RM32 interface {
		Operand
		RM32()
	}
	RM64 interface {
		Operand
		RM64()
	}
	RM128 interface {
		Operand
		RM128()
	}
)