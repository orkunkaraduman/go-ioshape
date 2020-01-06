package ioshape

import (
	"fmt"
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
	doneCh        chan struct{}
	tokenRequests chan *bucketTokenRequest
	tokenReturns  chan *bucketTokenReturn
}

// NewBucket returns a new Bucket.
func NewBucket() (bu *Bucket) {
	bu = &Bucket{
		ticker:        time.NewTicker(time.Second / freq),
		stopCh:        make(chan struct{}),
		doneCh:        make(chan struct{}),
		tokenRequests: make(chan *bucketTokenRequest),
		tokenReturns:  make(chan *bucketTokenReturn),
	}
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
	defer close(bu.doneCh)
	var n, k, b, m int64
	tokenRequests := bu.tokenRequests
	var pendingRequest *bucketTokenRequest
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

			if pendingRequest != nil {
				done := bu.handleReaquest(pendingRequest, m)
				if !done {
					fmt.Println("Warning: This should not happen")
					pendingRequest.callback <- pendingRequest.count
				}
				tokenRequests = bu.tokenRequests
				pendingRequest = nil
			}

		case tokenRequest := <-tokenRequests:
			done := bu.handleReaquest(tokenRequest, m)
			if !done {
				tokenRequests = nil
				pendingRequest = tokenRequest
			}

		case tokenReturn := <-bu.tokenReturns:
			count := tokenReturn.count
			bu.tokens += count
			if bu.tokens > b {
				bu.tokens = b
			}
		}
	}
}

// handleReaquest may only be called from the timer loop
func (bu *Bucket) handleReaquest(r *bucketTokenRequest, m int64) bool {
	count := r.count
	if count > bu.tokens {
		count = bu.tokens
	}
	if count > m {
		count = m
	}
	if count == 0 {
		return false
	}
	if r.priority > int((priorityScale*bu.ticks/freq)%priorityScale) {
		count = 0
	}
	r.callback <- count
	bu.tokens -= count
	return true
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
			select {
			case c := <-callback:
				return c
			case <-bu.doneCh:
				return count
			}
		case <-bu.doneCh:
			return count
		}
	}
	return count
}

func (bu *Bucket) giveTokens(count int64) {
	select {
	case bu.tokenReturns <- &bucketTokenReturn{
		count: count}:
	case <-bu.doneCh:
	}
}
