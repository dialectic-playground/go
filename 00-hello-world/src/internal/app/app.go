// Package app wires dp-hello's units together: parsing (cli), validation (validate),
// greeting construction (greeting), diagnostics, and output. Run is the only place that
// maps outcomes to exit statuses.
package app

import (
	"io"

	"github.com/dialectic-playground/go/00-hello-world/src/internal/cli"
	"github.com/dialectic-playground/go/00-hello-world/src/internal/diagnostics"
	"github.com/dialectic-playground/go/00-hello-world/src/internal/greeting"
	"github.com/dialectic-playground/go/00-hello-world/src/internal/output"
	"github.com/dialectic-playground/go/00-hello-world/src/internal/validate"
)

// Exit statuses (SRS §3.7, CONV-18).
const (
	ExitSuccess = 0
	// ExitOutputFailure is returned when standard output cannot be written (for example a
	// broken pipe, whose status CONV-25 leaves unspecified). No diagnostic is printed.
	ExitOutputFailure = 1
	ExitUsage         = 2
)

// Run executes dp-hello with args (the command line without the program name) and returns
// the exit status. Parsing and validation finish before anything is written, so a usage
// error never produces standard output (FR-11).
func Run(args []string, stdout, stderr io.Writer) int {
	invocation, usageErr := cli.Parse(args)
	if usageErr != nil {
		return reportUsage(stderr, usageErr)
	}
	var err error
	switch invocation.Action {
	case cli.Help:
		err = output.Write(stdout, greeting.HelpText)
	case cli.Version:
		err = output.Write(stdout, greeting.VersionLine)
	default:
		config, usageErr := validate.Validate(invocation)
		if usageErr != nil {
			return reportUsage(stderr, usageErr)
		}
		err = output.WriteRepeated(stdout, greeting.Line(config.Name), config.Count)
	}
	if err != nil {
		return ExitOutputFailure
	}
	return ExitSuccess
}

// reportUsage writes a usage error to stderr (CONV-13, CONV-14) and returns ExitUsage.
func reportUsage(stderr io.Writer, usageErr *diagnostics.UsageError) int {
	// A failure to write the diagnostic itself cannot be reported anywhere.
	_ = output.Write(stderr, usageErr.Report())
	return ExitUsage
}
