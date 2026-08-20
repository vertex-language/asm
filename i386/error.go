package i386

import (
	"errors"
	"fmt"
	"strings"
)

// Sentinels. Errors are sticky and first-wins: every builder call after a
// failure is a no-op, and Finalize surfaces the first error, positioned.
var (
	// ErrFeature is a form gated behind an extension not in the module's set.
	ErrFeature = errors.New("feature-gated form")

	// ErrForm is an operand combination no declared form accepts.
	ErrForm = errors.New("no matching form")

	// ErrDuplicate is a label defined twice in one section, or a symbol
	// defined in two sections.
	ErrDuplicate = errors.New("duplicate definition")

	// ErrUndefined is a reference to a name neither defined nor declared
	// Extern.
	ErrUndefined = errors.New("undefined reference")

	// ErrRange is a value that does not fit the field its form pins: a
	// too-big immediate at the call, or a label displacement at Finalize.
	// There is no branch relaxation and no silent form substitution:
	// the failure is loud instead of the bytes being different.
	ErrRange = errors.New("value out of range")

	// ErrAlign is an alignment that is not a power of two.
	ErrAlign = errors.New("bad alignment")

	// ErrFinalized is a builder call after Finalize.
	ErrFinalized = errors.New("module is finalized")
)

// Error is a positioned assembler error: which section, at what offset,
// doing what, under which sentinel, caused by what.
//
// Unwrap exposes both the sentinel and the cause, so errors.Is works
// against every sentinel above. The cause's type is internal; anything a
// caller might need from it is carried in Notes as text, which is the
// contract — internal packages never become API through the error chain.
type Error struct {
	Sentinel error
	Cause    error // underlying resolver or encoder error; may be nil
	Section  string
	Offset   int
	Context  string
	Notes    []string
}

func (e *Error) Error() string {
	var b strings.Builder
	if e.Section != "" {
		fmt.Fprintf(&b, "%s+%#x: ", e.Section, e.Offset)
	}
	b.WriteString(e.Context)
	b.WriteString(": ")
	b.WriteString(e.Sentinel.Error())
	for _, n := range e.Notes {
		b.WriteString("\n  note: ")
		b.WriteString(n)
	}
	return b.String()
}

func (e *Error) Unwrap() []error {
	if e.Cause == nil {
		return []error{e.Sentinel}
	}
	return []error{e.Sentinel, e.Cause}
}