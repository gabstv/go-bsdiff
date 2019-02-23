package bsdiff

import (
	"bytes"
	"testing"
)

func TestDiff(t *testing.T) {
	oldrdr := bytes.NewReader([]byte{0xFF, 0xFA, 0xB7, 0xDD})
	newrdr := bytes.NewReader([]byte{0xFF, 0xFA, 0x90, 0xB7, 0xDD, 0xFE})
	wr := &BufWriter{}
	if err := Diff(oldrdr, 4, newrdr, 6, wr); err != nil {
		t.Fatal(err)
	}
}
