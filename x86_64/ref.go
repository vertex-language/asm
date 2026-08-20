// x86_64/ref.go
package x86_64

import "github.com/vertex-language/asm/x86_64/operand"

// RefKind is link semantics, not a relocation number. `call puts@plt` and
// `call puts` are byte-identical, so the kind rides beside the bytes; a
// downstream table maps each kind to R_X86_64_*, IMAGE_REL_AMD64_* or
// X86_64_RELOC_* — or refuses by name.
//
// The type is defined over operand.RelocKind so a SymRef can carry it
// without operand/ importing the root. RefNone is operand.RelocNone: "the
// site's use picks" — a call gets RefPLT, a rip-relative load gets RefPC32,
// an absolute field gets the RefAbs of its width.
type RefKind operand.RelocKind

const (
	RefNone RefKind = iota

	RefAbs64
	RefAbs32
	RefAbs32S // sign-extended: a different claim the linker range-checks differently
	RefAbs16
	RefAbs8

	RefPC32
	RefPC16
	RefPC8

	RefPLT
	RefGOTPCREL

	RefTLSGD
	RefTLSLD
	RefGOTTPOFF
	RefTPOFF
	RefDTPOFF32

	numRefKinds
)

var refKindNames = [numRefKinds]string{
	RefNone: "none",
	RefAbs64: "abs64", RefAbs32: "abs32", RefAbs32S: "abs32s",
	RefAbs16: "abs16", RefAbs8: "abs8",
	RefPC32: "pc32", RefPC16: "pc16", RefPC8: "pc8",
	RefPLT: "plt", RefGOTPCREL: "gotpcrel",
	RefTLSGD: "tlsgd", RefTLSLD: "tlsld",
	RefGOTTPOFF: "gottpoff", RefTPOFF: "tpoff", RefDTPOFF32: "dtpoff32",
}

func (k RefKind) String() string {
	if int(k) >= len(refKindNames) {
		return "refkind(?)"
	}
	return refKindNames[k]
}

// Ref builds a symbol reference. RefNone lets the site's use pick the kind.
func Ref(name string, kind RefKind) operand.SymRef {
	return operand.Ref(name, operand.RelocKind(kind))
}

// Reference is one hole a linker fills: everything a downstream format needs
// to build a relocation record, in terms that name no format.
type Reference struct {
	Offset int   // where the hole starts, section-relative
	Size   int   // 1, 2, 4 or 8
	PCRel  bool
	Adjust int32 // field-position correction, already computed (-Tail)
	Sym    string
	Kind   RefKind
	Addend int64 // logical addend, never adjusted for the field's position
}