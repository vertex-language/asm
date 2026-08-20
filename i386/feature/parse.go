package feature

import (
	"errors"
	"fmt"
	"strings"
)

// ErrParse is the sentinel for a spelling that names nothing this target
// has. The message states the fact; how it is presented — usage error,
// diagnostic, log line — is the caller's business, because this package
// does not know a command line exists.
var ErrParse = errors.New("i386 features")

// aliases are alternate spellings accepted for extensions. They resolve
// here and nowhere below: past Parse only canonical values exist.
var aliases = map[string]Feature{
	"sse4":      SSE42,
	"sse4a":     SSE42,
	"bmi":       BMI1,
	"pclmulqdq": PCLMUL,
	"aesni":     AES,
	"abm":       LZCNT,
	"f16":       F16C,
	"sha_ni":    SHA,
	"avx512":    AVX512F,
	"avx512_f":  AVX512F,
}

// levelAliases are alternate spellings for the base levels.
var levelAliases = map[string]Level{
	"80386":      I386,
	"80486":      I486,
	"pentium":    I586,
	"i786":       I686,
	"pentiumpro": I686,
	"p6":         I686,
}

// inBaseline names spellings that are real x86 features but are part of a
// base level here, so asking for them as extensions gets an answer rather
// than "unknown".
var inBaseline = map[string]Level{
	"cmov":      I686,
	"fcmov":     I686,
	"fcomi":     I686,
	"cx8":       I586,
	"cmpxchg8b": I586,
	"cpuid":     I586,
	"rdtsc":     I586,
	"bswap":     I486,
	"cmpxchg":   I486,
	"xadd":      I486,
	"fpu":       I386,
	"x87":       I386,
}

// notThirtyTwoBit names extensions that exist but require 64-bit mode, so
// the diagnostic says why rather than claiming the name is unknown.
var notThirtyTwoBit = map[string]string{
	"cx16":       "CMPXCHG16B requires 64-bit mode",
	"cmpxchg16b": "CMPXCHG16B requires 64-bit mode",
	"sce":        "SYSCALL/SYSRET are 64-bit mode instructions",
	"syscall":    "SYSCALL/SYSRET are 64-bit mode instructions",
	"lahf_lm":    "LAHF/SAHF in 64-bit mode is the extension; in 32-bit mode they are baseline",
	"amx-tile":   "AMX requires 64-bit mode",
	"amx-int8":   "AMX requires 64-bit mode",
	"amx-bf16":   "AMX requires 64-bit mode",
	"uintr":      "user interrupts require 64-bit mode",
}

// withdrawn names extensions that shipped and were then removed from the
// architecture. Only current, ratified extensions are encodable, so these
// are errors with a reason rather than silent acceptance.
var withdrawn = map[string]string{
	"mpx":     "Intel MPX was removed from the architecture",
	"bnd":     "Intel MPX was removed from the architecture",
	"3dnow":   "AMD 3DNow! was removed from the architecture",
	"3dnowa":  "AMD 3DNow! was removed from the architecture",
	"pcommit": "PCOMMIT was removed from the architecture",
}

// Parse resolves a feature spelling against a starting set.
//
// Two grammars, which do not mix:
//
//	i686+sse2,aes      exact — start from the named level (or base's level)
//	                   with no extensions, and add each named one
//	+avx2,-sse4.2      adjust — start from base, apply left to right
//
// ',' and '+' both separate tokens, so the canonical String() spelling —
// "i686+sse2" — parses back to the set that printed it. A leading '+' or
// '-' selects the adjust grammar; a level may not appear there, because
// adding or removing a point on a cumulative ladder has no meaning.
func Parse(base Set, s string) (Set, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return base, nil
	}

	// The x86-64 microarchitecture levels contain '-' and would tokenize
	// into nonsense; recognize the whole family first and refuse it with
	// the reason. The levels are defined over the 64-bit baseline — v1
	// includes SCE and v2 requires CMPXCHG16B — so none has a 32-bit
	// member.
	if strings.Contains(strings.ToLower(s), "x86-64") {
		return base, fmt.Errorf("%w: %q is an x86-64 microarchitecture level and does not apply to i386: the base levels here are i386, i486, i586, i686", ErrParse, s)
	}

	adjust := s[0] == '+' || s[0] == '-'

	type token struct {
		name   string
		remove bool
	}
	var toks []token
	for i := 0; i < len(s); {
		for i < len(s) && (s[i] == ',' || s[i] == '+' || s[i] == ' ') {
			i++
		}
		if i >= len(s) {
			break
		}
		remove := false
		if s[i] == '-' {
			remove = true
			i++
		}
		start := i
		for i < len(s) && s[i] != ',' && s[i] != '+' && s[i] != '-' && s[i] != ' ' {
			i++
		}
		if start == i {
			return base, fmt.Errorf("%w: empty name in %q", ErrParse, s)
		}
		toks = append(toks, token{strings.ToLower(s[start:i]), remove})
	}

	out := base
	var levelSet bool
	if !adjust {
		out = New(base.level)
	}

	for _, t := range toks {
		if t.remove && !adjust {
			return base, fmt.Errorf("%w: %q removes from an exact set: an exact spelling names what is present; removal is the +/- grammar", ErrParse, t.name)
		}

		if l, ok := parseLevel(t.name); ok {
			if adjust {
				return base, fmt.Errorf("%w: %q is a base level, not an extension: levels are cumulative and cannot be added or removed", ErrParse, t.name)
			}
			if levelSet && l != out.level {
				return base, fmt.Errorf("%w: two base levels named: %s and %s", ErrParse, out.level, l)
			}
			out.level, levelSet = l, true
			continue
		}

		f, err := parseFeature(t.name)
		if err != nil {
			return base, err
		}
		if t.remove {
			out = out.Remove(f)
		} else {
			out = out.Add(f)
		}
	}
	return out, nil
}

func parseLevel(name string) (Level, bool) {
	for i, n := range levelNames {
		if name == n {
			return Level(i), true
		}
	}
	l, ok := levelAliases[name]
	return l, ok
}

func parseFeature(name string) (Feature, error) {
	for f := Feature(0); f < numFeatures; f++ {
		if name == featureNames[f] {
			return f, nil
		}
	}
	if f, ok := aliases[name]; ok {
		return f, nil
	}
	if l, ok := inBaseline[name]; ok {
		return 0, fmt.Errorf("%w: %q is not an extension: it is part of the %s base level", ErrParse, name, l)
	}
	if why, ok := notThirtyTwoBit[name]; ok {
		return 0, fmt.Errorf("%w: %q does not apply to i386: %s", ErrParse, name, why)
	}
	if why, ok := withdrawn[name]; ok {
		return 0, fmt.Errorf("%w: %q is not accepted: %s", ErrParse, name, why)
	}
	return 0, fmt.Errorf("%w: unknown extension %q for i386", ErrParse, name)
}

// Introspection. These replace a preformatted help text: the data is the
// library's, the formatting is the consumer's.

// Levels returns the base levels in ladder order.
func Levels() []Level { return []Level{I386, I486, I586, I686} }

// All returns every extension in canonical order.
func All() []Feature {
	out := make([]Feature, 0, numFeatures)
	for f := Feature(0); f < numFeatures; f++ {
		out = append(out, f)
	}
	return out
}

// Requires returns what f directly requires, in canonical order. The
// closure is what Add takes; this is the single edge, for anything that
// wants to print or walk the graph.
func Requires(f Feature) []Feature {
	if f >= numFeatures {
		return nil
	}
	out := make([]Feature, len(requires[f]))
	copy(out, requires[f])
	return out
}