package ioshape

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
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

	var wg sync.WaitGroup
	f := func(r io.Reader) {
		defer wg.Done()
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
		wg.Add(1)
		go f(resp.Body)
	}

	wg.Wait()
	bu.Stop()
}

func TestWriter(t *testing.T) {
	bu := NewBucket()
	bu.Set(128*1024, 0)
	size := 4 * 128 * 1024

	var wg sync.WaitGroup
	f := func(r io.Reader) {
		defer wg.Done()
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
		wg.Add(1)
		go f(resp.Body)
	}

	wg.Wait()
	bu.Stop()
}

func TestStopping(t *testing.T) {
	bu := NewBucket()
	bu.Set(128*1024, 0)
	size := 4 * 128 * 1024

	var wg sync.WaitGroup
	f := func(r io.Reader) {
		defer wg.Done()
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
		wg.Add(1)
		go f(resp.Body)
	}
	time.Sleep(8 * time.Second)
	bu.Stop()

	wg.Wait()
}
