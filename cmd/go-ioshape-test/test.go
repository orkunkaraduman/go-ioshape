package main

import (
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/orkunkaraduman/go-ioshape"
)

var count = int32(0)

func get(bu *ioshape.Bucket, size int) {
	atomic.AddInt32(&count, 1)
	defer atomic.AddInt32(&count, -1)
	resp, err := http.Get("http://ipv4.download.thinkbroadband.com/1GB.zip")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	r := &ioshape.Reader{B: bu, R: resp.Body}

	start := time.Now()
	b := make([]byte, size)
	_, err = io.ReadFull(r, b[:])
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Now().Sub(start))
}

func main() {
	bu := ioshape.NewBucket()
	bu.Set(16*1024, 0)
	for i := 0; i < 4; i++ {
		go get(bu, 64*1024)
	}
	time.Sleep(1 * time.Second)

	/*bu.Stop()
	fmt.Println("stopped")*/
	for count > 0 {
		time.Sleep(1 * time.Second)
	}
}
