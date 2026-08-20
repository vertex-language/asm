// x86_64/error.go
package x86_64

import (
	"errors"
	"fmt"
	"strings"
)

// The common sentinels per the repo README, plus the one this architecture
// demands. errors.Is works against every one of them through *Error.
var (
	ErrFeature   = errors.New("form gated behind an extension not in the module's set")
	ErrForm      = errors.New("operand combination no declared form accepts")
	ErrDuplicate = errors.New("label defined twice in one section, or a symbol in two")
	ErrUndefined = errors.New("reference to a name neither defined nor declared Extern")
	ErrRange     = errors.New("label displacement does not fit its pinned field")
	ErrAlign     = errors.New("alignment is not a power of two")
	ErrFinalized = errors.New("builder call after Finalize")

	// ErrEncoding is a legal-looking combination the silicon refuses: AH
	// with a REX prefix, XMM16+ outside EVEX, {z} without a mask, rounding
	// control off 512 bits. encode/ names the rule; this names the family.
	ErrEncoding = errors.New("operand combination the silicon refuses")
)

// Error is the concrete error type. Errors are sticky and first-wins: every
// builder call after a failure is a no-op, and Finalize surfaces the first
// error, positioned.
type Error struct {
	Section string // section name, or "" for a module-level failure
	Offset  int    // where the section stood when the call failed
	Context string // what was being built, and what went wrong
	Err     error  // the sentinel, reachable through errors.Is
	Notes   []string
}

func (e *Error) Error() string {
	var b strings.Builder
	if e.Section != "" {
		fmt.Fprintf(&b, "%s+%#x: ", e.Section, e.Offset)
	}
	b.WriteString(e.Context)
	for _, n := range e.Notes {
		b.WriteString("\n  ")
		b.WriteString(n)
	}
	return b.String()
}

func (e *Error) Unwrap() error { return e.Err }