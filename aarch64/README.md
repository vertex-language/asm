# asm/aarch64

The ARM64 assembler module: the register files, the A64 instruction
table, the fixed-width encoder, and an in-memory module of finished
machine code. For the design — the stability gradient, the module
model, what is deliberately absent, and how to test against it — see
the [repo README](../README.md). This file is the arch-specific
surface only.

```go
import "github.com/vertex-language/asm/aarch64"
```

---

## Quick start

```go
m := aarch64.NewModule()

t := m.Section(aarch64.Text)
t.Align(16)
t.Label("_start", aarch64.Global, aarch64.Func)
t.Adrp(aarch64.X0, aarch64.Ref("msg"))                       // page of the address
t.AddImm64(aarch64.X0, aarch64.X0,
    aarch64.PageOff(aarch64.Ref("msg")))                     // offset within it
t.BlRef(aarch64.Ref("puts"))                                 // 94xxxxxx + reference
t.MovzImm64(aarch64.X0, 0)                                   // d2800000 — decided now
t.Ret()                                                      // operand omitted: X30

r := m.Section(aarch64.ROData)
r.Label("msg", aarch64.Local, aarch64.Object)
r.Asciz("hello\n")

m.Extern("puts")

if err := m.Finalize(); err != nil {
    log.Fatal(err)
}
```

After `Finalize`, the module is immutable, pure data:

```go
for _, s := range m.Sections() {
    s.Kind()     // aarch64.Text, aarch64.ROData, aarch64.Data, aarch64.BSS
    s.Name()     // ".text"
    s.Bytes()    // []byte — finished machine code, same-section labels patched
    s.Refs()     // []aarch64.Reference — the holes a linker fills
    s.Symbols()  // []aarch64.Symbol — name, offset, size, binding, type
}
```

One section per kind, listed in creation order, as everywhere in this
repo.

Note the first two instructions. Materializing an address on this
target is usually **two instructions and therefore two references** —
`adrp` for the 4KiB page, `add` (or a load) for the offset within it —
and each carries a *role* naming which half of the address it wants.
A bare `Ref` in the ADRP slot is read as the page, because the form is
the page form; the ADD's immediate has no such form to lean on, so its
role must be written — `PageOff(...)` — and a bare symbol there is
refused. This is the idiom to internalize, and it sits between the
other two targets: position-independence is native, as on x86_64, but
it costs an instruction pair your lowering spells out, as on i386 —
the difference being that here the pair is the architecture's own
design rather than a GOT dance, and this package carries both halves
as data rather than making you fake them.

---

## Layout

```
asm/aarch64/
├── module.go            Module, NewModule, Section dispatch, Finalize
├── section.go           Section: Align, Label, data calls, resolve/fold
├── inst.go              helper machinery: form lookup by GoName, place,
│                        the sentinel mapping over encode/'s typed errors
├── helpers_base_gen.go  typed helpers — one method per form, generated
├── emit.go              Emit(mnemonic, ops...) — runtime-resolved escape
│                        hatch — and Inst(word)
├── reloc.go             RelocKind constants (R_AARCH64_*, arc-named), Reference
├── symbol.go            Symbol, Binding, SymbolType, Visibility
├── error.go             Error + sentinels; sticky, first-error-wins
├── alias.go             re-exports: aarch64.X0, aarch64.SP, aarch64.Page, ...
│
├── feature/             armv8-a..armv9.5-a levels, extension vocabulary, gating
├── reg/                 X, W, Xsp, Wsp, V..B scalar views, Vec, VLane, Z, P, FFR, Sys
├── operand/             Imm, Mem, Shift, Extend, Cond, Barrier, PrfOp, roles
│
└── internal/
    ├── isa/             the form table: classes, base words, masks, gates, Resolve
    └── encode/           form + operands → one word + fixups; the nop word
```

Import discipline as everywhere: reg/ imports nothing, operand/
imports reg/, isa/ imports both, encode/ imports all three, none of
them import the root. The helpers are generated from the table and
bound **by GoName lookup at the first call**, exactly as in x86_64 —
appended rows break nothing, a removed or renamed form panics loudly by
name. isa/ and encode/ additionally sit under internal/, so that
boundary is Go-enforced rather than just convention: nothing outside
this module can import them, whether or not it would have respected
the layering.

Three things live in `reg/` that the other packages have no analogue
for, because the architecture demands them: register number 31 is two
different registers (`X`/`W` read it as the zero register, `Xsp`/`Wsp`
as the stack pointer, and `WithSP`/`WithZR` are the explicit, fallible
conversions between the readings); AAPCS64 preservation is a package
function `Save(r, variant)` rather than a method, because whether Z8 is
callee-saved depends on the calling function's interface and no
property of Z8 decides it; and `DWARF(r)` reports whether a number
*exists* — XZR has none, because the zero register is not a location.

### The table has no tiebreak

Every A64 instruction is one 32-bit word, so there is no shortest-form
search and nothing for table order to decide. Two forms of one mnemonic
that accept the same operand classes are an **ambiguity refused when
the table is built** — a panic at init, not an arbitrary answer at
encode time. That is the whole difference in kind between this
package's `Resolve` and x86_64's, and it is why `Emit` here never
surprises you the way `Emit`-picks-rel8 can there.

---

## Instructions

### Typed helpers — the primary surface

One method per declared form, named as the table names it — the
ARM ARM's form title, with the datasize spelled out:

```go
t.SubImm64(aarch64.SP, aarch64.SP, 32)
t.AddShifted64(aarch64.X0, aarch64.X1, aarch64.X2, aarch64.Shifted(aarch64.LSL, 3))
t.LdrImm64(aarch64.X0, aarch64.Mem64(aarch64.SP).Off(16))
t.StpPre64(aarch64.X29, aarch64.X30, aarch64.Mem64(aarch64.SP).Pre(-16))
t.Csel64(aarch64.X0, aarch64.X1, aarch64.X2, aarch64.EQ)
t.CbzLabel64(aarch64.X0, "done")
t.Ret()
```

The 32- and 64-bit variants are separate helpers taking different
register types, so a width mismatch is a compile error:
`AddShifted64(aarch64.W0, …)` does not build — `W0` is a `reg.W`.

**The one slot that is not static is the register-31 slot.** A slot
that reads 31 as the stack pointer accepts two Go types — `SP` (a
`reg.Xsp`) or a numbered `reg.X` — and Go has no union value type, so,
exactly like x86_64's `RM64` family, `RegSP64` and `RegSP32` are
documented `any` aliases and the check runs at the call. Handing such a
slot `XZR` is a sticky, positioned `ErrForm` whose note says which
register 31 the slot reads — not "wrong register", but "this slot reads
register 31 as the stack pointer; xzr is a different register that
shares the encoding." Every other parameter is a red squiggle.

Optional operands are trailing variadics: `Ret()` or `Ret(aarch64.X1)`,
a shift written or omitted. `NoShift` exists for callers building
argument lists at runtime.

Naming and shape conventions, beyond what i386 and x86_64 establish:

- **A negative ADD immediate is not a SUB.** `AddImm64` takes the
  twelve unsigned bits (optionally `LSL #12`) the encoding has, and a
  value outside them is `ErrRange`. Rewriting `add #-32` into
  `sub #32` would be substituting an instruction, which nothing here
  does. Your lowering owns that if-statement.
- **Optional operands are the architecture's, stated in the encoding.**
  `Ret()` is RET with its operand omitted, which the encoding defines
  as X30 — a default, not an alias. An omitted shift is `LSL #0` the
  same way.
- **Branch targets split by where they resolve**, as everywhere:
  `BLabel`/`BlLabel` are same-section, patched at `Finalize`, no
  relocation; `BRef`/`BlRef` survive into `Refs()`. The conditional
  branch is one form with a condition operand — `BCondLabel(aarch64.NE,
  "loop")` — not a helper per condition, because that is how the table
  declares it. `HS` and `LO` are accepted spellings of `CS` and `CC`,
  the same encodings under the names source writes.
- **Shifts and extends are operands, not mnemonic suffixes.**
  `add x0, x1, x2, lsl #3` is `AddShifted64` with a `ShiftOp` argument,
  the same form with a field filled in.
- **Address materialization is spelled as roles.** `Page(...)`,
  `PageOff(...)`, `GotPage(...)`, `GotPageOff(...)` wrap a target and
  name which half of its address the operand wants. A symbolic value
  may land in two places: a memory displacement
  (`Mem64(X1).Off(PageOff(Ref("msg")))`) or a data-processing immediate
  (`AddImm64(X0, X0, PageOff(Ref("msg")))`), and in both only the
  page-offset roles are accepted. A bare symbol in either is refused at
  construction or at the call: the low-twelve half of an address is a
  fact the reader needs written down, and this package does not invent
  it. The one inference is ADRP's, where a bare `Ref` means the page
  because the form is the page form.

### The architecture's aliases — and only those

`Cmp`, `Cmn`, `Tst`, `Neg`, `MovReg64`, `MovSp64`, `Cset64`, `Mul64`,
the immediate shifts — each is in the ARM ARM's own alias table, each
is one-to-one with a word of the instruction it aliases, and each
carries the ARM ARM's own preferred-disassembly rule. When you call
`Cmp` you get the SUBS word with Rd pinned to 31: the same instruction
under another name, not a different instruction chosen for you.

The assembler-invented layer is deliberately absent, and this is worth
a table row of its own because every other AArch64 assembler has it:

| absent | why |
|---|---|
| `ldr x0, =value` | One mnemonic becoming a literal pool is selection plus layout. Emit the load and the pool yourself, or materialize the constant. |
| MOVZ/MOVK chains | A constant spanning two halfwords is refused with the range error saying so. Two instructions are the caller's to write, because *which* two is a decision. |
| `mov x0, #imm` dispatch | Whether a constant is a move-wide or a logical immediate changes which instruction is emitted. `MovzImm64` and `OrrImm64` are both there; picking is your lowering's. |

### Emit — the escape hatch

```go
t.Emit("add", aarch64.X0, aarch64.X1, 42)
t.Emit("add", aarch64.X0, aarch64.X0, aarch64.PageOff(aarch64.Ref("msg")))
t.Emit("ldr", aarch64.X0, aarch64.Mem64(aarch64.SP).Off(16))
t.Emit("b", aarch64.Label("done"))
```

Runtime resolution against the same table — and because the table
refuses ambiguity at build time, `Emit` matches **exactly one form or
fails**, naming the near-miss candidates. A form that matches but is
gated returns the gating error rather than "no such form": being told
an instruction does not exist when it exists and is disabled sends a
reader hunting for a typo that is not there.

A label target is spelled `aarch64.Label("done")`. A bare string is
refused, because a bare name could be a label or a register, and
guessing which is how a typo becomes an object file. Bare Go integers
coerce to immediates, because a mnemonic-as-data caller usually has
values, not operand types.

One layering fact worth knowing, because a diagnostic will eventually
show it to you: `Resolve` decides which *form* an operand lands in, so
an immediate slot accepts a symbolic reference there — and encode/ is
what then refuses every role but the page-offset ones. So
`Emit("svc", aarch64.PageOff(...))` resolves and fails in the encoder,
positioned, naming the operand. Resolve answers "which form"; encode/
answers "may this value fill it" — the same division the memory path
has, where `[x1, #8]` matches both LDR's and LDUR's shape and the
mnemonic you named is what picks.

### Data

```go
r.Byte(0x90)
r.Long(0xdeadbeef)         // little-endian, 4 bytes
r.Quad(1)
r.Ascii("no terminator")
r.Asciz("terminated")
r.Zero(64)
r.Data(blob)
```

`Align(n)` pads a code section with `d503201f` — there is one nop and
it is one word, so there is no multi-byte-sequence table to match GNU
as against — and a data section with zeros. `n` must be a power of two,
and in a code section the padding must be whole instructions; a text
alignment that would strand a partial word is `ErrAlign`, not rounded.

There is a `.inst` equivalent: `Inst(word)` states a whole word rather
than naming an instruction. It goes through the table's own row, and it
is the one place where emitting bytes nobody selected is exactly what
was asked for.

---

## References

The other two packages carry a flat `RefKind`. This one carries a
**role**, because the unit of "an address" here is usually two
instructions:

```go
aarch64.Ref("msg")                        // a symbol; the site's use picks
aarch64.Page(aarch64.Ref("msg"))          // :pg_hi21: — the ADRP half
aarch64.PageOff(aarch64.Ref("msg"))       // :lo12:    — the ADD/load half
aarch64.GotPage(aarch64.Ref("ext_var"))   // :got:     — GOT slot, page
aarch64.GotPageOff(aarch64.Ref("ext_var"))// :got_lo12:
aarch64.Ref("puts", aarch64.R_AARCH64_CALL26)  // a kind the caller insists on
```

The role is the portable fact — GNU as spells the pair `:pg_hi21:` /
`:lo12:`, Darwin spells it `@PAGE` / `@PAGEOFF`, and the relocation
number each becomes is the lowering's table. The caller states the
role, never the number. Naming a `RelocKind` is a request that blocks
folding, exactly as a named kind does in the other packages: asking for
`CALL26` to a symbol two lines up still gets a relocation, because
resolving it would answer a different question than the one asked.

Folding is narrower here than a first read suggests: **only a
`RoleDirect` reference to a same-section bare label folds** — a branch,
an `adr`, a literal load. A `Page` or `PageOff` reference to a bare
label is refused at `Finalize` with a note telling you to promote it,
because a page delta depends on where the section finally loads, and
nothing at this layer assigns an address. This is the same fact as
x86_64's absolute-field rule, seen from the other side.

One consequence with no analogue elsewhere: the low-twelve relocation
of a *load* depends on the access width — under ELF it is one of five
`LDST{8,16,32,64,128}_ABS_LO12_NC` kinds, chosen by what the immediate
is scaled by. The caller never states that width at the reference; the
form knows it, and the encoder copies it onto the reference. That
single field is what the fixup/relocation distinction pays for on this
target.

```go
type Reference struct {
    Offset int                 // the instruction's offset, section-relative
    Sym    string
    Role   aarch64.AddrRole    // direct, page, pageoff, gotpage, gotpageoff
    Kind   aarch64.RelocKind   // insisted-on kind, or RelocNone
    Addend int64               // logical addend, never field-corrected
    Access aarch64.Width       // memory access width, for the LDST lo12 family
    Branch bool                // a branch field rather than a data field
}
```

Two fields from the x86_64 struct are absent, and their absence is the
architecture's gift: there is no `Size`, because every hole is a field
inside one four-byte word, and there is no `Adjust`, because every
PC-relative relocation here is defined against the address of the
instruction itself — the field is never followed by bytes it must
correct for, so the entire `-4`/`-8` addend table collapses to nothing.

---

## Features

```go
m := aarch64.NewModule(aarch64.WithFeatures(
    aarch64.Armv9A.Plus(aarch64.SVE2BitPerm),
))
```

Default is `Baseline()`: Armv8-A, floating point and Advanced SIMD and
nothing above them.

The levels are names for closed sets, as everywhere — but the ladder
here is **not a chain and not monotone**, because the architecture's
isn't. Armv9-A is built on Armv8.5-A, not Armv8.9-A, so `Armv8_8A` has
MOPS and `Armv9A` does not and neither contains the other. And CRC is
optional at Armv8-A, mandatory from 8.1, so `Armv8A.Plus(CRC)` and
`Armv8_1A` are different overlapping sets, each printing as what the
world calls it: `armv8-a+crc` and `armv8.1-a`.

Every level and extension constant is re-exported at the root, because
a gating failure names the fix as a Go expression that compiles at your
call site:

```
sqrdmlah requires rdma, not in the active feature set
  active: armv8-a
  note: aarch64.WithFeatures(aarch64.Armv8A.Plus(aarch64.RDMA))
```

One spelling collision to know: the floating-point feature is
`aarch64.FeatFP` at the root, because `aarch64.FP` is the frame
pointer, X29 — a register the AAPCS names and this package prints
first. Inside `feature/` it remains `feature.FP`, and the printed
spelling is `fp` either way.

`ParseFeatures` accepts the `+ext` / `+noext` grammar GCC and gas
accept — `"armv8.2-a+sve+nofp16"`, `"armv9-a+sme2"`, `"+lse+crc"` —
applied left to right, each closed as applied, so `+sve2` pulls in SVE
and a later `+nosve` drops both. A spelling without a leading version
starts *empty*, not at Baseline: "sve" means SVE and its requirements
alone. `"crypto"` is accepted as the group it is and never printed
back, because nothing is encoded as "crypto".

---

## Errors

The common sentinels per the repo README, plus one this architecture
demands:

```
ErrFeature  ErrForm  ErrDuplicate  ErrUndefined
ErrRange    ErrAlign  ErrFinalized

ErrBitmask     a constant that is not a logical immediate
```

It is separate from `ErrRange` because the fix is different in kind. A
value out of range gets smaller; a value that is not a rotated run of
ones replicated across the register does not get more expressible, and
the caller has to materialize it into a register instead — which the
diagnostic says, naming movz/movk and the register form. All-zeros and
all-ones get their own sentences, because they are the two constants a
caller most plausibly writes and the encoding can name neither: the run
of ones must be neither empty nor complete.

The register-31 refusals surface under `ErrForm` with the note naming
the rule — which reading of 31 the slot has, and that the register you
passed is a different one sharing the encoding — rather than a bare
mismatch, because "wrong register" would send you counting parameters
when the count is fine.

The concrete type is `*aarch64.Error`, with `Section`, `Offset`,
`Context` and `Notes`. Its `Unwrap` exposes both the sentinel and the
underlying resolver or encoder error, so `errors.Is` works against
every sentinel and `errors.As` reaches the typed cause — an
`*encode.RangeError` still names its operand after crossing the wall.
Errors are sticky and first-wins, surfaced at `Finalize`, with
`Module.Err()` for callers that want to bail early.

---

## Lowering

```go
lower.AArch64ELF(w, am, lower.ELFOptions{GNUStack: elf.StackNonExec})
lower.AArch64MachO(w, am, lower.MachOOptions{})   // @PAGE spellings live there
lower.AArch64COFF(w, am, lower.COFFOptions{})
```

The role-to-kind tables live there: `RolePage` in a branch-less field →
`R_AARCH64_ADR_PREL_PG_HI21` or `ARM64_RELOC_PAGE21`;
`RolePageOff` in a load → the `LDST*` kind its `Access` selects, or
`ARM64_RELOC_PAGEOFF12`. A role a format cannot express is refused by
name at the lowering, not silently approximated — the same rule as
everywhere, and on this target the one Mach-O consumers hit first.