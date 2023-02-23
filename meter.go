package main

import (
	t "pinger/table"
)

const MPS = 4
const measurePeriod = 10

func runMeters(data []t.Data) {
	// ctx, cancel := context.WithTimeout(context.Background(), measurePeriod*time.Second)

	// defer func() {
	// 	w.Log <- "Context timed out. Calling cancel()."
	// 	cancel()
	// }()
	for i := 0; i < len(data); i++ {

	}
	// ticker := time.NewTicker(1000 * time.Millisecond / MPS).C

}

// type msrMsg struct {
// 	urlId  int
// 	v      uint64
// 	status string
// }

// func measurer(ctx context.Context, index int, ch chan<- msrMsg) {
// 	var duration uint64
// 	status := ""

// 	for {
// 		select {

// 		case <-ctx.Done():
// 			w.Log <- fmt.Sprintf("terminating: %d", index)
// 			return

// 		case <-ticker:

// 			start := time.Now()
// 			req, err := http.NewRequestWithContext(ctx, "GET", urls[index], nil)
// 			if err == nil {
// 				_, err = http.DefaultClient.Do(req)
// 				if err != nil {
// 					status = fmt.Sprintf("%s", err)
// 				} else {
// 					duration = uint64(time.Since(start).Milliseconds())
// 					status = "OK"
// 				}
// 			} else {
// 				status = fmt.Sprintf("%s", err)
// 			}
// 			ch <- msrMsg{urlId: index, v: duration, status: status}
// 		}
// 	}
// }

// func (d *tableItem) update(v uint64) {
// 	if d.count == 0 {
// 		d.min = v
// 		d.max = v
// 	} else {
// 		if d.min > v {
// 			d.min = v
// 		} else if d.max < v {
// 			d.max = v
// 		}
// 	}
// 	d.lastVal = v
// 	d.sum += v
// 	d.count++
// }

// 	ch := make(chan msrMsg)

// 	w.Log <- "Launching workers..."

// 	table.displayHeader()
// 	for i := range urls {
// 		table.displayRow(i)

// 		go measurer(ctx, i, ch)
// 	}

// 	var meas msrMsg

// 	w.Log <- "Starting..."

// loop:
// 	for {
// 		select {
// 		case meas = <-ch:
// 			table.update(meas)
// 			table.displayRow(meas.urlId)

// 		case <-ctx.Done():
// 			w.Log <- "Context expired. Breaking out of loop"
// 			break loop
// 		}
// 	}
