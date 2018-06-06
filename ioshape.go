/*
Package ioshape provides I/O structures and functions for Traffic Shaping using
token-bucket algorithm.
*/
package ioshape

import (
	"errors"
	"io"
)

const (
	freq          = 16
	freqMul       = 4
	priorityScale = 16
	chunkDiv      = 1
	chunkSize     = 32 * 1024
)

// ErrOutOfRange is the error used for the result of r/w is out of range.
var ErrOutOfRange = errors.New("out of range")

// CopyB is identical to io.Copy except that it shapes traffic by b *Bucket.
func CopyB(dst io.Writer, src io.Reader, b *Bucket) (written int64, err error) {
	return io.Copy(dst, &Reader{R: src, B: b})
}

// CopyBN is identical to io.CopyN except that it shapes traffic by b *Bucket.
func CopyBN(dst io.Writer, src io.Reader, b *Bucket, n int64) (written int64, err error) {
	return io.CopyN(dst, &Reader{R: src, B: b}, n)
}

// CopyRate is identical to io.Copy except that it shapes traffic with rate
// in bytes per second.
func CopyRate(dst io.Writer, src io.Reader, rate int64) (written int64, err error) {
	b := NewBucketRate(rate)
	written, err = io.Copy(dst, &Reader{R: src, B: b})
	b.Stop()
	return
}

// CopyRateN is identical to io.CopyN except that it shapes traffic with rate
// in bytes per second.
func CopyRateN(dst io.Writer, src io.Reader, rate int64, n int64) (written int64, err error) {
	b := NewBucketRate(rate)
	written, err = io.CopyN(dst, &Reader{R: src, B: b}, n)
	b.Stop()
	return
}
