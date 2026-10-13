package greeting

import (
	"strings"
	"testing"
)

func TestLine(t *testing.T) {
	cases := map[string]string{
		"World":        "Hello, World!\n",
		"Ada Lovelace": "Hello, Ada Lovelace!\n",
		"Zoë":          "Hello, Zo\xc3\xab!\n",
		"-x":           "Hello, -x!\n",
		"Hi!$`'\"":     "Hello, Hi!$`'\"!\n",
	}
	for name, want := range cases {
		if got := Line(name); got != want {
			t.Errorf("Line(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestHelpTextLayout(t *testing.T) {
	if !strings.HasPrefix(HelpText, "Usage: dp-hello [OPTION]...\n") {
		t.Errorf("help text does not start with the usage line: %q", HelpText)
	}
	if strings.Contains(HelpText, "\r") || !strings.HasSuffix(HelpText, "exit\n") {
		t.Errorf("help text must use LF only and end with LF: %q", HelpText)
	}
	if got := strings.Count(HelpText, "\n"); got != 8 {
		t.Errorf("help text has %d lines, want 8", got)
	}
}

func TestVersionLine(t *testing.T) {
	if VersionLine != "dp-hello (dialectic-playground) 1.0.0\n" {
		t.Errorf("VersionLine = %q", VersionLine)
	}
}
