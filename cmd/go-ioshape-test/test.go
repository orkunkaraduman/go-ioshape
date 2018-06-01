package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/orkunkaraduman/go-ioshape"
)

func f(bu *ioshape.Bucket) {
	resp, err := http.Get("http://ipv4.download.thinkbroadband.com/1GB.zip")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	r := &ioshape.Reader{B: bu, R: resp.Body}

	start := time.Now()
	b := make([]byte, 1024*1024)
	_, err = io.ReadFull(r, b[:])
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Now().Sub(start))
}

func main() {
	bu := ioshape.NewBucket()
	bu.Set(128*1024, 0)
	go f(bu)
	go f(bu)
	go f(bu)
	time.Sleep(1 * time.Second)
	//bu.Stop()
	//fmt.Println("stopped")
	time.Sleep(100 * time.Second)
}
