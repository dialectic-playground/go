// Package output writes dp-hello's prepared text to a stream.
//
// Go's os.File never translates line endings, so the LF-only bytes reach standard output
// and standard error unchanged on every platform, including Windows (CONV-03). A write
// failure is returned to the caller, which ends the tool without a diagnostic (CONV-25).
package output

import (
	"bufio"
	"io"
)

// Write writes all of text to w.
func Write(w io.Writer, text string) error {
	_, err := io.WriteString(w, text)
	return err
}

// WriteRepeated writes line to w count times through a buffer, so memory use stays
// bounded by the buffer size plus one line, whatever the count (NFR-02).
func WriteRepeated(w io.Writer, line string, count int) error {
	buffered := bufio.NewWriter(w)
	for range count {
		if _, err := buffered.WriteString(line); err != nil {
			return err
		}
	}
	return buffered.Flush()
}
