package ioshape

import "io"

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
	for i := 0; n < l && err == nil; n += i {
		i = n + 32*1024
		if i > l {
			i = l
		}
		i, err = rr.R.Read(p[n:i])
		rr.B.getTokens(int64(i))
	}
	return
}
