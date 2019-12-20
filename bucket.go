package ioshape

import (
	"sync"
	"time"
)

type bucketTokenRequest struct {
	count    int64
	callback chan int64
	priority int
}

type bucketTokenReturn struct {
	count int64
}

// Bucket shapes traffic by given rate, burst and Reader/Writer priorities.
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
	stopMu        sync.Mutex
	stopped       bool
	tokenRequests chan *bucketTokenRequest
	tokenReturns  chan *bucketTokenReturn
}

// NewBucket returns a new Bucket.
func NewBucket() (bu *Bucket) {
	bu = &Bucket{}
	bu.ticker = time.NewTicker(time.Second / freq)
	bu.stopCh = make(chan struct{}, 1)
	bu.tokenRequests = make(chan *bucketTokenRequest)
	bu.tokenReturns = make(chan *bucketTokenReturn)
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
	var n, k, b, m int64
	for {
		select {
		case <-bu.stopCh:
			return

		case <-bu.ticker.C:
			bu.setMu.RLock()
			n = bu.n
			k = bu.k
			b = bu.b
			m = bu.m
			bu.setMu.RUnlock()
			bu.tokens += n
			if bu.ticks%freq < k {
				bu.tokens++
			}
			if bu.tokens > b {
				bu.tokens = b
			}
			bu.ticks++
			if bu.ticks > freq {
				bu.ticks = 0
			}

		case tokenRequest := <-bu.tokenRequests:
			count := tokenRequest.count
			if count > bu.tokens {
				count = bu.tokens
			}
			if count > m {
				count = m
			}
			if tokenRequest.priority > int((priorityScale*bu.ticks/freq)%priorityScale) {
				count = 0
			}
			tokenRequest.callback <- count
			bu.tokens -= count

		case tokenReturn := <-bu.tokenReturns:
			count := tokenReturn.count
			bu.tokens += count
			if bu.tokens > b {
				bu.tokens = b
			}
		}
	}
}

// Stop turns off a bucket. After Stop, bucket won't shape traffic. Stop
// must be call to free resources, after the bucket doesn't be needing.
func (bu *Bucket) Stop() {
	bu.stopMu.Lock()
	defer bu.stopMu.Unlock()
	if bu.stopped {
		return
	}
	bu.ticker.Stop()
	close(bu.stopCh)
	bu.stopped = true
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
	bu.b = burst + bu.n
	bu.m = bu.n / chunkDiv
	if bu.m == 0 {
		bu.m = 1
	}
}

// SetRate sets rate and burst to the rate in bytes per second.
func (bu *Bucket) SetRate(rate int64) {
	bu.Set(rate, 0)
}

func (bu *Bucket) getTokens(count int64, priority int) int64 {
	callback := make(chan int64)
	if count > 0 {
		select {
		case bu.tokenRequests <- &bucketTokenRequest{
			count:    count,
			callback: callback,
			priority: priority}:
			return <-callback
		case <-bu.stopCh:
			return count
		}
	}
	return count
}

func (bu *Bucket) giveTokens(count int64) {
	select {
	case bu.tokenReturns <- &bucketTokenReturn{
		count: count}:
	case <-bu.stopCh:
	}
}
