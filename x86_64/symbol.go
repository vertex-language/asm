// x86_64/symbol.go
package x86_64

// Symbol is a promoted label: a name with link-time facts attached. Size is
// zero until Finalize closes it at the next symbol or section end.
type Symbol struct {
	Name       string
	Offset     int // section-relative
	Size       int // closed at Finalize
	Binding    Binding
	Type       SymbolType
	Visibility Visibility
}

// SymbolAttr is anything Label accepts after the name. Passing any attribute
// promotes a bare label into Symbols(); Binding, SymbolType and Visibility
// all qualify.
type SymbolAttr interface{ applyTo(*Symbol) }

// Binding is who wins at link time.
type Binding uint8

const (
	Local Binding = iota
	Global
	Weak
)

func (b Binding) applyTo(s *Symbol) { s.Binding = b }

func (b Binding) String() string {
	switch b {
	case Global:
		return "global"
	case Weak:
		return "weak"
	}
	return "local"
}

// SymbolType is what the symbol names.
type SymbolType uint8

const (
	NoType SymbolType = iota
	Func
	Object
	ThreadLocal
)

func (t SymbolType) applyTo(s *Symbol) { s.Type = t }

func (t SymbolType) String() string {
	switch t {
	case Func:
		return "func"
	case Object:
		return "object"
	case ThreadLocal:
		return "tls"
	}
	return "notype"
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

func (v Visibility) applyTo(s *Symbol) { s.Visibility = v }

func (v Visibility) String() string {
	switch v {
	case Hidden:
		return "hidden"
	case Protected:
		return "protected"
	case Internal:
		return "internal"
	}
	return "default"
}