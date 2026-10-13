// Package greeting builds dp-hello's output text. Its functions are pure (NFR-05).
package greeting

// HelpText is the help output from the user manual, byte for byte (FR-07, CONV-23).
const HelpText = `Usage: dp-hello [OPTION]...
Print a greeting to standard output.

Options:
  -n, --name=NAME     greet NAME instead of World
  -c, --count=COUNT   print the greeting COUNT times (1-1000, default 1)
      --help          display this help and exit
      --version       output version information and exit
`

// VersionLine is the --version output (FR-08, CONV-02).
const VersionLine = "dp-hello (dialectic-playground) 1.0.0\n"

// Line returns the greeting line for name: "Hello, " + name + "!" + LF (FR-01, FR-02).
// The name is used verbatim; Go strings carry the argument's UTF-8 bytes unchanged.
func Line(name string) string {
	return "Hello, " + name + "!\n"
}
