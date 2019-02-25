# go-bsdiff
Pure Go implementation of [bsdiff](http://www.daemonology.net/bsdiff/) 4.

bsdiff and bspatch are tools for building and applying patches to binary files. By using suffix sorting (specifically, Larsson and Sadakane's [qsufsort](http://www.larsson.dogma.net/ssrev-tr.pdf)) and taking advantage of how executable files change.

The package can be used as a library (pkg/bsdiff pkg/bspatch) or as a cli program (cmd/bsdiff cmd/bspatch).

## As a library

### Bsdiff Bytes
```Go
package main

import (
  "fmt"
  "bytes"

  "github.com/gabstv/go-bsdiff/pkg/bsdiff"
  "github.com/gabstv/go-bsdiff/pkg/bspatch"
)

func main(){
  oldfile := []byte{0xfa, 0xdd, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff}
  newfile := []byte{0xfa, 0xdd, 0x00, 0x00, 0x00, 0xee, 0xee, 0x00, 0x00, 0xff, 0xfe, 0xfe}
  patch, err := bsdiff.Bytes(oldfile, newfile)
  if err != nil {
    panic(err)
  }
  fmt.Println(patch)
  newfile2, err := bspatch.Bytes(oldfile, patch)
  if err != nil {
    panic(err)
  }
  if !bytes.Equal(newfile, newfile2) {
    panic()
  }
}
```
### Bsdiff Reader
```Go
package main

import (
  "fmt"
  "bytes"

  "github.com/gabstv/go-bsdiff/pkg/bsdiff"
  "github.com/gabstv/go-bsdiff/pkg/bspatch"
)

func main(){
  oldrdr := bytes.NewReader([]byte{0xfa, 0xdd, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff})
  newrdr := bytes.NewReader([]byte{0xfa, 0xdd, 0x00, 0x00, 0x00, 0xee, 0xee, 0x00, 0x00, 0xff, 0xfe, 0xfe})
  patch := new(bytes.Buffer)
  if err := bsdiff.Reader(oldrdr, newrdr, patch); err != nil {
    panic(err)
  }
  newpatchedf := new(bytes.Buffer)
  oldrdr.Seek(0, 0)
  if err := bspatch.Reader(oldrdr, newpatchedf, patch); err != nil {
    panic(err)
  }
  fmt.Println(newpatchedf.Bytes())
}
```

## As a program (CLI)
```sh
go get -u github.com/gabstv/go-bsdiff/cmd/...

bsdiff oldfile newfile patch
bspatch oldfile newfile2 patch
```