// Package validate checks the values of a parsed greeting request and turns them into a
// Config (SRS FR-03 … FR-06, FR-09).
package validate

import (
	"github.com/dialectic-playground/go/00-hello-world/src/internal/cli"
	"github.com/dialectic-playground/go/00-hello-world/src/internal/diagnostics"
)

const (
	// DefaultName is greeted when --name is not given (FR-01).
	DefaultName = "World"
	// DefaultCount is the number of greeting lines when --count is not given.
	DefaultCount = 1
	// MinCount and MaxCount bound COUNT (FR-04).
	MinCount = 1
	MaxCount = 1000
)

// Config is a validated greeting request: Name is non-empty and Count is in
// [MinCount, MaxCount].
type Config struct {
	Name  string
	Count int
}

// Validate checks a greeting invocation. It reports the first operand (FR-09), then an
// empty name (FR-03), then an invalid count (FR-05). Only the last occurrence of each
// option reaches this point, so earlier values are never validated (FR-06).
func Validate(inv cli.Invocation) (Config, *diagnostics.UsageError) {
	if len(inv.Operands) > 0 {
		return Config{}, diagnostics.ExtraOperand(inv.Operands[0])
	}
	config := Config{Name: DefaultName, Count: DefaultCount}
	if inv.Name != nil {
		if inv.Name.Text == "" {
			return Config{}, diagnostics.InvalidArgument(inv.Name.Text, inv.Name.Option)
		}
		config.Name = inv.Name.Text
	}
	if inv.Count != nil {
		count, ok := ParseCount(inv.Count.Text)
		if !ok {
			return Config{}, diagnostics.InvalidArgument(inv.Count.Text, inv.Count.Option)
		}
		config.Count = count
	}
	return config, nil
}

// ParseCount converts COUNT to an integer. COUNT must be one or more ASCII digits
// (leading zeros allowed) with a value in [MinCount, MaxCount]; otherwise ok is false.
// Arbitrarily long digit strings are handled without overflow.
func ParseCount(text string) (count int, ok bool) {
	if text == "" {
		return 0, false
	}
	for i := 0; i < len(text); i++ {
		digit := text[i]
		if digit < '0' || digit > '9' {
			return 0, false
		}
		if count <= MaxCount {
			count = count*10 + int(digit-'0')
		}
	}
	if count < MinCount || count > MaxCount {
		return 0, false
	}
	return count, true
}
