package diagnostics

import "testing"

func TestReportFormatsBothLines(t *testing.T) {
	got := ExtraOperand("foo").Report()
	want := "dp-hello: extra operand 'foo'\nTry 'dp-hello --help' for more information.\n"
	if got != want {
		t.Fatalf("Report() = %q, want %q", got, want)
	}
}

func TestMessages(t *testing.T) {
	cases := []struct {
		name string
		err  *UsageError
		want string
	}{
		{"invalid short option", InvalidShortOption("x"), "invalid option -- 'x'"},
		{"unrecognized option", UnrecognizedOption("--nme=Ada"), "unrecognized option '--nme=Ada'"},
		{"short requires argument", ShortOptionRequiresArgument("n"), "option requires an argument -- 'n'"},
		{"long requires argument", LongOptionRequiresArgument("count"), "option '--count' requires an argument"},
		{"long disallows argument", LongOptionDisallowsArgument("help"), "option '--help' doesn't allow an argument"},
		{"empty invalid argument", InvalidArgument("", "--name"), "invalid argument '' for '--name'"},
		{"invalid argument", InvalidArgument("abc", "-c"), "invalid argument 'abc' for '-c'"},
		{"extra operand", ExtraOperand("Zoë"), "extra operand 'Zoë'"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.err.Error(); got != tc.want {
				t.Errorf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}
