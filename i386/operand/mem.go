package operand

import (
	"errors"
	"fmt"

	"github.com/vertex-language/asm/i386/reg"
)

// ErrOperand is the sentinel every malformed-operand error unwraps to.
var ErrOperand = errors.New("operand")

// Memory is any memory operand, of any access width. It is what LEA takes
// and what the encoder's ModRM construction consumes.
type Memory interface {
	Operand
	Base() (reg.R32, bool)
	IndexReg() (index reg.R32, scale uint8, ok bool)
	Segment() (reg.Sreg, bool)
	Displacement() int32
	Symbol() (SymRef, bool)
	Err() error
}

// mem is the address computation every width shares:
// seg:[base + index*scale + disp/sym]. Construction errors are sticky on the
// value and surfaced by Err() when the operand is encoded, so a builder
// chain is not followed by a run of error checks.
type mem struct {
	reg.Seal
	base     reg.R32
	hasBase  bool
	index    reg.R32
	scale    uint8
	hasIndex bool
	seg      reg.Sreg
	hasSeg   bool
	disp     int32
	sym      SymRef
	hasSym   bool
	err      error
}

func based(b reg.R32) mem   { return mem{base: b, hasBase: true} }
func symbolic(s SymRef) mem { return mem{sym: s, hasSym: true} }
func direct(d int32) mem    { return mem{disp: d} }

// withDisp adds n to the address. With a symbol present the displacement is
// the symbol's addend, so it folds there and Displacement stays zero — one
// place for the number, whichever order the chain was written in.
func (m mem) withDisp(n int32) mem {
	if m.hasSym {
		m.sym = m.sym.Plus(n)
		return m
	}
	m.disp += n
	return m
}

func (m mem) withIndex(r reg.R32, scale int) mem {
	if m.err != nil {
		return m
	}
	if r == reg.ESP {
		// SIB.index = 100 means "no index", so ESP has no index encoding at
		// all. Refused here, at construction, rather than deep in the
		// encoder where the diagnostic would point at the instruction.
		m.err = fmt.Errorf("%w: esp cannot be an index register: SIB index 100 means no index", ErrOperand)
		return m
	}
	switch scale {
	case 1, 2, 4, 8:
	default:
		m.err = fmt.Errorf("%w: scale %d is not 1, 2, 4, or 8", ErrOperand, scale)
		return m
	}
	m.index, m.scale, m.hasIndex = r, uint8(scale), true
	return m
}

func (m mem) withSeg(s reg.Sreg) mem { m.seg, m.hasSeg = s, true; return m }

func (m mem) withSym(r SymRef) mem {
	if m.disp != 0 {
		r, m.disp = r.Plus(m.disp), 0
	}
	m.sym, m.hasSym = r, true
	return m
}

func (m mem) Base() (reg.R32, bool)                 { return m.base, m.hasBase }
func (m mem) IndexReg() (reg.R32, uint8, bool)      { return m.index, m.scale, m.hasIndex }
func (m mem) Segment() (reg.Sreg, bool)             { return m.seg, m.hasSeg }
func (m mem) Displacement() int32                   { return m.disp }
func (m mem) Symbol() (SymRef, bool)                { return m.sym, m.hasSym }
func (m mem) Err() error                            { return m.err }