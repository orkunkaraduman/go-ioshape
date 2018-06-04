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
	for i := 0; n < l && err == nil; n += i {
		j := n + chunkSize
		if j > l {
			j = l
		}
		i, err = wr.W.Write(p[n:j])
		if i != j-n {
			err = io.ErrShortWrite
		}
		wr.B.getTokens(int64(i))
	}
	return
}
