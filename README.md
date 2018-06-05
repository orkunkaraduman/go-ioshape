# Go Traffic Shaper

[![GoDoc](https://godoc.org/github.com/orkunkaraduman/go-ioshape?status.svg)](https://godoc.org/github.com/orkunkaraduman/go-ioshape)

The repository provides `ioshape` package shapes I/O traffic using
token-bucket algorithm. It is used for creating bandwidth limiting applications,
needing traffic limiting or throttling or prioritization.

## Examples

### Limit copy operation simply

It limits copy operation to 2 MBps.

```go
    n, err := ioshape.CopyRate(dst, src, 2*1024*1024)
```

### Limit multiple operations with bucket

It limits two copy operation to 3MBps totally. Traffic will be balanced equally.

```go
    bucket := ioshape.NewBucketRate(3*1024*1024)
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        ioshape.CopyB(dst1, src1, bucket)
        wg.Done()
    }()
    wg.Add(1)
    go func() {
        ioshape.CopyB(dst2, src2, bucket)
        wg.Done()
    }()
    wg.Wait()
    bucket.Stop() // its necessary to free resources
```

### Limit multiple operations with burst and priority

It limits three copy operation to 5MBps totally. Traffic will be balanced with
given priorities.

```go
    bucket := ioshape.NewBucket()
    rate := 5*1024*1024 // the rate is 5MBps
    burst := rate*10 // the burst is ten times of the rate
    bucket.Set(rate, burst)
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        rr1 := &ioshape.Reader{R: src1, B: bucket, Pr: 0} // highest priority
        io.Copy(dst1, rr1)
        wg.Done()
    }
    wg.Add(1)
    go func() {
        rr2 := &ioshape.Reader{R: src2, B: bucket, Pr: 15} // lowest priority
        io.Copy(dst2, rr2)
        wg.Done()
    }
    wg.Add(1)
    go func() {
        rr3 := &ioshape.Reader{R: src3, B: bucket, Pr: 2} // higher priority
        io.Copy(dst3, rr3)
        wg.Done()
    }
    wg.Wait()
    bucket.Stop()
```
