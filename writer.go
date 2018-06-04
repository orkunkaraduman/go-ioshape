package ioshape

import (
	"io"
)

// Writer is a traffic shaper struct that implements io.Writer interface. A
// Writer writes to W by B.
type Writer struct {
	W io.Writer // underlying reader
	B *Bucket   // bucket
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
		k := int(wr.B.giveTokens(int64(m)))
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
