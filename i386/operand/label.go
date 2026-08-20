package operand

import "github.com/vertex-language/asm/i386/reg"

// Label is a reference to a name in this section's namespace. It is patched
// at Finalize and leaves no trace in Refs(); a reference that must survive
// into a relocation is a SymRef.
type Label struct {
	reg.Seal
	name string
}

func NewLabel(name string) Label { return Label{name: name} }

func (l Label) Name() string   { return l.name }
func (l Label) Bits() int      { return 0 }
func (l Label) String() string { return l.name }