// Package cli turns dp-hello's command-line arguments into an Invocation, following the
// option syntax of CONV-07 … CONV-12. It checks syntax only; option values are checked by
// package validate.
package cli

import (
	"strings"
	"unicode/utf8"

	"github.com/dialectic-playground/go/00-hello-world/src/internal/diagnostics"
)

// Action is what the invocation asks dp-hello to do.
type Action int

const (
	// Greet prints the greeting (the default).
	Greet Action = iota
	// Help prints the help text (--help).
	Help
	// Version prints the version line (--version).
	Version
)

// Value is an option value together with the option as the user spelled it ("-n" or
// "--name"), which the invalid-argument diagnostic needs (CONV-13).
type Value struct {
	Text   string
	Option string
}

// Invocation is the result of a successful parse. Name and Count are nil when the option
// was not given; otherwise they hold the last occurrence (CONV-11). Operands lists every
// operand in order (CONV-09, CONV-10).
type Invocation struct {
	Action   Action
	Name     *Value
	Count    *Value
	Operands []string
}

// optionKind identifies one of dp-hello's options.
type optionKind int

const (
	nameOption optionKind = iota
	countOption
	helpOption
	versionOption
)

// option describes one option of dp-hello: its spellings and whether it takes a value.
type option struct {
	kind       optionKind
	short      string // "" when the option has no short form
	long       string
	takesValue bool
}

var options = []option{
	{kind: nameOption, short: "n", long: "name", takesValue: true},
	{kind: countOption, short: "c", long: "count", takesValue: true},
	{kind: helpOption, long: "help"},
	{kind: versionOption, long: "version"},
}

func findShort(letter string) (option, bool) {
	for _, candidate := range options {
		if candidate.short == letter {
			return candidate, true
		}
	}
	return option{}, false
}

func findLong(name string) (option, bool) {
	for _, candidate := range options {
		if candidate.long == name {
			return candidate, true
		}
	}
	return option{}, false
}

// store records a value for a valued option, replacing any earlier one (CONV-11).
func (inv *Invocation) store(opt option, text, spelled string) {
	value := &Value{Text: text, Option: spelled}
	switch opt.kind {
	case nameOption:
		inv.Name = value
	case countOption:
		inv.Count = value
	}
}

// Parse scans args (the command line without the program name) from left to right. It
// stops at the first --help or --version and returns that action (CONV-12). The first
// syntax error found before that point is returned as a usage error instead.
func Parse(args []string) (Invocation, *diagnostics.UsageError) {
	var inv Invocation
	for index := 0; index < len(args); index++ {
		token := args[index]
		switch {
		case token == "--":
			inv.Operands = append(inv.Operands, args[index+1:]...)
			return inv, nil
		case strings.HasPrefix(token, "--"):
			name, inline, hasInline := strings.Cut(token[2:], "=")
			opt, known := findLong(name)
			if !known {
				return Invocation{}, diagnostics.UnrecognizedOption(token)
			}
			if !opt.takesValue {
				if hasInline {
					return Invocation{}, diagnostics.LongOptionDisallowsArgument(name)
				}
				inv.Action = Help
				if opt.kind == versionOption {
					inv.Action = Version
				}
				return inv, nil
			}
			value := inline
			if !hasInline {
				if index+1 >= len(args) {
					return Invocation{}, diagnostics.LongOptionRequiresArgument(name)
				}
				index++
				value = args[index]
			}
			inv.store(opt, value, "--"+name)
		case strings.HasPrefix(token, "-") && token != "-":
			consumedNext, err := inv.parseShortBundle(token, args[index+1:])
			if err != nil {
				return Invocation{}, err
			}
			if consumedNext {
				index++
			}
		default:
			inv.Operands = append(inv.Operands, token)
		}
	}
	return inv, nil
}

// parseShortBundle parses one token of short options such as "-c2" or "-n" (CONV-07).
// rest holds the tokens after it; consumedNext reports whether the option's value was
// taken from rest[0]. Every short option of dp-hello takes a value, so the first letter's
// value is the rest of the token or, when that is empty, the next token.
func (inv *Invocation) parseShortBundle(token string, rest []string) (consumedNext bool, err *diagnostics.UsageError) {
	_, size := utf8.DecodeRuneInString(token[1:])
	letter := token[1 : 1+size]
	opt, known := findShort(letter)
	if !known {
		return false, diagnostics.InvalidShortOption(letter)
	}
	if attached := token[1+size:]; attached != "" {
		inv.store(opt, attached, "-"+letter)
		return false, nil
	}
	if len(rest) == 0 {
		return false, diagnostics.ShortOptionRequiresArgument(letter)
	}
	inv.store(opt, rest[0], "-"+letter)
	return true, nil
}
