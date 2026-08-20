# asm/i386

The Intel386 assembler module: registers, the ISA table, the encoder,
and an in-memory module of finished machine code. For the design — the
stability gradient, the module model, what is deliberately absent, and
how to test against it — see the [repo README](../README.md). This file
is the arch-specific surface only.

```go
import "github.com/vertex-language/asm/i386"
```

---

## Quick start

```go
m := i386.NewModule()

t := m.Section(i386.Text)
t.Align(16)
t.Label("_start", i386.Global, i386.Func)
t.MovR32Imm32(i386.EAX, 4)                      // b8 04 00 00 00 — decided now
t.MovR32Imm32(i386.EBX, 1)
t.CallRef(i386.Ref("puts", i386.RefPLT))        // e8 + 4-byte hole + reference
t.Ret()

r := m.Section(i386.ROData)
r.Label("msg", i386.Local, i386.Object)
r.Asciz("hello\n")

m.Extern("puts")                                // declared import

if err := m.Finalize(); err != nil {
    log.Fatal(err)
}
```

After `Finalize`, the module is immutable, pure data:

```go
for _, s := range m.Sections() {
    s.Kind()     // i386.Text, i386.Data, i386.ROData, i386.BSS
    s.Name()     // ".text"
    s.Bytes()    // []byte — finished machine code, same-section labels patched
    s.Refs()     // []i386.Reference — the holes a linker fills
    s.Symbols()  // []i386.Symbol — name, offset, size, binding, type
}
```

One section per kind: `m.Section(i386.Text)` returns the same section
every time, and `Sections()` lists them in creation order. The freeze
is total: after `Finalize`, existing sections still come back for
reading, but every builder call is a no-op recording `ErrFinalized`,
and asking for a kind that was never created is refused the same way —
`Sections()` is as frozen as the bytes are.

---

## Layout

```
asm/i386/
├── module.go        Module, NewModule, Section dispatch, Err, Finalize
├── section.go       Section: Align, Label, Offset, data calls, place
├── inst.go          typed helpers — one method per form, bound by name
├── emit.go          Emit(mnemonic, ops...) — runtime-resolved escape hatch
├── ref.go           Reference; RefKind aliased from operand
├── symbol.go        Symbol, Binding, SymbolType, Visibility, SymbolAttr
├── error.go         Error + sentinels; sticky, first-error-wins
├── alias.go         re-exports: i386.EAX, i386.Mem32, i386.Ref, ...
│
├── reg/             R8, R16, R32, Sreg, St, Mm, Xmm, Ymm, Zmm, K, Cr, Dr
├── operand/         Imm, M8..M512, Label, SymRef, RelocKind
├── feature/         I386..I686 ladder, MMX..AVX512 extensions, Parse
└── internal/
    ├── isa/         the form table: classes, opcodes, gates, Resolve
    └── encode/      form + operands → bytes + fixups; Nops for Align
```

Import discipline: **nothing imports upward.** reg imports nothing;
operand imports reg; feature stands alone; isa imports reg, operand,
and feature; encode imports all of them; none of them import the root.
The root sees everything; nothing sees the root.

isa and encode are `internal/` because they are implementation: the
typed helpers and `Emit` are the instruction surface, and nothing a
caller writes holds an isa or encode type. reg, operand, and feature
are public because caller code holds their types in its own signatures
— a lowering that passes registers around has `reg.R32` parameters.
What internal packages produce still reaches you: encoder failures
surface through `*i386.Error` as a sentinel, a message, and notes —
data, not internal types.

There is no code generation step. `inst.go` is written by hand against
the internal table and binds each helper to its form **by name lookup,
not by table index** — appending rows breaks nothing here, and a
removed or renamed form panics at the first call, naming the missing
form, rather than silently binding to the wrong row. The typed surface
is a few hundred one-line methods over a table that rarely moves; a
generator would be machinery to keep in sync with the thing it
generates. If the table ever churns weekly, revisit — `asm/x86_64` is
the revisit.

---

## Instructions

### Typed helpers — the primary surface

One method per declared form, named `MnemonicClassClass`:

```go
t.MovR32Imm32(i386.EAX, 60)
t.MovRM32R32(i386.Mem32(i386.EBX).Disp(8), i386.ESI)
t.AddR32RM32(i386.EAX, i386.Mem32(i386.ESP))
t.JzLabel("done")
t.CallRef(i386.Ref("puts", i386.RefPLT))
```

The operand classes are the parameter types, so a width or class
mismatch is a compile error: `MovR32Imm32(i386.AL, …)` does not build —
`AL` is a `reg.R8`. An isel bug is a red squiggle rather than a runtime
`ErrForm`.

A helper pins its form — exactly the encoding you named — and the
diagnostics follow from that. A helper checks operand *kinds* only and
hands values to the encoder: the wrong kind of operand is `ErrForm`; a
value that does not fit the field the form pins is `ErrRange`, with the
field width and range in the error's notes and the encoder's error as
the cause. You named the form, so a too-big constant is a value
problem, not a missing-form problem. Gated helpers still gate: a gated
helper on a module without the extension fails at the call with
`ErrFeature`, naming the gate.

Naming conventions beyond the class spelling:

- **Fixed operands are in the name, not the parameters.**
  `AddEAXImm32(60)` — the form names EAX and leaves no field to put
  another register in. Same for `AL`, `CL`, `DX`, and the literal `1`
  of `ShlRM32One`.
- **Branch targets split by where they resolve.** `JmpLabel("loop")`
  is same-section, patched at `Finalize`, no relocation.
  `JmpRef(i386.Ref("f", i386.RefPC32))` leaves the section and
  survives into `Refs()`. `CallLabel` is the bare-label call to a
  function compiled into this module — no reference, no relocation;
  `CallRef` is the one that crosses.
- **`Short` pins rel8.** `JmpShortLabel` and `JzShortLabel` are the
  2-byte forms; the plain names pin rel32. No relaxation between them —
  a short branch to a far target is `ErrRange` at `Finalize`. Mnemonics
  with only a rel8 form (`LoopLabel`, `JecxzLabel`) take the plain
  name.

### Emit — the escape hatch

```go
t.Emit("mov", i386.EAX, i386.Imm(60))
```

Runtime form resolution: shortest legal encoding among matching forms,
ties broken by table order. `add eax, 1` gets the four-byte
sign-extended imm8 form over the six-byte imm32 form — which is the
whole reason those rows exist. Exists for table-driven emission where
the mnemonic is data; if you know the instruction at compile time, use
the typed helper. Note that `Emit` with a label will happily pick a
rel8 jump because it is shortest, and pay for it at `Finalize` if the
target is far.

### Data

```go
r.Byte(0x90)
r.Long(0xdeadbeef)         // little-endian, 4 bytes
r.Quad(1)                  // little-endian, 8 bytes
r.Ascii("no terminator")
r.Asciz("terminated")
r.Zero(64)
r.Data(blob)               // raw bytes
```

The raw-bytes builder is `Data` because `Bytes()` is the read side of
the contract; there is no `Word` because a name whose width depends on
the package you're in is not a name.

`Align(n)` pads a code section with the arch's multi-byte nop
sequences — gated by the module's feature set, because below i686 the
0F 1F nop is an illegal instruction, not a slower one — and a data
section with zeros. `n` must be a power of two; anything else is
`ErrAlign`.

`Offset()` is the current end of the section: the offset the next byte
will land at, the value a `Label` placed now would name. It is exported
because jump tables and literal pools are yours to build — this package
deliberately refuses to build them — and building them requires knowing
where you are.

---

## Memory operands

Three constructors per access width, one question each:

```go
i386.Mem32(i386.EBX)                        // based:    [ebx]
i386.Abs32(i386.Ref("msg", i386.RefAbs32))  // symbolic: [msg], a relocation
i386.Addr32(0xb8000)                        // direct:   [0xb8000], no relocation
```

Chain methods refine the address and keep the width's type, so the
result still satisfies the helper parameter it is written into:

```go
i386.Mem32(i386.EBX).Disp(8)
i386.Mem32(i386.EBP).Index(i386.ESI, 4).Disp(-12)
i386.Addr32(0).Index(i386.EDI, 8)           // index-only
i386.Mem32(i386.EBX).Seg(i386.FS)
```

Construction errors — ESP as an index (SIB has no encoding for it), a
scale that is not 1/2/4/8 — are sticky on the operand and surface at
the instruction that uses it, positioned, so a builder chain is not
followed by a run of error checks.

---

## References

```go
i386.Ref("puts", i386.RefPLT)          // call through the PLT
i386.Ref("msg", i386.RefAbs32)         // absolute 32-bit address
i386.Ref("tls_var", i386.RefTLSIE)     // initial-exec TLS
```

```
RefAbs32  RefAbs16  RefAbs8              absolute, that many bytes
RefPC32   RefPC16   RefPC8               PC-relative
RefPLT                                   call via procedure linkage table
RefGOT    RefGOTOFF  RefGOTPC            global offset table forms
RefTLSGD  RefTLSLDM  RefTLSIE  RefTLSLE  RefDTPOFF  RefTPOFF
```

The kind is stated at construction because it is a decision — PLT or
direct, GOT or absolute — and this package does not make decisions.
`RefKind` is defined once, in `operand`, and aliased at the root: one
list, no conversion anywhere in the tree.

The idiom to know: loading `msg`'s address from `.text` is
cross-section, so it is always a reference — `RefAbs32` in a non-PIC
build, which is fully linkable for executables and the right place to
start. PIC on i386 (no PC-relative addressing) means the GOTOFF/GOTPC
dance with an EIP thunk — an *instruction sequence* your lowering
emits, not something this package can hide.

```go
type Reference struct {
    Offset int      // where the hole starts, section-relative
    Size   int      // 1, 2, or 4
    PCRel  bool
    Adjust int32    // field-position correction, already computed
    Sym    string
    Kind   RefKind
    Addend int32
}
```

`Adjust` is the reason you never write `-4`. A PC-relative field is
resolved against the end of the instruction, and the encoder that
placed the field knows how many bytes follow it; the downstream
lowering just adds.

---

## Features

```go
m := i386.NewModule(i386.WithFeatures(
    i386.NewFeatureSet(i386.I686).Add(i386.SSE2),
))
```

Default is `I686` with no extensions. Two different things live in the
vocabulary and are deliberately not flattened together: a **level** is
a point on the cumulative i386..i686 base-CPU ladder, and a **feature**
is an orthogonal extension with a CPUID bit of its own. Sets are closed
under requirements, as everywhere in this repo: adding SSE2 brings SSE,
removing SSE drops SSE2.

`feature.Parse` resolves the spellings the world writes, against a
starting set:

```go
feature.Parse(feature.Default(), "i586+mmx")     // exact: level plus extensions
feature.Parse(base, "+avx2,-sse4.2")             // adjust: applied left to right
```

`Set.String()` prints the canonical spelling — `i686+sse2+aes` — and
`Parse` accepts it back, so the two round-trip. Spellings that name
something real but wrong get the reason, not "unknown": `cmov` is part
of the i686 base level; `cx16` requires 64-bit mode; `mpx` was removed
from the architecture; the x86-64 microarchitecture levels are defined
over the 64-bit baseline and have no 32-bit member.

---

## Errors

The common sentinels, per the repo README:

```
ErrFeature  ErrForm  ErrDuplicate  ErrUndefined
ErrRange    ErrAlign  ErrFinalized
```

`ErrRange` covers both value failures: a too-big immediate at the call,
and a label displacement that does not fit its pinned field at
`Finalize`. Both carry the field width and reachable range in the
error's notes. There is no branch relaxation and no silent form
substitution — the failure is loud instead of the bytes being
different.

The concrete type is `*i386.Error`, with `Section`, `Offset`,
`Context`, and `Notes` fields. Errors are sticky and first-wins: every
builder call after a failure is a no-op, and `Finalize` surfaces the
first one, positioned. `Module.Err()` exists for callers that want to
bail early; it returns the same error `Finalize` will.

`errors.Is` works against every sentinel. Where an encoder or operand
error is the cause, it joins the chain — `Unwrap` exposes both the
sentinel and the cause — but its type is internal; anything a caller
might need from it is restated in `Notes` as text. Internal packages
do not become API through the error chain.

---

## Lowering

Not this package. A downstream consumer imports this package and an
object library and transcribes `Sections()`/`Symbols()`/`Refs()` into
`AddSection`/`AddSymbol`/`AddReloc` calls:

```go
lower.I386ELF(w, am, lower.ELFOptions{GNUStack: elf.StackNonExec})
lower.I386COFF(w, am, lower.COFFOptions{})   // '_' mangling lives there
```

The RefKind-to-relocation-number table lives there — `RefPLT` →
`R_386_PLT32`, or a refusal by name for a format with no PLT concept —
because format numbers are the format layer's vocabulary and this
package carries semantics.