package aarch64

import "github.com/vertex-language/asm/aarch64/operand"

// RelocKind is an arc identifier, not a format's number: R_AARCH64_ABS64 the
// ELF number is 257, the IMAGE_REL and ARM64_RELOC equivalents are other
// numbers entirely. The values here are names; the kind-to-number tables live
// in the lowerings. Naming a kind on a Ref is a request that blocks folding.
type RelocKind = operand.RelocKind

const RelocNone = operand.RelocNone

const (
	R_AARCH64_ABS64 RelocKind = iota + 1
	R_AARCH64_ABS32
	R_AARCH64_ABS16
	R_AARCH64_PREL32
	R_AARCH64_CALL26
	R_AARCH64_JUMP26
	R_AARCH64_CONDBR19
	R_AARCH64_TSTBR14
	R_AARCH64_ADR_PREL_LO21
	R_AARCH64_ADR_PREL_PG_HI21
	R_AARCH64_ADD_ABS_LO12_NC
	R_AARCH64_LDST8_ABS_LO12_NC
	R_AARCH64_LDST16_ABS_LO12_NC
	R_AARCH64_LDST32_ABS_LO12_NC
	R_AARCH64_LDST64_ABS_LO12_NC
	R_AARCH64_LDST128_ABS_LO12_NC
	R_AARCH64_ADR_GOT_PAGE
	R_AARCH64_LD64_GOT_LO12_NC
)

// AddrRole and Width are the operand package's, re-exported because Reference
// carries them across the wall.
type (
	AddrRole = operand.AddrRole
	Width    = operand.Width
)

// Reference is a hole a linker fills. No Size — every hole is a field inside
// one four-byte word. No Adjust — every pc-relative relocation here is defined
// against the instruction's own address.
type Reference struct {
	Offset int      // the instruction's offset, section-relative
	Sym    string
	Role   AddrRole // direct, page, pageoff, gotpage, gotpageoff
	Kind   RelocKind // insisted-on kind, or RelocNone
	Addend int64    // logical addend, never field-corrected
	Access Width    // memory access width, for the LDST lo12 family
	Branch bool     // a branch field rather than a data field
}