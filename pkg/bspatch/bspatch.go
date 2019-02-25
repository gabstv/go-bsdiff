package bspatch

import (
	"bytes"
	"fmt"
	"io"

	"github.com/dsnet/compress/bzip2"
)

func patchb(oldfile, patch []byte) ([]byte, error) {
	oldsize := len(oldfile)
	var newsize int
	header := make([]byte, 32)
	buf := make([]byte, 8)
	var lenread int
	var i int
	ctrl := make([]int, 3)

	f := bytes.NewReader(patch)

	/*
		File format:
			0	8	"BSDIFF40"
			8	8	X
			16	8	Y
			24	8	sizeof(newfile)
			32	X	bzip2(control block)
			32+X	Y	bzip2(diff block)
			32+X+Y	???	bzip2(extra block)
		with control block a set of triples (x,y,z) meaning "add x bytes
		from oldfile to x bytes from the diff block; copy y bytes from the
		extra block; seek forwards in oldfile by z bytes".
	*/

	/* Read header */
	if n, err := f.Read(header); err != nil || n < 32 {
		if err != nil {
			return nil, fmt.Errorf("corrupt patch %v", err.Error())
		}
		return nil, fmt.Errorf("corrupt patch")
	}
	/* Check for appropriate magic */
	if bytes.Compare(header[:8], []byte("BSDIFF40")) != 0 {
		return nil, fmt.Errorf("corrupt patch (header)")
	}

	/* Read lengths from header */
	// L 109
	bzctrllen := offtin(buf[8:])
	bzdatalen := offtin(header[16:])
	newsize = offtin(header[24:])

	if bzctrllen < 0 || bzdatalen < 0 || newsize < 0 {
		return nil, fmt.Errorf("corrupt patch")
	}

	/* Close patch file and re-open it via libbzip2 at the right places */
	f = nil
	cpf := bytes.NewReader(patch)
	if _, err := cpf.Seek(32, io.SeekStart); err != nil {
		return nil, err
	}
	cpfbz2, err := bzip2.NewReader(cpf, nil)
	if err != nil {
		return nil, err
	}
	dpf := bytes.NewReader(patch)
	if _, err := dpf.Seek(int64(32+bzctrllen), io.SeekStart); err != nil {
		return nil, err
	}
	dpfbz2, err := bzip2.NewReader(dpf, nil)
	if err != nil {
		return nil, err
	}
	epf := bytes.NewReader(patch)
	if _, err := epf.Seek(int64(32+bzctrllen+bzdatalen), io.SeekStart); err != nil {
		return nil, err
	}
	epfbz2, err := bzip2.NewReader(epf, nil)
	if err != nil {
		return nil, err
	}

	pnew := make([]byte, newsize) // newsize+1

	oldpos := 0
	newpos := 0

	// L 154:
	for newpos < newsize {
		/* Read control data */
		for i = 0; i <= 2; i++ {
			lenread, err = cpfbz2.Read(buf)
			if lenread != 8 || (err != nil && err != io.EOF) {
				e0 := ""
				if err != nil {
					e0 = err.Error()
				}
				return nil, fmt.Errorf("corrupt patch or bzstream ended: %s", e0)
			}
			ctrl[i] = offtin(buf)
		}
		/* Sanity-check */
		if newpos+ctrl[0] > newsize {
			return nil, fmt.Errorf("corrupt patch (sanity check)")
		}

		/* Read diff string */
		lenread, err = dpfbz2.Read(pnew[newpos : newpos+ctrl[0]])
		if lenread < ctrl[0] || (err != nil && err != io.EOF) {
			e0 := ""
			if err != nil {
				e0 = err.Error()
			}
			return nil, fmt.Errorf("corrupt patch or bzstream ended (2): %s", e0)
		}
		/* Add pold data to diff string */
		for i = 0; i < ctrl[0]; i++ {
			if oldpos+i >= 0 && oldpos+i < oldsize {
				pnew[newpos+i] += oldfile[oldpos+i]
			}
		}

		// Adjust pointers
		newpos += ctrl[0]
		oldpos += ctrl[0]

		// Sanity-check
		if newpos+ctrl[1] > newsize {
			return nil, fmt.Errorf("corrupt patch")
		}

		// Read extra string
		lenread, err = epfbz2.Read(pnew[newpos : newpos+ctrl[1]])
		if lenread < ctrl[1] || (err != nil && err != io.EOF) {
			e0 := ""
			if err != nil {
				e0 = err.Error()
			}
			return nil, fmt.Errorf("corrupt patch or bzstream ended (3): %s", e0)
		}
		// Adjust pointers
		newpos += ctrl[1]
		oldpos += ctrl[2]
	}

	// Clean up the bzip2 reads
	if err = cpfbz2.Close(); err != nil {
		return nil, err
	}
	if err = dpfbz2.Close(); err != nil {
		return nil, err
	}
	if err = epfbz2.Close(); err != nil {
		return nil, err
	}
	cpfbz2 = nil
	dpfbz2 = nil
	epfbz2 = nil
	cpf = nil
	dpf = nil
	epf = nil

	return pnew, nil
}

func offtin(buf []byte) int {

	y := int(buf[7] & 0x7f)
	y = y * 256
	y += int(buf[6])
	y = y * 256
	y += int(buf[5])
	y = y * 256
	y += int(buf[4])
	y = y * 256
	y += int(buf[3])
	y = y * 256
	y += int(buf[2])
	y = y * 256
	y += int(buf[1])
	y = y * 256
	y += int(buf[0])

	if (buf[7] & 0x80) != 0 {
		y = -y
	}
	return y
}
