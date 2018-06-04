package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/orkunkaraduman/go-ioshape"
)

var count = int32(0)

func testReader(bu *ioshape.Bucket, size int) {
	atomic.AddInt32(&count, 1)
	defer atomic.AddInt32(&count, -1)
	resp, err := http.Get("http://ipv4.download.thinkbroadband.com/1GB.zip")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	r := &ioshape.Reader{R: resp.Body, B: bu}

	start := time.Now()
	b := make([]byte, size)
	_, err = io.ReadFull(r, b[:])
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Now().Sub(start))
}

func testWriter(bu *ioshape.Bucket, size int) {
	atomic.AddInt32(&count, 1)
	defer atomic.AddInt32(&count, -1)
	resp, err := http.Get("http://ipv4.download.thinkbroadband.com/1GB.zip")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	w := &ioshape.Writer{W: ioutil.Discard, B: bu}

	start := time.Now()
	_, err = io.CopyN(w, resp.Body, int64(size))
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Now().Sub(start))
}

func main() {
	bu := ioshape.NewBucket()
	bu.Set(128*1024, 0)
	for i := 0; i < 4; i++ {
		go testWriter(bu, 512*1024)
	}
	time.Sleep(1 * time.Second)
	/*bu.Stop()
	fmt.Println("stopped")*/
	for count > 0 {
		time.Sleep(1 * time.Second)
	}
}
