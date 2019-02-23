package bsdiff

import (
	"fmt"
	"io"
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
	iii := make([]int32, len(oldbin)+1)
	vvv := make([]int32, len(oldbin)+1)
	qsufsort(iii, vvv, oldbin)
	// [C] free(V)
	vvv = nil
	return nil, fmt.Errorf("not implemented")
}

func qsufsort(iii, vvv []int32, buf []byte) {
	buckets := make([]int32, 256)
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
		iii[buckets[buf[i]]] = int32(i)
	}
	iii[0] = int32(bufzise)
	// [C] for (i = 0;i < oldsize;i++) V[i] = buckets[pold[i]];
	for i = 0; i < bufzise; i++ {
		vvv[i] = int32(buckets[buf[i]])
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
	for h = 1; iii[0] != int32(-(bufzise + 1)); h += h {
		ln = 0
		// [C] for (i = 0;i < oldsize + 1;) {
		i = 0
		for i < bufzise+1 {
			if iii[i] < 0 {
				ln -= int(iii[i])
				i -= int(iii[i])
			} else {
				if ln != 0 {
					iii[i-ln] = int32(-ln)
				}
				ln = int(vvv[iii[i]] + 1 - int32(i))
				split(iii, vvv, i, ln, h)
				i += ln
				ln = 0
			}
		}
		if ln != 0 {
			iii[i-ln] = int32(-ln)
		}
	}

	for i = 0; i < bufzise+1; i++ {
		iii[vvv[i]] = int32(i)
	}
}

func split(iii, vvv []int32, start, ln, h int) {
	var i, j, k, x, tmp, jj, kk int

	if ln < 16 {
		for k = start; k < start+ln; k += j {

		}
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
