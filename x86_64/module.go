// x86_64/module.go
package x86_64

import "github.com/vertex-language/asm/x86_64/feature"

// SectionKind is the load-time property bytes cannot carry.
type SectionKind uint8

const (
	Text SectionKind = iota
	ROData
	Data
	BSS

	numSectionKinds
)

var sectionNames = [numSectionKinds]string{
	Text: ".text", ROData: ".rodata", Data: ".data", BSS: ".bss",
}

func (k SectionKind) String() string {
	if int(k) >= len(sectionNames) {
		return "section(?)"
	}
	return sectionNames[k]
}

// IsCode reports whether Align pads with no-ops rather than zeros.
func (k SectionKind) IsCode() bool { return k == Text }

// Module sits after every decision an assembler can make and before every
// decision only a linker can make. Errors are sticky, first-wins, surfaced
// at Finalize; after Finalize the module is immutable, pure data.
type Module struct {
	feats   feature.Set
	byKind  map[SectionKind]*Section
	order   []*Section
	symOf   map[string]*Section // owner of every promoted symbol; module-global
	externs map[string]struct{}
	err     *Error
	done    bool
}

// Option configures NewModule.
type Option func(*Module)

// WithFeatures fixes the feature set at construction. A gate that changed
// mid-module would make already-emitted diagnostics unfalsifiable.
func WithFeatures(s feature.Set) Option { return func(m *Module) { m.feats = s } }

// NewModule builds an empty module. Default features are Baseline: V1,
// which is SSE2 and nothing above it.
func NewModule(opts ...Option) *Module {
	m := &Module{
		feats:   feature.Baseline(),
		byKind:  make(map[SectionKind]*Section),
		symOf:   make(map[string]*Section),
		externs: make(map[string]struct{}),
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// Features is the active feature set.
func (m *Module) Features() feature.Set { return m.feats }

// Section is the section of that kind, created on first use. One section per
// kind: the same *Section comes back every time.
//
// The freeze is total: after Finalize an existing section still comes back
// for reading, but asking for a kind that was never created is refused the
// same way any builder call is — ErrFinalized is recorded, and the section
// handed back is never registered, so Sections() is as frozen as the bytes.
func (m *Module) Section(k SectionKind) *Section {
	if s, ok := m.byKind[k]; ok {
		return s
	}
	s := &Section{m: m, kind: k, labels: make(map[string]int)}
	if m.done {
		s.fail("Section after Finalize", ErrFinalized)
		return s
	}
	m.byKind[k] = s
	m.order = append(m.order, s)
	return s
}

// Sections lists the sections in creation order.
func (m *Module) Sections() []*Section {
	return append([]*Section(nil), m.order...)
}

// Extern declares names resolved outside this module. Finalize refuses a
// reference to a name neither defined nor declared here.
func (m *Module) Extern(names ...string) {
	if m.err != nil {
		return
	}
	if m.done {
		m.fail(&Error{Context: "Extern after Finalize", Err: ErrFinalized})
		return
	}
	for _, n := range names {
		m.externs[n] = struct{}{}
	}
}

// Err is the sticky error, or nil. Finalize is the usual way to read it;
// this exists for callers that want to bail early.
func (m *Module) Err() error {
	if m.err != nil {
		return m.err
	}
	return nil
}

// Finalize does four jobs, then freezes the module:
//
//  1. surface the sticky error — every builder call after a failure was a
//     no-op, and the first error is the one you get, positioned
//  2. patch every same-section label reference into the bytes and drop it,
//     refusing with ErrRange any displacement that does not fit its field
//  3. close every symbol's size at the next symbol or section end
//  4. verify every remaining reference targets a defined symbol or a
//     declared Extern
//
// Finalize is idempotent: a second call returns the same answer.
func (m *Module) Finalize() error {
	if m.done {
		return m.Err()
	}
	m.done = true
	if m.err != nil {
		return m.err
	}
	for _, s := range m.order {
		if e := s.finalize(); e != nil {
			m.fail(e)
			return m.err
		}
	}
	return nil
}

func (m *Module) fail(e *Error) {
	if m.err == nil {
		m.err = e
	}
}

// defined reports whether name is a promoted symbol anywhere in the module.
func (m *Module) defined(name string) bool { return m.symOf[name] != nil }

// resolvable reports whether a reference to name may survive Finalize.
func (m *Module) resolvable(name string) bool {
	if m.defined(name) {
		return true
	}
	_, ok := m.externs[name]
	return ok
}