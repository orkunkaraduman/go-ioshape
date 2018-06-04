package ioshape

import (
	"io"
)

const (
	freq     = 16
	chunkDiv = 16
)

func CopyB(dst io.Writer, src io.Reader, b *Bucket) (written int64, err error) {
	return io.Copy(dst, &Reader{R: src, B: b})
}

func CopyBN(dst io.Writer, src io.Reader, b *Bucket, n int64) (written int64, err error) {
	return io.CopyN(dst, &Reader{R: src, B: b}, n)
}

func CopyRate(dst io.Writer, src io.Reader, rate int64) (written int64, err error) {
	return io.Copy(dst, &Reader{R: src, B: NewBucketRate(rate)})
}

func CopyRateN(dst io.Writer, src io.Reader, rate int64, n int64) (written int64, err error) {
	return io.CopyN(dst, &Reader{R: src, B: NewBucketRate(rate)}, n)
}
