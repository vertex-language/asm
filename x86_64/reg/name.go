// x86_64/reg/name.go
package reg

import "strconv"

var (
	name64 = [16]string{"rax", "rcx", "rdx", "rbx", "rsp", "rbp", "rsi", "rdi",
		"r8", "r9", "r10", "r11", "r12", "r13", "r14", "r15"}
	name32 = [16]string{"eax", "ecx", "edx", "ebx", "esp", "ebp", "esi", "edi",
		"r8d", "r9d", "r10d", "r11d", "r12d", "r13d", "r14d", "r15d"}
	name16 = [16]string{"ax", "cx", "dx", "bx", "sp", "bp", "si", "di",
		"r8w", "r9w", "r10w", "r11w", "r12w", "r13w", "r14w", "r15w"}
	name8 = [20]string{"al", "cl", "dl", "bl", "spl", "bpl", "sil", "dil",
		"r8b", "r9b", "r10b", "r11b", "r12b", "r13b", "r14b", "r15b",
		"ah", "ch", "dh", "bh"}
	nameSreg = [6]string{"es", "cs", "ss", "ds", "fs", "gs"}
)

func (r Reg64) Name() string { return name64[r] }
func (r Reg32) Name() string { return name32[r] }
func (r Reg16) Name() string { return name16[r] }
func (r Reg8) Name() string  { return name8[r] }
func (r Sreg) Name() string  { return nameSreg[r] }

func (r Xmm) Name() string { return "xmm" + strconv.Itoa(int(r)) }
func (r Ymm) Name() string { return "ymm" + strconv.Itoa(int(r)) }
func (r Zmm) Name() string { return "zmm" + strconv.Itoa(int(r)) }
func (r St) Name() string  { return "st" + strconv.Itoa(int(r)) }
func (r Mm) Name() string  { return "mm" + strconv.Itoa(int(r)) }
func (r K) Name() string   { return "k" + strconv.Itoa(int(r)) }
func (r Tmm) Name() string { return "tmm" + strconv.Itoa(int(r)) }
func (r Cr) Name() string  { return "cr" + strconv.Itoa(int(r)) }
func (r Dr) Name() string  { return "dr" + strconv.Itoa(int(r)) }

// String is Name. These are bare, lowercase names; diagnostics print them
// unadorned.
func (r Reg64) String() string { return r.Name() }
func (r Reg32) String() string { return r.Name() }
func (r Reg16) String() string { return r.Name() }
func (r Reg8) String() string  { return r.Name() }
func (r Sreg) String() string  { return r.Name() }
func (r Xmm) String() string   { return r.Name() }
func (r Ymm) String() string   { return r.Name() }
func (r Zmm) String() string   { return r.Name() }
func (r St) String() string    { return r.Name() }
func (r Mm) String() string    { return r.Name() }
func (r K) String() string     { return r.Name() }
func (r Tmm) String() string   { return r.Name() }
func (r Cr) String() string    { return r.Name() }
func (r Dr) String() string    { return r.Name() }