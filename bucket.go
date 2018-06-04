package ioshape

import (
	"sync"
	"sync/atomic"
	"time"
)

type bucketTokenRequest struct {
	count    int64
	callback chan int64
	priority int
}

// Bucket shapes traffic given rate, burst and Reader/Writer priorities.
type Bucket struct {
	tokens        int64
	n             int64
	k             int64
	b             int64
	m             int64
	setMu         sync.RWMutex
	ticker        *time.Ticker
	ticks         int64
	stopCh        chan struct{}
	stopped       int32
	tokenRequests chan *bucketTokenRequest
}

// NewBucket returns a new Bucket.
func NewBucket() (bu *Bucket) {
	bu = &Bucket{}
	bu.ticker = time.NewTicker(1000 * 1000 * time.Microsecond / freq)
	bu.stopCh = make(chan struct{}, 1)
	bu.tokenRequests = make(chan *bucketTokenRequest)
	go bu.timer()
	return
}

// NewBucketRate returns a new Bucket and sets rate.
func NewBucketRate(rate int64) (bu *Bucket) {
	bu = NewBucket()
	bu.SetRate(rate)
	return
}

func (bu *Bucket) timer() {
	for {
		select {
		case <-bu.stopCh:
			atomic.StoreInt32(&bu.stopped, 1)
			time.Sleep(10 * time.Millisecond)
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
			k := bu.k
			b := bu.b
			bu.setMu.RUnlock()
			bu.m = n / chunkDiv
			if bu.m == 0 {
				bu.m = 1
			}
			bu.tokens += n
			if bu.ticks%freq < k {
				bu.tokens++
			}
			if bu.tokens > b {
				bu.tokens = b
			}
			bu.ticks++
		case tokenRequest := <-bu.tokenRequests:
			count := tokenRequest.count
			if count > bu.tokens {
				count = bu.tokens
			}
			if count > bu.m {
				count = bu.m
			}
			if tokenRequest.priority > int(bu.ticks%freq) {
				count = 0
			}
			tokenRequest.callback <- count
			bu.tokens -= count
		}
	}
}

// Stop turns off a bucket. After Stop, bucket won't shape traffic. Stop
// must be call to free resources, after the bucket doesn't be needing.
func (bu *Bucket) Stop() {
	bu.ticker.Stop()
	select {
	case bu.stopCh <- struct{}{}:
	default:
	}
}

// Set sets buckets rate and burst in bytes per second. The burst should be
// greater or equal than the rate. Otherwise burst will be equal rate.
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
	bu.k = rate % freq
	bu.b = burst / freq
}

// SetRate sets rate and burst to the rate in bytes per second.
func (bu *Bucket) SetRate(rate int64) {
	bu.Set(rate, 0)
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

func (bu *Bucket) giveTokens(count int64) int64 {
	callback := make(chan int64)
	if count > 0 && bu.stopped == 0 {
		bu.tokenRequests <- &bucketTokenRequest{
			count:    count,
			callback: callback}
		return <-callback
	}
	return count
}

func (bu *Bucket) giveTokensPriority(count int64, priority int) int64 {
	callback := make(chan int64)
	if count > 0 && bu.stopped == 0 {
		bu.tokenRequests <- &bucketTokenRequest{
			count:    count,
			callback: callback,
			priority: priority}
		return <-callback
	}
	return count
}
