package ioshape

import "io"

type Writer struct {
	W io.Writer
	B *Bucket
}

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
