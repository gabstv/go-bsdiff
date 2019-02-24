package bsdiff

import (
	"bytes"
	"fmt"
	"io"

	"github.com/dsnet/compress/bzip2"
)

const (
	buffersize = 1024 * 16
)

// https://github.com/cnSchwarzer/bsdiff-win/blob/master/bsdiff-win/bsdiff.c

func DiffFiles(oldbin io.ReadSeeker, newbin io.ReadSeeker, diffbin io.Writer) error {
	oldsize, err := oldbin.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	pold := make([]byte, int(oldsize))               // pold = (u_char *)malloc(oldsize + 1);
	oldbin.Seek(0, io.SeekStart)                     // fseek(fs, 0, SEEK_SET);
	if err := copyReader(pold, oldbin); err != nil { // if (fread(pold, 1, oldsize, fs) == -1)	err(1, "Read failed :%s", argv[1]);
		return err
	}
	//
	newsize, err := newbin.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	pnew := make([]byte, int(newsize))
	newbin.Seek(0, io.SeekStart)
	if err := copyReader(pnew, newbin); err != nil {
		return err
	}
	diffbytes, err := diffb(pold, pnew)
	if err != nil {
		return err
	}
	return putWriter(diffbin, diffbytes)
}

func diffb(oldbin, newbin []byte) ([]byte, error) {
	iii := make([]int, len(oldbin)+1)
	vvv := make([]int, len(oldbin)+1)
	qsufsort(iii, vvv, oldbin)
	// [C] free(V)
	vvv = nil

	//var db
	var dblen, eblen int

	// create the patch file
	pf := new(bytes.Buffer)

	// Header is
	//	0	8	 "BSDIFF40"
	//	8	8	length of bzip2ed ctrl block
	//	16	8	length of bzip2ed diff block
	//	24	8	length of pnew file */
	//	/* File is
	//		0	32	Header
	//		32	??	Bzip2ed ctrl block
	//		??	??	Bzip2ed diff block
	//		??	??	Bzip2ed extra block

	header := make([]byte, 32)

	copy(header, []byte("BSDIFF40"))
	offtout(0, header[8:])
	offtout(0, header[16:])
	offtout(len(newbin), header[24:])
	if _, err := pf.Write(header); err != nil {
		return nil, err
	}
	// Compute the differences, writing ctrl as we go
	pfbz2, err := bzip2.NewWriter(pf, nil)
	if err != nil {
		return nil, err
	}
	var scan, ln, lastscan, lastpos, lastoffset int
	newsize := len(newbin)
	oldsize := len(oldbin)

	var oldscore, scsc int

	for scan < newsize {
		oldscore = 0

		scan += ln
		scsc = scan
		for scan < newsize {
			scan++
			//ln = search(iii, oldbin, oldsize, newbin[scan:], newsize - scan, 0, oldsize, &pos)
		}
	}

	return nil, fmt.Errorf("not implemented")
}

func offtout(x int, buf []byte) {
	var y int
	if x < 0 {
		y = -x
	} else {
		y = x
	}
	buf[0] = byte(y % 256)
	y -= int(buf[0])
	y = y / 256
	buf[1] = byte(y % 256)
	y -= int(buf[1])
	y = y / 256
	buf[2] = byte(y % 256)
	y -= int(buf[2])
	y = y / 256
	buf[3] = byte(y % 256)
	y -= int(buf[3])
	y = y / 256
	buf[4] = byte(y % 256)
	y -= int(buf[4])
	y = y / 256
	buf[5] = byte(y % 256)
	y -= int(buf[5])
	y = y / 256
	buf[6] = byte(y % 256)
	y -= int(buf[6])
	y = y / 256
	buf[7] = byte(y % 256)

	if x < 0 {
		buf[7] |= 0x80
	}
}

func qsufsort(iii, vvv []int, buf []byte) {
	buckets := make([]int, 256)
	var i, h, ln int
	bufzise := len(buf)
	// [C] for (i = 0;i < 256;i++) buckets[i] = 0;
	// [C] for (i = 0;i < oldsize;i++) buckets[pold[i]]++;
	for i = 0; i < bufzise; i++ {
		buckets[buf[i]]++
	}
	// [C] for (i = 1;i < 256;i++) buckets[i] += buckets[i - 1];
	for i = 1; i < 256; i++ {
		buckets[i] += buckets[i-1]
	}
	// [C] for (i = 255;i > 0;i--) buckets[i] = buckets[i - 1];
	for i = 255; i > 0; i-- {
		buckets[i] = buckets[i-1]
	}
	buckets[0] = 0

	// [C] for (i = 0;i < oldsize;i++) I[++buckets[pold[i]]] = i;
	for i = 0; i < bufzise; i++ {
		buckets[buf[i]]++
		iii[buckets[buf[i]]] = i
	}
	iii[0] = bufzise
	// [C] for (i = 0;i < oldsize;i++) V[i] = buckets[pold[i]];
	for i = 0; i < bufzise; i++ {
		vvv[i] = buckets[buf[i]]
	}
	vvv[bufzise] = 0
	// [C] for (i = 1;i < 256;i++) if (buckets[i] == buckets[i - 1] + 1) I[buckets[i]] = -1;
	for i = 1; i < 256; i++ {
		if buckets[i] == buckets[i-1]+1 {
			iii[buckets[i]] = -1
		}
	}
	iii[0] = -1
	// [C] for (h = 1;I[0] != -(oldsize + 1);h += h) {
	for h = 1; iii[0] != (-(bufzise + 1)); h += h {
		ln = 0
		// [C] for (i = 0;i < oldsize + 1;) {
		i = 0
		for i < bufzise+1 {
			if iii[i] < 0 {
				ln -= (iii[i])
				i -= (iii[i])
			} else {
				if ln != 0 {
					iii[i-ln] = (-ln)
				}
				ln = (vvv[iii[i]] + 1 - (i))
				split(iii, vvv, i, ln, h)
				i += ln
				ln = 0
			}
		}
		if ln != 0 {
			iii[i-ln] = (-ln)
		}
	}

	for i = 0; i < bufzise+1; i++ {
		iii[vvv[i]] = (i)
	}
}

func split(iii, vvv []int, start, ln, h int) {
	var i, j, k, x int
	//var tmp int32

	if ln < 16 {
		for k = start; k < start+ln; k += j {
			j = 1
			x = (vvv[(iii[k] + (h))])
			for i = 1; k+i < start+ln; i++ {
				if vvv[(iii[k+i]+(h))] < (x) {
					x = int(vvv[int(iii[k+i]+(h))])
					j = 0
				}
				if vvv[int(iii[k+i]+(h))] == x {
					//tmp = iii[k+j]
					//iii[k+j]=iii[k+i]
					//iii[k+i]=tmp
					iii[k+j], iii[k+i] = iii[k+i], iii[k+j]
					j++
				}
				for i = 0; i < j; i++ {
					vvv[iii[k+i]] = k + j - 1
				}
				if j == 1 {
					iii[k] = -1
				}
			}
		}
		return
	}

	x = vvv[iii[start+ln/2]+h]
	var jj, kk int
	for i = start; i < start+ln; i++ {
		if vvv[iii[i]+h] < x {
			jj++
		} else if vvv[iii[i]+h] == x {
			kk++
		}
	}
	jj += start
	kk += jj

	i = start
	k = 0
	for i < jj {
		if vvv[iii[i]+h] < x {
			i++
		} else if vvv[iii[i]+h] == x {
			iii[i], iii[jj+j] = iii[jj+j], iii[i]
			j++
		} else {
			iii[i], iii[kk+k] = iii[kk+k], iii[i]
			k++
		}
	}
	for jj+j < kk {
		if vvv[iii[jj+j]+h] == x {
			j++
		} else {
			iii[jj+j], iii[kk+k] = iii[kk+k], iii[jj+j]
			k++
		}
	}
	if jj > start {
		split(iii, vvv, start, jj-start, h)
	}

	for i = 0; i < kk-jj; i++ {
		vvv[iii[jj+i]] = kk - 1
	}
	if jj == kk-1 {
		iii[jj] = -1
	}

	if start+ln > kk {
		split(iii, vvv, kk, start+ln-kk, h)
	}
}

func putWriter(target io.Writer, b []byte) error {
	lb := len(b)
	if lb < buffersize {
		n, err := target.Write(b)
		if err != nil {
			return err
		}
		if lb != n {
			return fmt.Errorf("%v of %v bytes written", n, lb)
		}
		return nil
	}
	offs := 0

	for offs < lb {
		n := min(buffersize, lb-offs)
		n2, err := target.Write(b[offs:n])
		if err != nil {
			return err
		}
		if n2 != n {
			return fmt.Errorf("%v of %v bytes written", offs+n2, lb)
		}
		offs += n
	}
	return nil
}

func copyReader(target []byte, rdr io.Reader) error {
	offs := 0
	buf := make([]byte, buffersize)
	for {
		nread, err := rdr.Read(buf)
		if nread > 0 {
			copy(target[offs:], buf[:nread])
			offs += nread
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
