package reg

import "sync"

// all is every register this package declares, in declaration order — a
// deterministic enumeration for tools and tests.
//
// There is deliberately no Lookup(name): resolving "eax" from text is a
// parser's job, and this repo has no text front end by design. The name map
// below exists only to panic at first use if two registers ever collide on
// a spelling, which is a table-integrity check, not an API.
var (
	all     []Reg
	allOnce sync.Once
)

func buildAll() {
	add := func(rs ...Reg) { all = append(all, rs...) }
	for i := 0; i < 8; i++ {
		add(R32(i))
	}
	for i := 0; i < 8; i++ {
		add(R16(i))
	}
	for i := 0; i < 8; i++ {
		add(R8(i))
	}
	for i := 0; i < 6; i++ {
		add(Sreg(i))
	}
	for i := 0; i < 8; i++ {
		add(St(i), Mm(i))
	}
	for i := 0; i < 8; i++ {
		add(Xmm(i), Ymm(i), Zmm(i), K(i))
	}
	for i := 0; i < 8; i++ {
		add(Cr(i), Dr(i))
	}

	seen := make(map[string]Reg, len(all))
	for _, r := range all {
		if prev, dup := seen[r.Name()]; dup {
			panic("i386/reg: duplicate register name " + r.Name() + " (" + prev.Class().String() + ")")
		}
		seen[r.Name()] = r
	}
}

// All returns every declared register in declaration order.
func All() []Reg {
	allOnce.Do(buildAll)
	out := make([]Reg, len(all))
	copy(out, all)
	return out
}