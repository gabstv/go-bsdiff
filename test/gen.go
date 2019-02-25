package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gabstv/go-bsdiff/pkg/bsdiff"
	"github.com/gabstv/go-bsdiff/pkg/bspatch"
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
	var ddd []byte
	if _, err := os.Stat("diff.bin"); err == nil {
		ddd, _ = ioutil.ReadFile("diff.bin")
		fmt.Println(ddd)
	}
	newbs2, err := bspatch.Bytes(oldbs, hhh)
	if err != nil {
		fmt.Println("go patch error:", err.Error())
		return
	}
	newbs3, err := bspatch.Bytes(oldbs, ddd)
	if err != nil {
		fmt.Println("'brew' patch error:", err.Error())
		return
	}
	if bytes.Equal(newbs, newbs2) {
		fmt.Println("go diff success!")
	} else {
		fmt.Println("go diff FAILED!")
	}
	if bytes.Equal(newbs, newbs3) {
		fmt.Println("'brew' diff success!")
	} else {
		fmt.Println("'brew' diff FAILED!")
	}
}
