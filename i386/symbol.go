package i386

// Binding is a symbol's link-time binding.
type Binding uint8

const (
	Local Binding = iota
	Global
	Weak
)

var bindingNames = [...]string{"local", "global", "weak"}

func (b Binding) String() string {
	if int(b) < len(bindingNames) {
		return bindingNames[b]
	}
	return "?"
}

// SymbolType is a symbol's link-time type.
type SymbolType uint8

const (
	NoType SymbolType = iota
	Func
	Object
	ThreadLocal
)

var symbolTypeNames = [...]string{"notype", "func", "object", "tls"}

func (t SymbolType) String() string {
	if int(t) < len(symbolTypeNames) {
		return symbolTypeNames[t]
	}
	return "?"
}

// Visibility is carried verbatim; whether a format can express it is the
// lowering's question and refusal.
type Visibility uint8

const (
	DefaultVisibility Visibility = iota
	Hidden
	Protected
	Internal
)

var visibilityNames = [...]string{"default", "hidden", "protected", "internal"}

func (v Visibility) String() string {
	if int(v) < len(visibilityNames) {
		return visibilityNames[v]
	}
	return "?"
}

// SymbolAttr is anything Label accepts as an attribute: a Binding, a
// SymbolType, or a Visibility. Passing any attribute promotes a label into
// Symbols().
type SymbolAttr interface {
	applySymbol(*Symbol)
}

func (b Binding) applySymbol(s *Symbol)    { s.Binding = b }
func (t SymbolType) applySymbol(s *Symbol) { s.Type = t }
func (v Visibility) applySymbol(s *Symbol) { s.Visibility = v }

// Symbol is an attributed label: a link-time fact with no cell in the bytes.
// Size is closed at Finalize: the next symbol in the section, or section end.
type Symbol struct {
	Name       string
	Offset     int
	Size       int
	Binding    Binding
	Type       SymbolType
	Visibility Visibility
}