// Package bsdiff is a pure Go implementation of Bsdiff 4.
//
// Example:
//  package main
//
//  import (
//    "fmt"
//    "bytes"
//
//    "github.com/gabstv/go-bsdiff/pkg/bsdiff"
//    "github.com/gabstv/go-bsdiff/pkg/bspatch"
//  )
//
//  func main(){
//    // example files
//    oldfile := []byte{0xfa, 0xdd, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff}
//    newfile := []byte{0xfa, 0xdd, 0x00, 0x00, 0x00, 0xee, 0xee, 0x00, 0x00, 0xff, 0xfe, 0xfe}
//
//    // generate a BSDIFF4 patch
//    patch, err := bsdiff.Bytes(oldfile, newfile)
//    if err != nil {
//      panic(err)
//    }
//    fmt.Println(patch)
//
//    // Apply a BSDIFF4 patch
//    newfile2, err := bspatch.Bytes(oldfile, patch)
//    if err != nil {
//      panic(err)
//    }
//    if !bytes.Equal(newfile, newfile2) {
//      panic()
//    }
//  }
package bsdiff

func v() int {
	return 1
}
