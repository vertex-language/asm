package aarch64

import (
	"errors"
	"fmt"
	"strings"
)

// The sentinels. Sticky and first-wins: after any failure every builder call
// is a no-op, and Finalize surfaces the first error, positioned.
var (
	ErrFeature   = errors.New("form gated behind an extension not in the module's set")
	ErrForm      = errors.New("operand combination no declared form accepts")
	ErrDuplicate = errors.New("duplicate definition")
	ErrUndefined = errors.New("reference to a name neither defined nor declared extern")
	ErrRange     = errors.New("value does not fit its pinned field")
	ErrAlign     = errors.New("bad alignment")
	ErrFinalized = errors.New("builder call after Finalize")

	// ErrBitmask is separate from ErrRange because the fix is different in
	// kind: the constant has to be materialized, not made smaller.
	ErrBitmask = errors.New("constant is not a logical immediate")
)

// Error is the concrete error type. errors.Is works against every sentinel.
type Error struct {
	Section string
	Offset  int
	Context string
	Notes   []string

	sentinel error
	cause    error
}

func (e *Error) Error() string {
	var b strings.Builder
	if e.Section != "" {
		fmt.Fprintf(&b, "%s+%#x: ", e.Section, e.Offset)
	}
	if e.Context != "" {
		b.WriteString(e.Context + ": ")
	}
	if e.cause != nil {
		b.WriteString(e.cause.Error())
	} else if e.sentinel != nil {
		b.WriteString(e.sentinel.Error())
	}
	for _, n := range e.Notes {
		b.WriteString("\n  note: " + n)
	}
	return b.String()
}

// Unwrap exposes both the sentinel and the underlying encoder/resolver error,
// so errors.Is(err, ErrRange) and errors.As(err, &rangeErr) both work.
func (e *Error) Unwrap() []error {
	var out []error
	if e.sentinel != nil {
		out = append(out, e.sentinel)
	}
	if e.cause != nil {
		out = append(out, e.cause)
	}
	return out
}