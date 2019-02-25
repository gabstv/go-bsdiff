package bsdiff

import (
	"bytes"
	"testing"

	"github.com/gabstv/go-bsdiff/pkg/bsdiff"
	"github.com/gabstv/go-bsdiff/pkg/bspatch"
)

func TestDiffPatch(t *testing.T) {
	oldbs := []byte{0xFF, 0xFA, 0xB7, 0xDD}
	newbs := []byte{0xFF, 0xFA, 0x90, 0xB7, 0xDD, 0xFE}
	patch, err := bsdiff.Bytes(oldbs, newbs)
	if err != nil {
		t.Fatal(err.Error())
	}
	newbs2, err := bspatch.Bytes(oldbs, patch)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !bytes.Equal(newbs, newbs2) {
		t.Fatal(newbs2, "!=", newbs)
	}
}
