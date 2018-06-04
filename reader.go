package ioshape

import (
	"io"
)

type Reader struct {
	R io.Reader
	B *Bucket
}

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
