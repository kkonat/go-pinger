package main

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// type msrMsg struct {
// 	urlId  int
// 	v      uint64
// 	status string
// }

// func measurer(ctx context.Context, index int, ch chan<- msrMsg) {
// 	var duration uint64
// 	status := ""

// 	ticker := time.NewTicker(1000 * time.Millisecond / MPS).C

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
