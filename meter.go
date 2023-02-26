package main

import (
	"context"
	"fmt"
	"math/rand"
	w "pinger/window"
	"time"
)

const MPS = 2
const minRepeatTime = 777

type measResult struct {
	val  uint64
	indx int
}

// simulate
func measure(ctx context.Context, resultStream chan<- measResult, i int) {
	start := time.Now()
	timeout := time.After(20 + time.Duration(rand.Intn(3000))*time.Millisecond)

	select {
	case <-timeout:
		dur := uint64(time.Since(start).Milliseconds())
		resultStream <- measResult{val: dur, indx: i}
		return
	case <-ctx.Done():
		return
	}
}

func runMeters(ctx context.Context, data []pinger) {
	w.Log <- "running  meters..."
	resultStream := make(chan measResult, len(data))
	//resultStream := make(chan measResult)

	for i := range data {
		go measure(ctx, resultStream, i)
	}

	for {
		select {
		case r := <-resultStream:
			data[r.indx].setVal(uint64(r.val))
			if r.val < minRepeatTime {
				delay := minRepeatTime - r.val
				w.Log <- "delaying..."
				time.Sleep(time.Millisecond * time.Duration(delay))
			}
			w.Log <- fmt.Sprintf("respawning measurement after delay: %d", r.val)
			go measure(ctx, resultStream, r.indx)

		case <-ctx.Done():
			close(resultStream)
			return
		}
	}
}

// TODO various pinger options
// https://medium.com/@deeeet/trancing-http-request-latency-in-golang-65b2463f548c
// https://stackoverflow.com/questions/41423637/go-ping-library-for-unprivileged-icmp-ping-in-golang
// https://github.com/davecheney/httpstat
// https://pkg.go.dev/search?q=netstat
