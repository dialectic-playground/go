package validate

import (
	"strings"
	"testing"

	"github.com/dialectic-playground/go/00-hello-world/src/internal/cli"
)

func TestParseCountAccepted(t *testing.T) {
	cases := map[string]int{"1": 1, "2": 2, "999": 999, "1000": 1000, "007": 7, "0001000": 1000}
	for text, want := range cases {
		got, ok := ParseCount(text)
		if !ok || got != want {
			t.Errorf("ParseCount(%q) = %d, %v; want %d, true", text, got, ok, want)
		}
	}
}

func TestParseCountRejected(t *testing.T) {
	rejected := []string{
		"", "0", "000", "1001", "99999", strings.Repeat("9", 40),
		"abc", "-1", "+1", " 1", "1 ", "1.0", "1e3", "٣", "１",
	}
	for _, text := range rejected {
		if got, ok := ParseCount(text); ok {
			t.Errorf("ParseCount(%q) = %d, true; want rejection", text, got)
		}
	}
}

func TestValidateDefaults(t *testing.T) {
	got, err := Validate(cli.Invocation{})
	if err != nil || got != (Config{Name: "World", Count: 1}) {
		t.Fatalf("Validate(default) = %+v, %v", got, err)
	}
}

func TestValidateValues(t *testing.T) {
	inv := cli.Invocation{
		Name:  &cli.Value{Text: "Ada Lovelace", Option: "--name"},
		Count: &cli.Value{Text: "007", Option: "-c"},
	}
	got, err := Validate(inv)
	if err != nil || got != (Config{Name: "Ada Lovelace", Count: 7}) {
		t.Fatalf("Validate = %+v, %v", got, err)
	}
}

func TestValidateRejected(t *testing.T) {
	cases := []struct {
		name string
		inv  cli.Invocation
		want string
	}{
		{"empty name (long)", cli.Invocation{Name: &cli.Value{Text: "", Option: "--name"}},
			"invalid argument '' for '--name'"},
		{"empty name (short)", cli.Invocation{Name: &cli.Value{Text: "", Option: "-n"}},
			"invalid argument '' for '-n'"},
		{"count zero", cli.Invocation{Count: &cli.Value{Text: "0", Option: "--count"}},
			"invalid argument '0' for '--count'"},
		{"count too large", cli.Invocation{Count: &cli.Value{Text: "1001", Option: "-c"}},
			"invalid argument '1001' for '-c'"},
		{"count not digits", cli.Invocation{Count: &cli.Value{Text: "abc", Option: "-c"}},
			"invalid argument 'abc' for '-c'"},
		{"count empty", cli.Invocation{Count: &cli.Value{Text: "", Option: "--count"}},
			"invalid argument '' for '--count'"},
		{"operand first", cli.Invocation{
			Name:     &cli.Value{Text: "", Option: "-n"},
			Operands: []string{"extra", "more"},
		}, "extra operand 'extra'"},
		{"name before count", cli.Invocation{
			Name:  &cli.Value{Text: "", Option: "-n"},
			Count: &cli.Value{Text: "0", Option: "-c"},
		}, "invalid argument '' for '-n'"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Validate(tc.inv)
			if err == nil {
				t.Fatalf("Validate = %+v, want error %q", got, tc.want)
			}
			if err.Error() != tc.want {
				t.Errorf("error %q, want %q", err.Error(), tc.want)
			}
		})
	}
}
