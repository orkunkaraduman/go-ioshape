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
		j := n + chunkSize
		if j > l {
			j = l
		}
		i, err = rr.R.Read(p[n:j])
		rr.B.getTokens(int64(i))
	}
	return
}
