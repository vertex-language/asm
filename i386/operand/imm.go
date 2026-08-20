package operand

import (
	"strconv"

	"github.com/vertex-language/asm/i386/reg"
)

// Imm is an immediate value. It carries no width: width is a property of the
// field the form pins, not of the value, and whether the value fits that
// field is encode's range check.
type Imm struct {
	reg.Seal
	v int64
}

// NewImm builds an immediate. The typed helpers take plain integers and call
// this themselves; it is exported for Emit, whose mnemonic-as-data callers
// have values, not operand types.
func NewImm(v int64) Imm { return Imm{v: v} }

func (i Imm) Int64() int64   { return i.v }
func (i Imm) Bits() int      { return 0 }
func (i Imm) String() string { return strconv.FormatInt(i.v, 10) }