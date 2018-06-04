package ioshape

import (
	"io"
)

// Writer is a traffic shaper struct that implements io.Writer interface. A
// Writer writes to W by B.
// Priority changes between 0 and 15: 0 is higher, 15 is lower.
type Writer struct {
	W  io.Writer // underlying reader
	B  *Bucket   // bucket
	Pr int       // priority
}

// Write writes to W by b.
func (wr *Writer) Write(p []byte) (n int, err error) {
	if wr.B == nil {
		n, err = wr.W.Write(p)
		return
	}

	l := len(p)
	m := l
	for n < l && err == nil {
		k := int(wr.B.giveTokensPriority(int64(m), wr.Pr))
		var nn int
		nn, err = wr.W.Write(p[n : n+k])
		if nn != k {
			err = io.ErrShortWrite
		}
		n += nn
		m -= nn
	}
	return
}
