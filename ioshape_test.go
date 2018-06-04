package ioshape

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

const (
	testURL = "http://ipv4.download.thinkbroadband.com/1GB.zip"
)

func TestReader(t *testing.T) {
	bu := NewBucket()
	bu.Set(128*1024, 0)
	size := 4 * 128 * 1024
	count := int32(0)

	f := func(r io.Reader) {
		atomic.AddInt32(&count, 1)
		defer atomic.AddInt32(&count, -1)
		defer func() {
			if rr, ok := r.(io.Closer); ok {
				rr.Close()
			}
		}()
		start := time.Now()
		rr := &Reader{R: r, B: bu}
		_, err := io.CopyN(ioutil.Discard, rr, int64(size))
		if err != nil {
			panic(err)
		}
		fmt.Println(time.Now().Sub(start))
	}

	for i := 0; i < 4; i++ {
		resp, err := http.Get(testURL)
		if err != nil {
			panic(err)
		}
		go f(resp.Body)
	}
	time.Sleep(10 * time.Millisecond)

	for count > 0 {
		time.Sleep(10 * time.Millisecond)
	}
}

func TestWriter(t *testing.T) {
	bu := NewBucket()
	bu.Set(128*1024, 0)
	size := 4 * 128 * 1024
	count := int32(0)

	f := func(r io.Reader) {
		atomic.AddInt32(&count, 1)
		defer atomic.AddInt32(&count, -1)
		defer func() {
			if wr, ok := r.(io.Closer); ok {
				wr.Close()
			}
		}()
		start := time.Now()
		wr := &Writer{W: ioutil.Discard, B: bu}
		_, err := io.CopyN(wr, r, int64(size))
		if err != nil {
			panic(err)
		}
		fmt.Println(time.Now().Sub(start))
	}

	for i := 0; i < 4; i++ {
		resp, err := http.Get(testURL)
		if err != nil {
			panic(err)
		}
		go f(resp.Body)
	}
	time.Sleep(10 * time.Millisecond)

	for count > 0 {
		time.Sleep(10 * time.Millisecond)
	}
}

func TestStopping(t *testing.T) {
	bu := NewBucket()
	bu.Set(128*1024, 0)
	size := 4 * 128 * 1024
	count := int32(0)

	f := func(r io.Reader) {
		atomic.AddInt32(&count, 1)
		defer atomic.AddInt32(&count, -1)
		defer func() {
			if rr, ok := r.(io.Closer); ok {
				rr.Close()
			}
		}()
		start := time.Now()
		rr := &Reader{R: r, B: bu}
		_, err := io.CopyN(ioutil.Discard, rr, int64(size))
		if err != nil {
			panic(err)
		}
		fmt.Println(time.Now().Sub(start))
	}

	for i := 0; i < 4; i++ {
		resp, err := http.Get(testURL)
		if err != nil {
			panic(err)
		}
		go f(resp.Body)
	}
	time.Sleep(10 * time.Millisecond)

	time.Sleep(8 * time.Second)
	bu.Stop()

	for count > 0 {
		time.Sleep(10 * time.Millisecond)
	}
}
