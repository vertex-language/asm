package aarch64

import "github.com/vertex-language/asm/aarch64/feature"

type SectionKind uint8

const (
	Text SectionKind = iota
	ROData
	Data
	BSS
	sectionKindCount
)

var sectionNames = [sectionKindCount]string{".text", ".rodata", ".data", ".bss"}

func (k SectionKind) String() string {
	if k >= sectionKindCount {
		return "?"
	}
	return sectionNames[k]
}

// Module sits after every decision an assembler can make and before every
// decision only a linker can make. See the repo README for the model.
type Module struct {
	features feature.Set

	sections []*Section
	byKind   [sectionKindCount]*Section

	symbols map[string]*Section // symbol name -> defining section
	externs map[string]bool

	err       *Error
	finalized bool
}

type Option func(*Module)

// WithFeatures fixes the feature set at construction. A gate that changed
// mid-module would make already-emitted diagnostics unfalsifiable.
func WithFeatures(s feature.Set) Option { return func(m *Module) { m.features = s } }

func NewModule(opts ...Option) *Module {
	m := &Module{
		features: feature.Baseline(),
		symbols:  map[string]*Section{},
		externs:  map[string]bool{},
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

func (m *Module) Features() feature.Set { return m.features }

// Section returns the one section of a kind, creating it on first use.
// Sections() lists them in creation order.
func (m *Module) Section(k SectionKind) *Section {
	if k >= sectionKindCount {
		panic("aarch64: no such section kind")
	}
	if s := m.byKind[k]; s != nil {
		return s
	}
	s := &Section{m: m, kind: k, labels: map[string]int{}}
	if m.finalized {
		m.setErr(&Error{Section: k.String(), Context: "Section", sentinel: ErrFinalized})
		return s // detached: its builder calls no-op, and it never joins Sections()
	}
	m.byKind[k] = s
	m.sections = append(m.sections, s)
	return s
}

func (m *Module) Sections() []*Section { return m.sections }

// Extern declares an imported name. Finalize refuses a reference to a name
// neither defined nor declared.
func (m *Module) Extern(name string) {
	if m.err != nil {
		return
	}
	if m.finalized {
		m.setErr(&Error{Context: "Extern " + name, sentinel: ErrFinalized})
		return
	}
	m.externs[name] = true
}

// Err is the sticky error, for callers that want to bail early.
func (m *Module) Err() error {
	if m.err == nil {
		return nil
	}
	return m.err
}

func (m *Module) setErr(e *Error) {
	if m.err == nil {
		m.err = e
	}
}

// Finalize: surface the sticky error, patch same-section label references,
// close symbol sizes, verify every remaining reference. Then the module is
// immutable, pure data.
func (m *Module) Finalize() error {
	if m.finalized {
		return m.Err()
	}
	m.finalized = true
	if m.err != nil {
		return m.err
	}
	for _, s := range m.sections {
		s.resolve()
		if m.err != nil {
			return m.err
		}
	}
	for _, s := range m.sections {
		s.closeSizes()
	}
	return m.Err()
}