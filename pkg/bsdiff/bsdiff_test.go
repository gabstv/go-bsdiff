package bsdiff

import (
	"testing"
)

func TestDiff(t *testing.T) {
	oldbs := []byte{0xFF, 0xFA, 0xB7, 0xDD}
	newbs := []byte{0xFF, 0xFA, 0x90, 0xB7, 0xDD, 0xFE}
	var diffbs []byte
	var err error
	if diffbs, err = Bytes(oldbs, newbs); err != nil {
		t.Fatal(err)
	}
	t.Fatal(diffbs)
}
