package main

import (
	"os"

	"github.com/gabstv/go-bsdiff/pkg/bspatch"
)

func main() {
	if len(os.Args) != 4 {
		printusage(1)
	}
	err := bspatch.File(os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		println(err.Error())
		printusage(1)
	}
}

func printusage(exitcode int) {
	println("usage: " + os.Args[0] + " oldfile newfile patchfile")
	os.Exit(exitcode)
}
