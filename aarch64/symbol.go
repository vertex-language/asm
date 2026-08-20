package aarch64

// SymbolAttr is anything Label accepts beyond the name. A bare label lives in
// its section's namespace; any attribute promotes it into Symbols().
type SymbolAttr interface{ apply(*Symbol) }

type Binding uint8

const (
	Local Binding = iota
	Global
	Weak
)

func (b Binding) apply(s *Symbol) { s.Binding = b }

type SymbolType uint8

const (
	NoType SymbolType = iota
	Func
	Object
	ThreadLocal
)

func (t SymbolType) apply(s *Symbol) { s.Type = t }

// Visibility is carried verbatim; whether a format can express it is the
// lowering's question and refusal.
type Visibility uint8

const (
	DefaultVisibility Visibility = iota
	Hidden
	Protected
	Internal
)

func (v Visibility) apply(s *Symbol) { s.Visibility = v }

// Symbol is a promoted label. Size is closed at Finalize: the next symbol's
// offset or the section's end.
type Symbol struct {
	Name       string
	Offset     int
	Size       int
	Binding    Binding
	Type       SymbolType
	Visibility Visibility
}