package aarch64

import "github.com/vertex-language/asm/aarch64/internal/encode"

// Emit is the runtime-resolved escape hatch. Because the table refuses
// ambiguity when it is built, Emit matches exactly one form or fails, naming
// the near-miss candidates; a form that matches but is gated returns the
// gating error rather than "no such form".
//
// A label target is spelled aarch64.Label("done"); a bare string is refused
// by the encoder, because a bare name could be a label or a register.
func (s *Section) Emit(mnem string, ops ...any) {
	if !s.ready("emit " + mnem) {
		return
	}
	word, fixups, err := encode.EncodeWith(
		encode.Opts{Offset: len(s.buf)}, s.m.features, mnem, ops...)
	if err != nil {
		s.fail(sentinelFor(err), err, mnem)
		return
	}
	var b [4]byte
	b[0], b[1], b[2], b[3] = byte(word), byte(word>>8), byte(word>>16), byte(word>>24)
	s.buf = append(s.buf, b[:]...)
	s.fixups = append(s.fixups, fixups...)
}

// Inst states a whole word rather than naming an instruction — the .inst
// equivalent, routed through the table's own row.
func (s *Section) Inst(word uint32) {
	s.inst("Inst", uint64(word))
}