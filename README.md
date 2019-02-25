# go-bsdiff
Pure Go implementation of [bsdiff](http://www.daemonology.net/bsdiff/) 4.0

bsdiff and bspatch are tools for building and applying patches to binary files. By using suffix sorting (specifically, Larsson and Sadakane's [qsufsort](http://www.larsson.dogma.net/ssrev-tr.pdf)) and taking advantage of how executable files change.

```Go
package main

import (
  "fmt"

  "github.com/gabstv/go-bsdiff/pkg/bsdiff"
)

func main(){
  oldfile := []byte{0xfa, 0xdd, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff}
  newfile := []byte{0xfa, 0xdd, 0x00, 0x00, 0x00, 0xee, 0xee, 0x00, 0x00, 0xff, 0xfe, 0xfe}
  patch, err := bsdiff.Bytes(oldfile, newfile)
  if err != nil {
    panic(err)
  }
  fmt.Println(patch)
}
```
