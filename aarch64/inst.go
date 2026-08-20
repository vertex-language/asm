package aarch64

import (
	"encoding/binary"
	"errors"
	"sync"

	"github.com/vertex-language/asm/aarch64/internal/encode"
	"github.com/vertex-language/asm/aarch64/internal/isa"
)

// The helpers bind to their forms by GoName lookup, at the first call —
// appended table rows break nothing, and a removed or renamed form panics
// loudly by name rather than silently binding to the wrong row.
var (
	formOnce   sync.Once
	formByName map[string]*isa.Form
)

func formFor(goName string) *isa.Form {
	formOnce.Do(func() {
		all := isa.All()
		formByName = make(map[string]*isa.Form, len(all))
		for _, f := range all {
			formByName[f.GoName()] = f
		}
	})
	f := formByName[goName]
	if f == nil {
		panic("aarch64: no form in the table generates the helper " + goName +
			"; a row was removed or renamed")
	}
	return f
}

// inst is what every typed helper calls: gate, encode, append, record.
func (s *Section) inst(goName string, ops ...any) {
	if !s.ready(goName) {
		return
	}
	s.place(formFor(goName), ops)
}

func (s *Section) place(f *isa.Form, ops []any) {
	// EncodeForm does not consult the feature set — the gate is checked here,
	// where the helper was called, so the diagnostic is positioned.
	if !f.Enabled(s.m.features) {
		cause := &isa.GateError{Form: f, Active: s.m.features}
		s.fail(ErrFeature, cause, f.Mnem)
		return
	}
	word, fixups, err := encode.EncodeForm(f, ops, encode.Opts{Offset: len(s.buf)})
	if err != nil {
		s.fail(sentinelFor(err), err, f.Mnem)
		return
	}
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], word)
	s.buf = append(s.buf, b[:]...)
	s.fixups = append(s.fixups, fixups...)
}

// sentinelFor maps the resolver's and encoder's typed errors onto the
// sentinels errors.Is answers for.
func sentinelFor(err error) error {
	var (
		bitmask *encode.BitmaskError
		rng     *encode.RangeError
		gate    *isa.GateError
	)
	switch {
	case errors.As(err, &bitmask):
		return ErrBitmask
	case errors.As(err, &rng):
		return ErrRange
	case errors.As(err, &gate):
		return ErrFeature
	}
	return ErrForm
}