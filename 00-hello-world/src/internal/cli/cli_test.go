package cli

import (
	"reflect"
	"testing"
)

func value(text, option string) *Value {
	return &Value{Text: text, Option: option}
}

func TestParseAccepted(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want Invocation
	}{
		{"no arguments", nil, Invocation{}},
		{"short separate", []string{"-n", "Grace"}, Invocation{Name: value("Grace", "-n")}},
		{"short attached", []string{"-nAda"}, Invocation{Name: value("Ada", "-n")}},
		{"long separate", []string{"--name", "Ada"}, Invocation{Name: value("Ada", "--name")}},
		{"long inline", []string{"--name=Linus"}, Invocation{Name: value("Linus", "--name")}},
		{"long inline keeps later equals signs", []string{"--name=a=b"}, Invocation{Name: value("a=b", "--name")}},
		{"long inline empty", []string{"--name="}, Invocation{Name: value("", "--name")}},
		{"bundled values", []string{"-c2", "-nKen"}, Invocation{Name: value("Ken", "-n"), Count: value("2", "-c")}},
		{"bundle consumes rest of token", []string{"-nc3"}, Invocation{Name: value("c3", "-n")}},
		{"value may start with dash (attached)", []string{"-n-x"}, Invocation{Name: value("-x", "-n")}},
		{"value may start with dash (inline)", []string{"--name=-x"}, Invocation{Name: value("-x", "--name")}},
		{"value may start with dash (next token)", []string{"--name", "-x"}, Invocation{Name: value("-x", "--name")}},
		{"next-token value may be --", []string{"-n", "--"}, Invocation{Name: value("--", "-n")}},
		{"non-ASCII value", []string{"-n", "Zoë"}, Invocation{Name: value("Zoë", "-n")}},
		{"last occurrence wins", []string{"-n", "A", "--name", "B"}, Invocation{Name: value("B", "--name")}},
		{"count last wins with spelling", []string{"--count=0", "-c", "5"}, Invocation{Count: value("5", "-c")}},
		{"operand", []string{"extra"}, Invocation{Operands: []string{"extra"}}},
		{"single dash is an operand", []string{"-"}, Invocation{Operands: []string{"-"}}},
		{"operands interleaved", []string{"a", "-n", "B", "c"}, Invocation{Name: value("B", "-n"), Operands: []string{"a", "c"}}},
		{"terminator makes operands", []string{"--", "-n", "--help"}, Invocation{Operands: []string{"-n", "--help"}}},
		{"terminator alone", []string{"--"}, Invocation{}},
		{"help", []string{"--help"}, Invocation{Action: Help}},
		{"version", []string{"--version"}, Invocation{Action: Version}},
		{"help stops the scan", []string{"--help", "extra", "-x"}, Invocation{Action: Help}},
		{"version stops the scan", []string{"-n", "A", "--version", "--bogus"}, Invocation{Action: Version, Name: value("A", "-n")}},
		{"first of help and version wins", []string{"--version", "--help"}, Invocation{Action: Version}},
		{"help after operand", []string{"extra", "--help"}, Invocation{Action: Help, Operands: []string{"extra"}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse(tc.args)
			if err != nil {
				t.Fatalf("Parse(%q) error %q", tc.args, err.Error())
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Parse(%q) = %+v, want %+v", tc.args, got, tc.want)
			}
		})
	}
}

func TestParseRejected(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{"unknown short option", []string{"-x"}, "invalid option -- 'x'"},
		{"unknown short option is a whole letter", []string{"-ë"}, "invalid option -- 'ë'"},
		{"unknown long option", []string{"--nme", "Ada"}, "unrecognized option '--nme'"},
		{"no abbreviations", []string{"--nam=Ada"}, "unrecognized option '--nam=Ada'"},
		{"bare double dash with value", []string{"--=x"}, "unrecognized option '--=x'"},
		{"short missing argument", []string{"-n"}, "option requires an argument -- 'n'"},
		{"short count missing argument", []string{"-c"}, "option requires an argument -- 'c'"},
		{"long missing argument", []string{"--name"}, "option '--name' requires an argument"},
		{"long count missing argument", []string{"--count"}, "option '--count' requires an argument"},
		{"help with argument", []string{"--help=1"}, "option '--help' doesn't allow an argument"},
		{"version with argument", []string{"--version=1"}, "option '--version' doesn't allow an argument"},
		{"error before help wins", []string{"-x", "--help"}, "invalid option -- 'x'"},
		{"error stops at first problem", []string{"--bogus", "-x"}, "unrecognized option '--bogus'"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse(tc.args)
			if err == nil {
				t.Fatalf("Parse(%q) = %+v, want error %q", tc.args, got, tc.want)
			}
			if err.Error() != tc.want {
				t.Errorf("Parse(%q) error %q, want %q", tc.args, err.Error(), tc.want)
			}
		})
	}
}
