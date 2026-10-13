// Command dp-hello prints a greeting to standard output (curriculum project 00).
//
// os.Args holds the arguments as UTF-8 on every platform: on Windows the Go runtime
// decodes the UTF-16 command line itself (CONV-05). os.Stdout and os.Stderr are written
// as raw bytes, without newline translation (CONV-03).
package main

import (
	"os"

	"github.com/dialectic-playground/go/00-hello-world/src/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:], os.Stdout, os.Stderr))
}
