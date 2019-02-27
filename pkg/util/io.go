package util

import (
	"fmt"
	"io"
	"sync"
)

const (
	buffersize = 1024 * 16
)

// PutWriter writes all content from b to target
func PutWriter(target io.Writer, b []byte) error {
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
		n := Min(buffersize, lb-offs)
		n2, err := target.Write(b[offs : offs+n])
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

// Min int
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// BufWriter is byte slice buffer that implements io.WriteSeeker
type BufWriter struct {
	lock sync.Mutex
	buf  []byte
	pos  int
}

// Write the contents of p and return the bytes written
func (m *BufWriter) Write(p []byte) (n int, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.buf == nil {
		m.buf = make([]byte, 0)
		m.pos = 0
	}
	minCap := m.pos + len(p)
	if minCap > cap(m.buf) { // Make sure buf has enough capacity:
		buf2 := make([]byte, len(m.buf), minCap+len(p)) // add some extra
		copy(buf2, m.buf)
		m.buf = buf2
	}
	if minCap > len(m.buf) {
		m.buf = m.buf[:minCap]
	}
	copy(m.buf[m.pos:], p)
	m.pos += len(p)
	return len(p), nil
}

// Seek to a position on the byte slice
func (m *BufWriter) Seek(offset int64, whence int) (int64, error) {
	newPos, offs := 0, int(offset)
	switch whence {
	case io.SeekStart:
		newPos = offs
	case io.SeekCurrent:
		newPos = m.pos + offs
	case io.SeekEnd:
		newPos = len(m.buf) + offs
	}
	if newPos < 0 {
		return 0, fmt.Errorf("negative result pos")
	}
	m.pos = newPos
	return int64(newPos), nil
}

// Len returns the length of the internal byte slice
func (m *BufWriter) Len() int {
	return len(m.buf)
}

// Bytes return a copy of the internal byte slice
func (m *BufWriter) Bytes() []byte {
	b2 := make([]byte, len(m.buf))
	copy(b2, m.buf)
	return b2
}
