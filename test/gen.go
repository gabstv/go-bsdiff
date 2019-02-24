package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gabstv/go-bsdiff/pkg/bsdiff"
)

func main() {
	oldbs := []byte{0xFF, 0xFA, 0xB7, 0xDD}
	newbs := []byte{0xFF, 0xFA, 0x90, 0xB7, 0xDD, 0xFE}
	ioutil.WriteFile("old.bin", oldbs, 0644)
	ioutil.WriteFile("new.bin", newbs, 0644)
	f, _ := os.Create("godiff.bin")
	oldb, _ := os.Open("old.bin")
	newb, _ := os.Open("new.bin")
	defer f.Close()
	defer oldb.Close()
	defer newb.Close()
	bsdiff.Stream(oldb, newb, f)
	f.Close()
	hhh, _ := ioutil.ReadFile("godiff.bin")
	fmt.Println(hhh)
	if _, err := os.Stat("diff.bin"); err == nil {
		hhh, _ = ioutil.ReadFile("diff.bin")
		fmt.Println(hhh)
	}
}
