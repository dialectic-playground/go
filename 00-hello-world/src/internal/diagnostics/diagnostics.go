// Package diagnostics defines the usage error of dp-hello and its exact CONV-13/CONV-14
// wording. It is the only place where diagnostic text is spelled out.
package diagnostics

// Tool is the canonical tool name used in every message (CONV-01).
const Tool = "dp-hello"

// UsageError is a CONV-13 usage error. Message is the text after "dp-hello: ".
type UsageError struct {
	Message string
}

// Error returns the message without the tool prefix.
func (e *UsageError) Error() string {
	return e.Message
}

// Report returns the two lines written to standard error: the diagnostic and the
// "Try ..." hint (CONV-13, CONV-14), each terminated by LF.
func (e *UsageError) Report() string {
	return Tool + ": " + e.Message + "\nTry '" + Tool + " --help' for more information.\n"
}

// InvalidShortOption reports an unknown short option letter.
func InvalidShortOption(letter string) *UsageError {
	return &UsageError{Message: "invalid option -- '" + letter + "'"}
}

// UnrecognizedOption reports an unknown long option token, exactly as given.
func UnrecognizedOption(token string) *UsageError {
	return &UsageError{Message: "unrecognized option '" + token + "'"}
}

// ShortOptionRequiresArgument reports a short option that is missing its value.
func ShortOptionRequiresArgument(letter string) *UsageError {
	return &UsageError{Message: "option requires an argument -- '" + letter + "'"}
}

// LongOptionRequiresArgument reports a long option (name without "--") missing its value.
func LongOptionRequiresArgument(name string) *UsageError {
	return &UsageError{Message: "option '--" + name + "' requires an argument"}
}

// LongOptionDisallowsArgument reports "--name=value" for an option that takes no value.
func LongOptionDisallowsArgument(name string) *UsageError {
	return &UsageError{Message: "option '--" + name + "' doesn't allow an argument"}
}

// InvalidArgument reports an invalid value for the option as the user spelled it
// (for example "-n" or "--count"). The value is inserted verbatim (CONV-15).
func InvalidArgument(value, spelledOption string) *UsageError {
	return &UsageError{Message: "invalid argument '" + value + "' for '" + spelledOption + "'"}
}

// ExtraOperand reports the first surplus operand.
func ExtraOperand(operand string) *UsageError {
	return &UsageError{Message: "extra operand '" + operand + "'"}
}
