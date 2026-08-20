package i386

import "github.com/vertex-language/asm/i386/operand"

// RefKind is link semantics, not a relocation number. The values are
// defined once, in operand — the lowest package every layer can see — and
// aliased here, so `call puts@plt` and `call puts` differ by exactly one
// datum that rides beside identical bytes, and no conversion exists
// anywhere in the tree. A downstream lowering maps RefPLT to R_386_PLT32
// for ELF, or refuses by name for a format with no PLT concept.
type RefKind = operand.RelocKind

const (
	RefAbs32 = operand.RefAbs32
	RefAbs16 = operand.RefAbs16
	RefAbs8  = operand.RefAbs8

	RefPC32 = operand.RefPC32
	RefPC16 = operand.RefPC16
	RefPC8  = operand.RefPC8

	RefPLT = operand.RefPLT

	RefGOT    = operand.RefGOT
	RefGOTOFF = operand.RefGOTOFF
	RefGOTPC  = operand.RefGOTPC

	RefTLSGD  = operand.RefTLSGD
	RefTLSLDM = operand.RefTLSLDM
	RefTLSIE  = operand.RefTLSIE
	RefTLSLE  = operand.RefTLSLE
	RefDTPOFF = operand.RefDTPOFF
	RefTPOFF  = operand.RefTPOFF
)

// Reference is a hole a linker fills: where it is, how wide, how it
// resolves. Adjust is the field-position correction, already computed by
// the encoder that placed the field.
type Reference struct {
	Offset int
	Size   int
	PCRel  bool
	Adjust int32
	Sym    string
	Kind   RefKind
	Addend int32
}