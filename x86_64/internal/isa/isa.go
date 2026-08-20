package isa

import (
	"fmt"
	"sort"
	"strings"

	"github.com/vertex-language/asm/x86_64/feature"
)

// forms is every declared form, in table order. Order is load-bearing:
// Resolve breaks size ties by it, so a form declared earlier wins.
var forms []*Form

// byMnemonic indexes forms by mnemonic, preserving table order within each.
var byMnemonic map[string][]*Form

// The tables are registered from one place rather than from an init per
// file, because Go orders per-file init by filename and the table order is
// part of the contract. table_avx.go must not sort ahead of table_base.go.
func init() {
	register(baseForms())
	register(sseForms())
	register(avxForms())
	register(avx512Forms())

	byMnemonic = make(map[string][]*Form, 512)
	for _, f := range forms {
		byMnemonic[f.Op] = append(byMnemonic[f.Op], f)
	}
}

func register(fs []*Form) {
	for _, f := range fs {
		f.finish() // derives what can be derived, panics on what cannot be encoded
		f.index = len(forms)
		forms = append(forms, f)
	}
}

func All() []*Form { return forms }

// Forms is every declared encoding of one mnemonic, in table order,
// regardless of what is enabled. Nil for a mnemonic this target has no form
// for; that is the caller's "unknown mnemonic".
func Forms(mnemonic string) []*Form { return byMnemonic[mnemonic] }

// Mnemonics is every mnemonic with at least one form, sorted. This is what
// `arc isa` lists.
func Mnemonics() []string {
	out := make([]string, 0, len(byMnemonic))
	for m := range byMnemonic {
		out = append(out, m)
	}
	sort.Strings(out)
	return out
}

// Enabled is every form encodable under s, in table order.
func Enabled(s feature.Set) []*Form {
	var out []*Form
	for _, f := range forms {
		if f.Enabled(s) {
			out = append(out, f)
		}
	}
	return out
}

// Count is len(All()), for the generator's benefit and for a test that
// notices when a table file silently stops being registered.
func Count() int { return len(forms) }

func init() {
	// Every mnemonic is lowercase here. gas and nasm disagree about case and
	// both fold before they reach this package; a mixed-case entry would be
	// unreachable rather than wrong, which is worse.
	for _, f := range forms {
		if f.Op != strings.ToLower(f.Op) {
			panic(fmt.Sprintf("isa: mnemonic %q is not lowercase", f.Op))
		}
	}
}