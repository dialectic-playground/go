package app

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/dialectic-playground/go/00-hello-world/src/internal/greeting"
)

// run invokes Run with in-memory streams and returns (stdout, stderr, status).
func run(args ...string) (string, string, int) {
	var stdout, stderr bytes.Buffer
	status := Run(args, &stdout, &stderr)
	return stdout.String(), stderr.String(), status
}

func usage(message string) string {
	return "dp-hello: " + message + "\nTry 'dp-hello --help' for more information.\n"
}

func TestRunSuccess(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{"default greeting", nil, "Hello, World!\n"},
		{"named", []string{"--name", "Ada Lovelace"}, "Hello, Ada Lovelace!\n"},
		{"non-ASCII name", []string{"-c2", "-nZoë"}, "Hello, Zoë!\nHello, Zoë!\n"},
		{"leading zeros", []string{"--count=003"}, strings.Repeat("Hello, World!\n", 3)},
		{"maximum count", []string{"-c", "1000"}, strings.Repeat("Hello, World!\n", 1000)},
		{"last wins; earlier invalid values ignored", []string{"-n", "", "-n", "B", "-c", "0", "-c", "1"}, "Hello, B!\n"},
		{"help", []string{"--help"}, greeting.HelpText},
		{"help ends the scan", []string{"--help", "extra", "-x"}, greeting.HelpText},
		{"help beats earlier invalid values", []string{"-c", "0", "--help"}, greeting.HelpText},
		{"help beats earlier operands", []string{"extra", "--help"}, greeting.HelpText},
		{"version", []string{"--version"}, greeting.VersionLine},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, status := run(tc.args...)
			if stdout != tc.want || stderr != "" || status != ExitSuccess {
				t.Errorf("Run(%q) = (%q, %q, %d), want (%q, \"\", 0)", tc.args, stdout, stderr, status, tc.want)
			}
		})
	}
}

func TestRunUsageErrors(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		message string
	}{
		{"unknown short", []string{"-x"}, "invalid option -- 'x'"},
		{"unknown long", []string{"--nme", "Ada"}, "unrecognized option '--nme'"},
		{"missing value", []string{"--name"}, "option '--name' requires an argument"},
		{"version with value", []string{"--version=1"}, "option '--version' doesn't allow an argument"},
		{"error before help", []string{"-x", "--help"}, "invalid option -- 'x'"},
		{"empty name", []string{"--name", ""}, "invalid argument '' for '--name'"},
		{"count zero", []string{"--count", "0"}, "invalid argument '0' for '--count'"},
		{"count non-digit", []string{"-c", "abc"}, "invalid argument 'abc' for '-c'"},
		{"count above maximum", []string{"-c", "1001"}, "invalid argument '1001' for '-c'"},
		{"operand", []string{"extra"}, "extra operand 'extra'"},
		{"operand after terminator", []string{"--", "-n"}, "extra operand '-n'"},
		{"non-ASCII operand", []string{"Zoë"}, "extra operand 'Zoë'"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, status := run(tc.args...)
			if stdout != "" || stderr != usage(tc.message) || status != ExitUsage {
				t.Errorf("Run(%q) = (%q, %q, %d), want (\"\", %q, 2)", tc.args, stdout, stderr, status, usage(tc.message))
			}
		})
	}
}

type brokenPipe struct{}

func (brokenPipe) Write([]byte) (int, error) { return 0, errors.New("broken pipe") }

func TestRunOutputFailureIsSilent(t *testing.T) {
	var stderr bytes.Buffer
	if status := Run(nil, brokenPipe{}, &stderr); status != ExitOutputFailure {
		t.Errorf("status = %d, want %d", status, ExitOutputFailure)
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
}
