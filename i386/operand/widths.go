package operand

import "github.com/vertex-language/asm/i386/reg"

// One type per access width. The chain methods return their own type so a
// builder expression stays typed — Mem32(EBX).Disp(8) is an M32 and
// satisfies RM32 — which is why they are spelled out per width rather than
// shared on the embedded core: a method on mem would return mem and drop
// the width on the floor.
//
// Three constructors per width:
//
//	Mem32(EBX)              based:    [ebx]
//	Abs32(Ref("msg", k))    symbolic: [msg], a relocation
//	Addr32(0xb8000)         direct:   [0xb8000], no relocation
//
// Index-only addressing is Addr32(0).Index(ESI, 4).

type (
	M8   struct{ mem }
	M16  struct{ mem }
	M32  struct{ mem }
	M64  struct{ mem }
	M80  struct{ mem }
	M128 struct{ mem }
	M256 struct{ mem }
	M512 struct{ mem }
)

func (M8) Bits() int   { return 8 }
func (M16) Bits() int  { return 16 }
func (M32) Bits() int  { return 32 }
func (M64) Bits() int  { return 64 }
func (M80) Bits() int  { return 80 }
func (M128) Bits() int { return 128 }
func (M256) Bits() int { return 256 }
func (M512) Bits() int { return 512 }

// The r/m markers, mirroring reg's. A width with no register half (M80) or
// no i386 form yet (M256, M512) carries no marker; it is still a Memory.
func (M8) RM8()     {}
func (M16) RM16()   {}
func (M32) RM32()   {}
func (M64) RM64()   {}
func (M128) RM128() {}

func Mem8(base reg.R32) M8     { return M8{based(base)} }
func Mem16(base reg.R32) M16   { return M16{based(base)} }
func Mem32(base reg.R32) M32   { return M32{based(base)} }
func Mem64(base reg.R32) M64   { return M64{based(base)} }
func Mem80(base reg.R32) M80   { return M80{based(base)} }
func Mem128(base reg.R32) M128 { return M128{based(base)} }
func Mem256(base reg.R32) M256 { return M256{based(base)} }
func Mem512(base reg.R32) M512 { return M512{based(base)} }

func Abs8(sym SymRef) M8     { return M8{symbolic(sym)} }
func Abs16(sym SymRef) M16   { return M16{symbolic(sym)} }
func Abs32(sym SymRef) M32   { return M32{symbolic(sym)} }
func Abs64(sym SymRef) M64   { return M64{symbolic(sym)} }
func Abs80(sym SymRef) M80   { return M80{symbolic(sym)} }
func Abs128(sym SymRef) M128 { return M128{symbolic(sym)} }
func Abs256(sym SymRef) M256 { return M256{symbolic(sym)} }
func Abs512(sym SymRef) M512 { return M512{symbolic(sym)} }

func Addr8(a int32) M8     { return M8{direct(a)} }
func Addr16(a int32) M16   { return M16{direct(a)} }
func Addr32(a int32) M32   { return M32{direct(a)} }
func Addr64(a int32) M64   { return M64{direct(a)} }
func Addr80(a int32) M80   { return M80{direct(a)} }
func Addr128(a int32) M128 { return M128{direct(a)} }
func Addr256(a int32) M256 { return M256{direct(a)} }
func Addr512(a int32) M512 { return M512{direct(a)} }

func (m M8) Disp(n int32) M8                 { m.mem = m.mem.withDisp(n); return m }
func (m M8) Index(r reg.R32, scale int) M8   { m.mem = m.mem.withIndex(r, scale); return m }
func (m M8) Seg(s reg.Sreg) M8               { m.mem = m.mem.withSeg(s); return m }
func (m M8) Sym(r SymRef) M8                 { m.mem = m.mem.withSym(r); return m }

func (m M16) Disp(n int32) M16               { m.mem = m.mem.withDisp(n); return m }
func (m M16) Index(r reg.R32, scale int) M16 { m.mem = m.mem.withIndex(r, scale); return m }
func (m M16) Seg(s reg.Sreg) M16             { m.mem = m.mem.withSeg(s); return m }
func (m M16) Sym(r SymRef) M16               { m.mem = m.mem.withSym(r); return m }

func (m M32) Disp(n int32) M32               { m.mem = m.mem.withDisp(n); return m }
func (m M32) Index(r reg.R32, scale int) M32 { m.mem = m.mem.withIndex(r, scale); return m }
func (m M32) Seg(s reg.Sreg) M32             { m.mem = m.mem.withSeg(s); return m }
func (m M32) Sym(r SymRef) M32               { m.mem = m.mem.withSym(r); return m }

func (m M64) Disp(n int32) M64               { m.mem = m.mem.withDisp(n); return m }
func (m M64) Index(r reg.R32, scale int) M64 { m.mem = m.mem.withIndex(r, scale); return m }
func (m M64) Seg(s reg.Sreg) M64             { m.mem = m.mem.withSeg(s); return m }
func (m M64) Sym(r SymRef) M64               { m.mem = m.mem.withSym(r); return m }

func (m M80) Disp(n int32) M80               { m.mem = m.mem.withDisp(n); return m }
func (m M80) Index(r reg.R32, scale int) M80 { m.mem = m.mem.withIndex(r, scale); return m }
func (m M80) Seg(s reg.Sreg) M80             { m.mem = m.mem.withSeg(s); return m }
func (m M80) Sym(r SymRef) M80               { m.mem = m.mem.withSym(r); return m }

func (m M128) Disp(n int32) M128               { m.mem = m.mem.withDisp(n); return m }
func (m M128) Index(r reg.R32, scale int) M128 { m.mem = m.mem.withIndex(r, scale); return m }
func (m M128) Seg(s reg.Sreg) M128             { m.mem = m.mem.withSeg(s); return m }
func (m M128) Sym(r SymRef) M128               { m.mem = m.mem.withSym(r); return m }

func (m M256) Disp(n int32) M256               { m.mem = m.mem.withDisp(n); return m }
func (m M256) Index(r reg.R32, scale int) M256 { m.mem = m.mem.withIndex(r, scale); return m }
func (m M256) Seg(s reg.Sreg) M256             { m.mem = m.mem.withSeg(s); return m }
func (m M256) Sym(r SymRef) M256               { m.mem = m.mem.withSym(r); return m }

func (m M512) Disp(n int32) M512               { m.mem = m.mem.withDisp(n); return m }
func (m M512) Index(r reg.R32, scale int) M512 { m.mem = m.mem.withIndex(r, scale); return m }
func (m M512) Seg(s reg.Sreg) M512             { m.mem = m.mem.withSeg(s); return m }
func (m M512) Sym(r SymRef) M512               { m.mem = m.mem.withSym(r); return m }