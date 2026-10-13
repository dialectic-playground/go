package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {
	var out bytes.Buffer
	if err := Write(&out, "a\nb\n"); err != nil || out.String() != "a\nb\n" {
		t.Fatalf("Write: %q, %v", out.String(), err)
	}
}

func TestWriteRepeated(t *testing.T) {
	var out bytes.Buffer
	if err := WriteRepeated(&out, "Hello, Zoë!\n", 1000); err != nil {
		t.Fatal(err)
	}
	if want := strings.Repeat("Hello, Zoë!\n", 1000); out.String() != want {
		t.Fatalf("WriteRepeated wrote %d bytes, want %d", out.Len(), len(want))
	}
}

type failingWriter struct{}

var errClosed = errors.New("reader went away")

func (failingWriter) Write([]byte) (int, error) { return 0, errClosed }

func TestWriteReportsFailure(t *testing.T) {
	if err := Write(failingWriter{}, "x\n"); !errors.Is(err, errClosed) {
		t.Errorf("Write error = %v, want %v", err, errClosed)
	}
	if err := WriteRepeated(failingWriter{}, "x\n", 3); !errors.Is(err, errClosed) {
		t.Errorf("WriteRepeated error = %v, want %v", err, errClosed)
	}
}
