package util

import (
	"fmt"
	"io"
)

const (
	buffersize = 1024 * 16
)

// PutWriter writes all content from b to target
func PutWriter(target io.Writer, b []byte) error {
	lb := len(b)
	if lb < buffersize {
		n, err := target.Write(b)
		if err != nil {
			return err
		}
		if lb != n {
			return fmt.Errorf("%v of %v bytes written", n, lb)
		}
		return nil
	}
	offs := 0

	for offs < lb {
		n := Min(buffersize, lb-offs)
		n2, err := target.Write(b[offs:n])
		if err != nil {
			return err
		}
		if n2 != n {
			return fmt.Errorf("%v of %v bytes written", offs+n2, lb)
		}
		offs += n
	}
	return nil
}

// Min int
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
