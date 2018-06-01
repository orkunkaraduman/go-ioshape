package ioshape

import (
	"sync"
	"sync/atomic"
	"time"
)

type bucketTokenRequest struct {
	count    int64
	callback chan int64
}

type Bucket struct {
	tokens        int64
	n             int64
	b             int64
	setMu         sync.RWMutex
	ticker        *time.Ticker
	stopCh        chan struct{}
	stopped       int32
	tokenRequests chan *bucketTokenRequest
}

func NewBucket() (bu *Bucket) {
	bu = &Bucket{}
	bu.ticker = time.NewTicker(1000 * time.Millisecond / freq)
	bu.stopCh = make(chan struct{}, 1)
	bu.tokenRequests = make(chan *bucketTokenRequest)
	go bu.timer()
	return
}

func (bu *Bucket) timer() {
	for {
		select {
		case <-bu.stopCh:
			atomic.StoreInt32(&bu.stopped, 1)
			time.Sleep(1 * time.Millisecond)
			for ok := true; ok; {
				select {
				case tokenRequest := <-bu.tokenRequests:
					tokenRequest.callback <- tokenRequest.count
				default:
					ok = false
				}
			}
			return
		case <-bu.ticker.C:
			bu.setMu.RLock()
			n := bu.n
			b := bu.b
			bu.setMu.RUnlock()

			bu.tokens += n
			if bu.tokens > b {
				bu.tokens = b
			}

			for ok := true; ok; {
				select {
				case tokenRequest := <-bu.tokenRequests:
					count := tokenRequest.count
					if count > bu.tokens {
						count = bu.tokens
					}
					tokenRequest.callback <- count
					bu.tokens -= count
					if bu.tokens <= 0 {
						ok = false
					}
				default:
					ok = false
				}
			}
		}
	}
}

func (bu *Bucket) Stop() {
	select {
	case bu.stopCh <- struct{}{}:
	default:
	}
	bu.ticker.Stop()
}

func (bu *Bucket) Set(rate, burst int64) {
	if rate < 0 {
		return
	}
	bu.setMu.Lock()
	defer bu.setMu.Unlock()
	if rate > burst {
		burst = rate
	}
	bu.n = rate / freq
	bu.b = burst / freq
}

func (bu *Bucket) getTokens(count int64) {
	callback := make(chan int64)
	for count > 0 && bu.stopped == 0 {
		bu.tokenRequests <- &bucketTokenRequest{
			count:    count,
			callback: callback}
		count -= <-callback
	}
}
