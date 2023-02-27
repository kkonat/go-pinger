package main

import (
	"context"
	"fmt"
	w "pinger/window"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

const MPS = 2
const minRepeatTime = 777

type measResult struct {
	tts  time.Duration
	indx int
	err  error
}

func measure(ctx context.Context, resultStream chan<- measResult, data []pinger, i int) {

	result := make(chan measResult)

	pinger, err := probing.NewPinger(data[i].url)

	pinger.SetPrivileged(true)
	if err != nil {
		w.Log <- fmt.Sprint("Error 1:", err)
		resultStream <- measResult{tts: 0, indx: i, err: err}
		return
	}

	pinger.Count = 3
	pinger.Interval = 200 * time.Millisecond
	pinger.Timeout = 3 * time.Second
	go func() {

		err = pinger.Run() // Blocks until finished.

		if err != nil {
			w.Log <- fmt.Sprint("Error 2:", err)
			resultStream <- measResult{tts: 0, indx: i, err: err}
			return
		}
		stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
		result <- measResult{tts: stats.AvgRtt, indx: i, err: nil}
	}()

	select {
	case res := <-result:
		resultStream <- res
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
		w.Log <- "pinging " + data[i].url + " ..."
		data[i].SetStatus("connecting...")
		go measure(ctx, resultStream, data, i)
	}

	for {
		select {
		case res := <-resultStream:
			if res.err != nil {
				data[res.indx].Invalidate(res.err.Error())
			} else {
				data[res.indx].SetVal(float64(res.tts.Nanoseconds()) / 1.0e6)
			}

			w.Log <- "repinging " + data[res.indx].url + " ..."
			time.Sleep(time.Millisecond * 500)
			go measure(ctx, resultStream, data, res.indx)
			continue
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
