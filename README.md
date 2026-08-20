# asm

In-memory assembler modules for native code generation, one package per
architecture:

```
asm/
├── i386/        Intel386: R32 registers, legacy encoding, GOTOFF PIC
├── x86_64/      AMD64: legacy/VEX/EVEX, RIP-relative PIC, x86-64-v1..v4
└── aarch64/     ARM64: fixed-width A64 encoding, page+offset PIC, armv8-a..armv9.5-a
```

```
go get github.com/vertex-language/asm
```

Each package produces the same thing: **decided bytes plus the three
facts bytes cannot carry** — references, symbols, section kinds. None
of them write object files, and none of them know your IR exists.

Zero dependencies outward. These packages import nothing from
`objectfile`, `linker`, each other, or any compiler. They are the
bottom of a stack, and the bottom imports nothing.

---

## Why this repo exists: the stability gradient

Order a native-code pipeline by how often each layer changes;
dependencies must point down the gradient — volatile code depends on
stable code, never the reverse:

```
your IR                     changes constantly — it is where you design
isel · regalloc · frame     rewritten whenever the IR's shape changes
────────────────────────── the wall ──────────────────────────────────
asm/<arch>  (this repo)     terms that don't move: RAX, mov, .text, PLT
objectfile                  frozen file formats: ELF, COFF, Mach-O
linker
```

Intel defined `EAX` in 1985; AMD defined `RAX` in 2003; the encoding of
`mov` is a published specification; the psABI's relocation semantics
are frozen documents. An API built only from these inherits their
stability.

The failure mode this repo is designed against: an encoder built around
an IR gets rewritten every time the IR changes, because it is built
against the IR's shape. The rule, stated once:

> **Don't build your machine code encoder around an IR. Build your IR —
> and its lowering — around a machine code package.**

The measurable consequence: when your IR changes, your lowering package
rewrites and these packages do not lose a line. The helpers still
compile, still emit the same bytes, still carry the same references,
and every test against `Section.Bytes()` still passes untouched.

---

## The flow

```
   ir.Module
       │
       │   isel + regalloc + frame layout
       │   ── your compiler's backend; lives in YOUR repo,
       │      imports one arch package, pattern-matches on YOUR nodes
       ▼
   <arch>.Module                      ←  this repo
       │
       │   transcription
       │   ── a ~60-line loop: Sections()/Symbols()/Refs()
       │      against a writer's AddSection/AddSymbol/AddReloc
       ▼
   objectfile (ELF / COFF / Mach-O)
       │
       ▼
   linker
```

The asymmetry is the design. **IR → Module is the hard arrow**:
choosing instructions, assigning physical registers, computing stack
offsets — all compiler work, all yours, all coupled to your IR and
rewritten with it. **Module → object is the trivial arrow**:
transcription of facts already decided. If the second arrow is ever
making a decision — picking an encoding, choosing a section, resolving
a symbol — something leaked upstream; file the bug here.

Isel and regalloc feel "backend-ish" but do not live in this repo:
their algorithms are textbook-stable, but their *implementations* are
functions over your IR's node shapes, so they move when your IR moves.
They belong on the volatile side of the wall, beside the thing they are
coupled to. This repo has never heard of them.

---

## The model

A `Module` sits **after every decision an assembler can make** and
**before every decision only a linker can make.** This is the
definition, not a description, and it is the same in every arch
package.

Decided at the call, irrevocably:

- **which instruction** — nothing substitutes, folds, or reorders later
- **which registers** — physical only; there is no virtual register
  type to pass
- **which encoding** — a typed helper pins its form; `Emit` picks the
  shortest legal one at the call and commits
- **which offset** — `len(section)` at the moment of the call

Undecidable here, carried as the minimal residue:

- **references** — `puts` has no address anywhere in the universe yet;
  the site is a zero-filled hole plus offset, width, PC-relativity,
  semantic kind, addend
- **symbols** — "global", "function" are link-time facts with no cell
  in the bytes
- **section kinds** — "read-only" is a load-time property

That closed list is the whole cargo. A section is an append-only byte
buffer and a reference list. Nothing is symbolic-and-rewritable, so no
pass can run over it — which is precisely what makes this not an IR.

The litmus test for anything you're tempted to add: *would it survive
into a flat binary?* A module whose `Refs()` is empty after `Finalize`
is a complete program as raw bytes.

### References are semantics, not numbers

A `RefKind` is **link semantics, not a relocation number**. `call
puts@plt` and `call puts` are byte-identical — `e8` either way — so the
kind must ride beside the bytes. It is the one fact in this repo that
neither the CPU nor the file format owns: your codegen decides it (PLT
or direct? GOT or absolute?), so it crosses the wall as data. A
downstream table maps each kind to its format's relocation number —
`RefPLT` → `R_386_PLT32` or `R_X86_64_PLT32` or
`IMAGE_REL_AMD64_REL32` — or refuses by name. One table, one place,
visible holes. The kinds themselves are per-architecture, because the
architectures genuinely differ; see each package's README.

### Labels, symbols, externs

Identical rules in every package:

- a **bare label** lives only in its section's namespace; a branch to
  it is patched at `Finalize` and leaves no trace
- **any attribute** (`Global`, `Func`, ...) promotes it into
  `Symbols()`; size is closed at the next symbol or section end
- symbols are **module-global**; a duplicate is `ErrDuplicate`
- a call to a function you compiled into this module is a **bare-label
  call** — no reference, no relocation. Only calls outside the module
  take `Ref(...)` and need an `Extern`; `Finalize` refuses a reference
  to a name neither defined nor declared

Bindings: `Global`, `Local`, `Weak`. Types: `Func`, `Object`,
`ThreadLocal`. Visibility (`Hidden`, `Protected`, `Internal`) is
carried verbatim; whether a format can express it is the lowering's
question and refusal.

### Finalize

Four jobs, then the module is immutable:

1. surface the sticky error — every builder call after a failure was a
   no-op; the first error is the one you get, positioned
2. patch every same-section label reference into the bytes and drop it,
   refusing with `ErrRange` any displacement that does not fit its
   pinned field
3. close every symbol's size at the next symbol or section end
4. verify every remaining reference targets a defined symbol or a
   declared `Extern`

What remains in `Refs()` is exactly the set every downstream format
must turn into relocation records — the same set regardless of format.

Decided ≠ resolved: the module commits to *what* the bytes are; the
numeric values inside label displacements and relocation fields arrive
at `Finalize` and link time respectively, because they depend on
addresses nothing at this layer assigns.

### Errors

Sticky, first-wins, surfaced at `Finalize` — a run of builder calls is
not followed by a run of `if err != nil`. Every error is positioned:
section, offset, and what was being built. `errors.Is` works against
every sentinel. The common sentinels:

```
ErrFeature     form gated behind an extension not in the module's set
ErrForm        operand combination no declared form accepts
ErrDuplicate   label defined twice in one section, or a symbol in two
ErrUndefined   reference to a name neither defined nor declared Extern
ErrRange       label displacement does not fit its pinned field
ErrAlign       alignment is not a power of two
ErrFinalized   builder call after Finalize
```

Architectures add sentinels only where the silicon demands one (see
`x86_64`'s `ErrEncoding`).

### Features

```go
m := x86_64.NewModule(x86_64.WithFeatures(...))
```

Fixed at construction — a gate that changed mid-module would make
already-emitted diagnostics unfalsifiable. The vocabulary is
instruction-set extensions only — "may I emit these bytes" — never
CPUID detection: the feature set says what you *may emit*, not what the
host has. Detecting the host is a runtime library's job; conflating
them would make cross-compilation a special case instead of the
default. Sets are always closed under requirements: adding an extension
brings its whole closure, and removing one removes everything that
requires it.

---

## What is deliberately absent, and why

Every exclusion below is load-bearing. Adding any of them back moves
this repo up the gradient, and the gradient is the point.

| absent | why |
|---|---|
| **IR, passes, optimizer** | Anything symbolic-and-rewritable invites rewriting, and rewriting couples this repo to whoever rewrites. A section is bytes; there is nothing to walk. Your optimizer runs over your IR, above the wall. |
| **virtual registers, register allocation** | A vreg type would be an IR type. Allocation is a function over *your* liveness, *your* node shapes — it moves when your IR moves, so it lives beside your IR. Every register here is physical because physical registers are the terms that don't move. |
| **instruction selection** | `Emit` picks an *encoding* of the instruction you named. It never picks a different instruction. Selection pattern-matches on IR nodes this repo has never heard of. |
| **object file writing, relocation numbers** | Format numbers (`R_*`, `IMAGE_REL_*`, `*_RELOC_*`) are the format layer's vocabulary. These packages carry semantic `RefKind`s; the mapping table lives downstream. |
| **text input: GAS/NASM parsers, directives, macros** | Go is the input language. A compiler lowers by calling methods, not by printing assembly and re-parsing it. A text layer is a second front door that must agree with the first forever. |
| **disassembly, `Decode`, `Explain`** | A tool's job. A builder only encodes; carrying the inverse machine doubles the surface for zero lowering value. |
| **branch relaxation** | Widening a pinned `rel8` to `rel32` behind your back would change bytes you decided. `Finalize` fails loudly with `ErrRange` instead. Policy belongs above the wall. |
| **CPUID detection** | See Features above. |
| **cross-arch anything** | A shared operand type that worked for two architectures would be an IR — the thing these packages are defined by not being. Each arch redeclares everything under its own types, on purpose. |
| **name mangling** | The cdecl `_` prefix is a COFF file decoration, not a fact about the code. It lives in the COFF lowering. |

---

## Lowering to an object file

Not this repo. A downstream consumer imports an arch package and an
object library and transcribes:

```go
am := x86_64.NewModule()
// ... your IR backend builds am ...
am.Finalize()

lower.X8664ELF(w, am, lower.ELFOptions{GNUStack: elf.StackNonExec})
lower.X8664COFF(w, am, lower.COFFOptions{})
lower.X8664MachO(w, am, lower.MachOOptions{})
```

Format-only decisions surface as *parameters of the lowering call* —
`e_flags`, OSABI, GNU-stack notes, bigobj policy — because they are
facts about files, not about code, and they have no safe defaults these
packages could pick for you. The same module lowers to every format; a
`RefKind` a format cannot express is refused by name at the lowering,
not silently approximated.

---

## Testing your backend against this

Every boundary is a plain value you can stop at and inspect, which is
what the wall buys you day to day:

- assert your isel by comparing `Section.Bytes()` — no file written
- diff a function's bytes against `gcc -c` output; encoder prefix order
  is pinned to GNU as for exactly this reason
- print `Refs()` and check against `readelf -r` of a reference object
- unit-test regalloc by walking the helpers it called, not the object
  it produced

Compiler backends are miserable to debug because these seams are
usually fused. Here they are not, and that — more than any purity
argument — is what the separation is for.