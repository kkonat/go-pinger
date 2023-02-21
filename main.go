package main

import (
	"context"
	"fmt"
	"net/http"
	w "pinger/window"
	"time"
)

const measurePeriod = 10
const MPS = 4 // Measurements per second
var urls []string = []string{
	"http://www.onet.pl",
	"http://www.gazeta.pl",
	"http://www.wsj.com",
	"http://www.google.com",
	"http://nonexistent.com",
	"http://www.icm.edu.pl",
}

type msrMsg struct {
	urlId  int
	v      uint64
	status string
}

func measurer(ctx context.Context, index int, ch chan<- msrMsg) {
	var duration uint64
	status := ""

	ticker := time.NewTicker(1000 * time.Millisecond / MPS).C

	for {
		select {

		case <-ctx.Done():
			w.Log <- fmt.Sprintf("terminating: %d", index)
			return

		case <-ticker:

			start := time.Now()
			req, err := http.NewRequestWithContext(ctx, "GET", urls[index], nil)
			if err == nil {
				_, err = http.DefaultClient.Do(req)
				if err != nil {
					status = fmt.Sprintf("%s", err)
				} else {
					duration = uint64(time.Since(start).Milliseconds())
					status = "OK"
				}
			} else {
				status = fmt.Sprintf("%s", err)
			}
			ch <- msrMsg{urlId: index, v: duration, status: status}
		}
	}
}

func main() {

	table := NewTable()

	ctx, cancel := context.WithTimeout(context.Background(), measurePeriod*time.Second)
	defer func() {
		w.Log <- "Context timed out. Calling cancel()."
		cancel()
	}()

	w.Init(len(urls) + 4)

	ch := make(chan msrMsg)

	w.Log <- "Launching workers..."

	table.displayHeader()
	for i := range urls {
		table.displayRow(i)

		go measurer(ctx, i, ch)
	}

	var meas msrMsg

	w.Log <- "Starting..."

loop:
	for {
		select {
		case meas = <-ch:
			table.update(meas)
			table.displayRow(meas.urlId)

		case <-ctx.Done():
			w.Log <- "Context expired. Breaking out of loop"
			break loop
		}
	}

}
