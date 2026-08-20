# asm/x86_64

The AMD64 assembler module: registers at every width, the ISA tables,
the legacy/VEX/EVEX encoder, and an in-memory module of finished
machine code. For the design — the stability gradient, the module
model, what is deliberately absent, and how to test against it — see
the [repo README](../README.md). This file is the arch-specific
surface only.

```go
import "github.com/vertex-language/asm/x86_64"
```

---

## Quick start

```go
m := x86_64.NewModule()

t := m.Section(x86_64.Text)
t.Align(16)
t.Label("_start", x86_64.Global, x86_64.Func)
t.LeaR64M(x86_64.RDI, x86_64.RIPRel(x86_64.Ref("msg", x86_64.RefPC32)))
t.CallRef(x86_64.Ref("puts", x86_64.RefPLT))     // e8 + 4-byte hole + reference
t.XorR32RM32(x86_64.EAX, x86_64.EAX)             // 31 c0 — decided now
t.Ret()

r := m.Section(x86_64.ROData)
r.Label("msg", x86_64.Local, x86_64.Object)
r.Asciz("hello\n")

m.Extern("puts")                                 // declared import

if err := m.Finalize(); err != nil {
    log.Fatal(err)
}
```

After `Finalize`, the module is immutable, pure data:

```go
for _, s := range m.Sections() {
    s.Kind()     // x86_64.Text, x86_64.ROData, x86_64.Data, x86_64.BSS
    s.Name()     // ".text"
    s.Bytes()    // []byte — finished machine code, same-section labels patched
    s.Refs()     // []x86_64.Reference — the holes a linker fills
    s.Symbols()  // []x86_64.Symbol — name, offset, size, binding, type
}
```

One section per kind: `m.Section(x86_64.Text)` returns the same
section every time, and `Sections()` lists them in creation order.

Note the first instruction above: RIP-relative addressing is how
position-independent code reaches static data on this target — one
`lea`, one reference, no thunk. This is the single biggest ergonomic
difference from i386, where PIC is a GOTOFF instruction sequence your
lowering has to spell out.

---

## Layout

```
asm/x86_64/
├── module.go            Module, NewModule, Section dispatch, Finalize
├── section.go           Section: Align, Label, data calls, place
├── inst.go              helper machinery: form lookup, the RM types,
│                        EVEX option functions, branch/call helpers
├── helpers_*_gen.go     typed helpers — one method per form, generated
├── emit.go              Emit(mnemonic, ops...) — runtime-resolved escape hatch
├── ref.go               RefKind (RefPC32, RefPLT, RefGOTPCREL, ...), Reference
├── symbol.go            Symbol, Binding, SymbolType, Visibility
├── error.go             Error + sentinels; sticky, first-error-wins
├── alias.go             re-exports: x86_64.RAX, x86_64.Mem64, x86_64.V2, ...
│
├── feature/             extension vocabulary, x86-64-v1..v4 levels, gating
├── reg/                 Reg8..Reg64, Sreg, St, Mm, Xmm, Ymm, Zmm, K, Tmm, Cr, Dr
├── operand/             Imm, M8..M512, RIP-relative forms, Label, SymRef
│
└── internal/
    ├── isa/             the form table: classes, opcodes, gates, Resolve
    └── encode/          form + operands → bytes + fixups; Nops for Align
```

Import discipline inside the tree: **nothing imports upward.** reg/
imports nothing; operand/ imports reg/; isa/ imports both; encode/
imports all three; none of them import the root. That is why encode/
produces fixups and not relocations, and why isa/ matches its own
`Arg` type rather than the root's operand values. The root sees
everything; nothing sees the root. isa/ and encode/ additionally sit
under internal/, so that boundary is Go-enforced rather than just
convention: nothing outside this module can import them, whether or
not it would have respected the layering.

One divergence from i386, in machinery, not surface: the typed helpers
here are **generated**, from the same table everything else reads.
i386's table is a few hundred forms and rarely moves; this target's is
thousands of forms across three encodings — SSE through AVX-512 — so
the generator earns its keep, and `Form.GoName()` is the single source
of what each helper is called. Binding is the same in both packages: a
helper finds its form **by GoName lookup, at the first call** — so
appending table rows breaks nothing, and a removed or renamed form
panics loudly, naming the missing form, rather than silently binding
to the wrong row.

---

## Instructions

### Typed helpers — the primary surface

One method per declared form, named `MnemonicClassClass`:

```go
t.MovR64Imm64(x86_64.RAX, 0x123456789abc)
t.MovRM64R64(x86_64.Mem64(x86_64.RBX).Disp(8), x86_64.RSI)
t.AddR64RM64(x86_64.RAX, x86_64.Mem64(x86_64.RSP))
t.VaddpsYmmYmmYmmM256(x86_64.YMM0, x86_64.YMM1, x86_64.Mem256(x86_64.RDI))
t.JeLabel("done")
t.CallRef(x86_64.Ref("puts", x86_64.RefPLT))
```

What the compiler checks, precisely: **register, immediate, memory-only
and label slots are static.** `MovR64Imm64(x86_64.EAX, …)` does not
build — `EAX` is a `reg.Reg32`, and the parameter is a `reg.Reg64`.
The one slot that is not static is **r/m** — the slot that accepts a
register *or* memory. Go has no union value type, and nothing below
the root may implement a root-declared interface, so `RM64` and its
siblings (`RM8..RM64`, `XmmM128`, `YmmM256`, `ZmmM512`, `KM64`,
`MmM64`, `Memory`) are documented `any` aliases. The class check for
that slot runs at the call: a mismatch is a sticky, positioned
`ErrForm` at `Finalize`, not a compile error. Everything else — which
is most parameters — is a red squiggle.

A helper pins its form. `MovR64Imm64` is the ten-byte imm64 and
nothing quietly relaxes it to the seven-byte `MOV r/m64, imm32`; if
you want the short one when the value fits, that is your lowering's
if-statement or `Emit`'s tie-break, stated where you can see it.
Immediate parameters are plain `int64`s — the helper pinned the field
width, so the only question left is whether the value fits it, and
that is `encode/`'s range check, surfaced as `ErrRange`. Gated helpers
still gate: an AVX-512 helper on a V2 module fails at the call with
`ErrFeature`, naming the gate.

Naming conventions — i386's, plus what this target adds:

- **Fixed operands are in the name, not the parameters.**
  `AddRAXImm32(60)`, `ShlRM64CL(...)`, the literal `1` of `ShlRM64One`.
- **Branch targets split by where they resolve.** `JmpLabel("loop")`
  is same-section, patched at `Finalize`, no relocation.
  `JmpRef(...)` leaves the module and survives into `Refs()`.
  `CallLabel` is the bare-label call to a function compiled into this
  module — no reference, no relocation; `CallRef` is the one that
  crosses.
- **`Short` pins rel8**; the plain conditional and jump names pin
  rel32. No relaxation between them — a short branch to a far target
  is `ErrRange` at `Finalize`.
- **Conditions are the canonical spellings** the table declares —
  `JeLabel`, `JneLabel`, `JbLabel`, ... — plus the aliases the world
  writes: `JzLabel` is `JeLabel` and `JnzLabel` is `JneLabel`, the
  same rows under the other name.
- **EVEX modifiers ride as options, masks as operands.** A writemask
  names a register, so it is a `reg.K` parameter on the masked helper.
  Zeroing, broadcast, SAE and rounding are one bit each with no
  register behind them, so they are option arguments:

  ```go
  t.VaddpsZmmKZmmZmmM512(ZMM0, K1, ZMM1, mem, x86_64.Zeroing())
  t.VaddpsZmmKZmmZmmM512(ZMM0, K1, ZMM1, ZMM2, x86_64.RoundNearest())
  ```

  An option a form does not take, `{z}` without a mask, rounding off
  512 bits — `encode/` refuses each by name, surfaced as `ErrEncoding`.

### Emit — the escape hatch

```go
t.Emit("mov", x86_64.RAX, x86_64.Imm(60))
t.Emit("add", x86_64.RAX, 1)                // bare Go integers coerce to Imm
```

Runtime form resolution: shortest legal encoding among matching forms,
ties broken by table order. `add rax, 1` gets the sign-extended imm8
group form; `add eax, imm32` gets the one-byte-shorter accumulator
form — which is the whole reason those rows exist. Bare Go integer
types are accepted where an immediate is meant, because a
mnemonic-as-data caller usually has values, not operand types. Note
that `Emit` with a label will happily pick a rel8 jump because it is
shortest, and pay for it at `Finalize` if the target is far.

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

`Align(n)` pads a code section with the canonical multi-byte nop
sequences — the ones GNU as emits, because the differential suite
compares padded sections byte for byte — and a data section with
zeros. `n` must be a power of two; anything else is `ErrAlign`.

---

## References

```go
x86_64.Ref("puts", x86_64.RefPLT)           // call through the PLT
x86_64.Ref("msg", x86_64.RefPC32)           // rip-relative data
x86_64.Ref("ext_var", x86_64.RefGOTPCREL)   // address via the GOT, rip-relative
x86_64.Ref("tls_var", x86_64.RefTPOFF)      // initial-exec TLS
x86_64.Ref("puts", x86_64.RefNone)          // the site's use picks
```

```
RefAbs64  RefAbs32  RefAbs32S  RefAbs16  RefAbs8    absolute
RefPC32   RefPC16   RefPC8                          PC-relative
RefPLT                                              call via the PLT
RefGOTPCREL                                         GOT entry, rip-relative
RefTLSGD  RefTLSLD  RefGOTTPOFF  RefTPOFF  RefDTPOFF32   TLS models
```

The kinds are this target's, not i386's, and the differences are the
architecture's: `RefAbs32S` exists because a 32-bit absolute in a
64-bit address space is two different claims — zero-extended or
sign-extended — and the linker range-checks them differently.
`RefGOTPCREL` replaces the whole GOTOFF/GOTPC family, because on this
target the GOT is reached rip-relatively like everything else. And the
default for static data is `RefPC32` through `RIPRel`, not an
absolute — position-independence is the native idiom here, not a
dance.

`RefNone` means the site's use picks: a call gets `RefPLT`, a
rip-relative load gets `RefPC32`, an absolute field gets the `RefAbs`
of its width. Naming a kind overrides that — and asking for one is
asking for a relocation, so a `Ref(...)` never folds into the bytes,
even when the name happens to resolve in the same section. Only the
`Label` helpers fold. One consequence worth knowing: an *absolute*
field holding a same-section address still needs a relocation — no
address exists at this layer — and a relocation needs a symbol, so an
absolute reference to a bare label is refused at `Finalize` with
`ErrUndefined` and a note telling you to promote it.

```go
type Reference struct {
    Offset int      // where the hole starts, section-relative
    Size   int      // 1, 2, 4 or 8
    PCRel  bool
    Adjust int32    // field-position correction, already computed
    Sym    string
    Kind   RefKind
    Addend int64    // logical addend, never adjusted for the field
}
```

`Adjust` is the reason you never write `-4`. A rel32 field is relative
to the *end of the instruction*, but the field is not always the last
thing in it: a call's displacement ends its instruction, so the
correction is zero, while the displacement of
`mov dword [rip+x], 5` is followed by four bytes of immediate, so the
correction is −4 and the raw ELF addend comes out −8. The encoder
knows because it placed the field; the downstream lowering just adds.

---

## Features

```go
m := x86_64.NewModule(x86_64.WithFeatures(
    x86_64.V2.Plus(x86_64.AVX512F),
))
```

Default is `Baseline`, which is **V1**: SSE2 and nothing above it. The
psABI's CMOV, CX8, FPU, FXSR, OSFXSR, SCE and OSXSAVE bits report CPU
identification or OS enablement, not that an encoding exists, so they
are absent. The microarchitecture levels V1 ⊂ V2 ⊂ V3 ⊂ V4 are
shorthand for closed sets, exactly as the psABI states them; gating
happens against features, and levels only ever appear in spelling.

A gating failure is `ErrFeature`, and its note line prints a Go
expression you can paste back into your build:

```
vpaddd requires avx512vl, not in the active feature set
  active: x86-64-v2
  note: enable with x86_64.V2.Plus(x86_64.AVX512VL)
```

Every level and extension constant is re-exported at the root for
exactly this reason: the note line has to name something that
compiles at your call site.

`ParseFeatures` accepts the spellings the world writes —
`"x86-64-v3"`, `"x86-64-v4-avx512vl"`, `"sse2+aes"` — and `String()`
round-trips through it. A spelling without a leading level starts
*empty*, not at Baseline: "sse2" means sse2, because silently widening
the set would make a gating diagnostic unfalsifiable.

---

## Errors

The common sentinels per the repo README, plus one this architecture
demands:

```
ErrFeature  ErrForm  ErrDuplicate  ErrUndefined
ErrRange    ErrAlign  ErrFinalized

ErrEncoding    a legal-looking combination the silicon refuses:
               AH with a REX prefix, XMM16+ outside EVEX, {z} without
               a mask, rounding control off 512 bits
```

This architecture has corners where operands that each exist cannot
coexist — AH/CH/DH/BH occupy the encodings REX reassigns to
SPL/BPL/SIL/DIL, vector registers 16–31 have no bits outside EVEX,
embedded rounding steals the length field — and each refusal names the
rule, not just the failure.

The concrete type is `*x86_64.Error`, with `Section`, `Offset`,
`Context` and `Notes` fields; `errors.Is` works against every
sentinel. Errors are sticky and first-wins: every builder call after a
failure is a no-op, and `Finalize` surfaces the first one, positioned.
`Module.Err()` exists for callers that want to bail early.

---

## Lowering

```go
lower.X8664ELF(w, am, lower.ELFOptions{GNUStack: elf.StackNonExec})
lower.X8664COFF(w, am, lower.COFFOptions{})
lower.X8664MachO(w, am, lower.MachOOptions{})
```