package ioshape

import (
	"io"
	"time"
)

// Writer is a traffic shaper struct that implements io.Writer interface. A
// Writer writes to W by B.
// Priority changes between 0(highest) and 15(lowest).
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
		if k <= 0 {
			time.Sleep(1000 * 1000 * time.Microsecond / freq)
			continue
		}
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
