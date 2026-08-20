// Package i386 is an in-memory module of finished Intel386 machine code:
// decided bytes plus the three facts bytes cannot carry — references,
// symbols, and section kinds. It writes no object files and knows no IR.
package i386

import (
	"fmt"

	"github.com/vertex-language/asm/i386/feature"
)

// SectionKind is a load-time property of a section, carried as data for
// the lowering to interpret. It is not a file-format section header.
type SectionKind uint8

const (
	Text SectionKind = iota
	Data
	ROData
	BSS
)

var sectionNames = [...]string{".text", ".data", ".rodata", ".bss"}

func (k SectionKind) String() string {
	if int(k) < len(sectionNames) {
		return sectionNames[k]
	}
	return "?"
}

// Module sits after every decision an assembler can make and before every
// decision only a linker can make. Builder calls append; Finalize patches
// same-section labels, verifies references, and freezes the module.
type Module struct {
	features  feature.Set
	sections  []*Section
	byKind    map[SectionKind]*Section
	externs   map[string]bool
	err       error
	finalized bool
}

// Option configures a Module at construction. Nothing is configurable
// after: a gate that changed mid-module would make already-emitted
// diagnostics unfalsifiable.
type Option func(*Module)

// WithFeatures fixes the module's feature set.
func WithFeatures(s feature.Set) Option {
	return func(m *Module) { m.features = s }
}

// NewModule builds an empty module at the default feature set.
func NewModule(opts ...Option) *Module {
	m := &Module{
		features: feature.Default(),
		byKind:   make(map[SectionKind]*Section),
		externs:  make(map[string]bool),
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// Features is the set fixed at construction.
func (m *Module) Features() feature.Set { return m.features }

// Section returns the section of a kind, creating it on first use. There
// is one section per kind; Sections() lists them in creation order.
//
// After Finalize, an existing section is still returned — the read side is
// the point of a finalized module — but a kind not yet created is refused:
// the handle handed back is usable and every builder call on it is a no-op,
// and it is not added to the module, so Sections() is as frozen as the
// bytes are.
func (m *Module) Section(k SectionKind) *Section {
	if s, ok := m.byKind[k]; ok {
		return s
	}
	s := &Section{m: m, kind: k, labels: make(map[string]int)}
	if m.finalized {
		m.failAt(s.Name(), 0, ErrFinalized, nil, nil, "Section(%s) after Finalize", k)
		return s
	}
	m.byKind[k] = s
	m.sections = append(m.sections, s)
	return s
}

// Extern declares an imported name. Finalize refuses a reference to a name
// neither defined as a symbol nor declared here.
func (m *Module) Extern(name string) {
	if m.finalized {
		m.failAt("", 0, ErrFinalized, nil, nil, "Extern(%q) after Finalize", name)
		return
	}
	if m.err != nil {
		return
	}
	m.externs[name] = true
}

// Sections returns every section in creation order.
func (m *Module) Sections() []*Section {
	out := make([]*Section, len(m.sections))
	copy(out, m.sections)
	return out
}

// Err returns the sticky error, if any, without finalizing — for callers
// that want to bail out of a long build early rather than discover the
// failure at Finalize. It is the same error Finalize will surface.
func (m *Module) Err() error { return m.err }

// failAt records the first error with its position. Errors are sticky:
// every builder call after the first failure is a no-op, so the error a
// caller sees is the first thing that went wrong, not the last symptom
// of it.
func (m *Module) failAt(section string, offset int, sentinel, cause error, notes []string, format string, args ...any) {
	if m.err != nil {
		return
	}
	m.err = &Error{
		Sentinel: sentinel,
		Cause:    cause,
		Section:  section,
		Offset:   offset,
		Context:  fmt.Sprintf(format, args...),
		Notes:    notes,
	}
}

// Finalize surfaces the sticky error, patches every same-section label
// reference into the bytes, closes symbol sizes, and verifies that every
// remaining reference targets a defined symbol or a declared Extern. The
// module is immutable afterwards, whether or not an error is returned.
func (m *Module) Finalize() error {
	if m.finalized {
		return m.err
	}
	m.finalized = true
	if m.err != nil {
		return m.err
	}

	// Symbols are module-global; a name defined twice is refused here even
	// though each section's labels are a namespace of their own.
	defined := make(map[string]bool)
	for _, s := range m.sections {
		for i := range s.syms {
			name := s.syms[i].Name
			if defined[name] {
				m.failAt(s.Name(), s.syms[i].Offset, ErrDuplicate, nil, nil,
					"symbol %q defined in more than one section", name)
				return m.err
			}
			defined[name] = true
		}
	}

	for _, s := range m.sections {
		if !s.patchLabels() {
			return m.err
		}
	}

	for _, s := range m.sections {
		s.closeSizes()
	}

	for _, s := range m.sections {
		for _, r := range s.refs {
			if !defined[r.Sym] && !m.externs[r.Sym] {
				m.failAt(s.Name(), r.Offset, ErrUndefined, nil, nil,
					"reference to %q, neither defined nor declared Extern", r.Sym)
				return m.err
			}
		}
	}
	return m.err
}