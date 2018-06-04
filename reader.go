package ioshape

import (
	"io"
)

// Reader is a traffic shaper struct that implements io.Reader interface. A
// Reader reads from R by B.
type Reader struct {
	R io.Reader // underlying reader
	B *Bucket   // bucket
}

// Read reads from R by b.
func (rr *Reader) Read(p []byte) (n int, err error) {
	if rr.B == nil {
		n, err = rr.R.Read(p)
		return
	}

	l := len(p)
	m := l
	for n < l && err == nil {
		k := int(rr.B.giveTokens(int64(m)))
		var nn int
		nn, err = rr.R.Read(p[n : n+k])
		n += nn
		m -= nn
	}
	return
}
