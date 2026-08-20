package operand

import "github.com/vertex-language/asm/i386/reg"

// RelocKind is link semantics, not a relocation number. `call puts@plt` and
// `call puts` are byte-identical — e8 either way — so the kind rides beside
// the bytes as data; a downstream lowering maps RefPLT to R_386_PLT32 for
// ELF, or refuses by name for a format with no PLT concept.
//
// The zero value names no semantics. A SymRef is always built with a kind.
type RelocKind uint8

const (
	// Absolute, by field width.
	RefAbs32 RelocKind = 1 + iota
	RefAbs16
	RefAbs8

	// PC-relative, by field width.
	RefPC32
	RefPC16
	RefPC8

	// Call via the procedure linkage table.
	RefPLT

	// Global offset table forms.
	RefGOT
	RefGOTOFF
	RefGOTPC

	// Thread-local storage models.
	RefTLSGD
	RefTLSLDM
	RefTLSIE
	RefTLSLE
	RefDTPOFF
	RefTPOFF
)

var relocNames = map[RelocKind]string{
	RefAbs32: "abs32", RefAbs16: "abs16", RefAbs8: "abs8",
	RefPC32: "pc32", RefPC16: "pc16", RefPC8: "pc8",
	RefPLT: "plt",
	RefGOT: "got", RefGOTOFF: "gotoff", RefGOTPC: "gotpc",
	RefTLSGD: "tlsgd", RefTLSLDM: "tlsldm", RefTLSIE: "tlsie",
	RefTLSLE: "tlsle", RefDTPOFF: "dtpoff", RefTPOFF: "tpoff",
}

func (k RelocKind) String() string {
	if n, ok := relocNames[k]; ok {
		return n
	}
	return "?"
}

// SymRef names a symbol and how a reference to it resolves.
type SymRef struct {
	reg.Seal
	name   string
	kind   RelocKind
	addend int32
}

// Ref builds a symbol reference. The kind is stated at construction because
// it is a decision — PLT or direct, GOT or absolute — and this package does
// not make decisions.
func Ref(name string, kind RelocKind) SymRef { return SymRef{name: name, kind: kind} }

// Plus returns the reference with n added to its addend.
func (r SymRef) Plus(n int32) SymRef { r.addend += n; return r }

func (r SymRef) Name() string    { return r.name }
func (r SymRef) Kind() RelocKind { return r.kind }
func (r SymRef) Addend() int32   { return r.addend }
func (r SymRef) Bits() int       { return 0 }
func (r SymRef) String() string  { return r.name + "@" + r.kind.String() }