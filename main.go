package main

import (
	"fmt"
	t "pinger/table"
)

type tableItem struct {
	url     string `td:"name='website',  w= 30,   sep='|'"`
	lastVal uint64 `td:"name='ping (ms)',w= 9, R, sep='>'"`
	min     uint64 `td:"w= 8, R"`
	max     uint64 `td:"w= 8, R"`
	avg     uint64 `td:"w= 8, R"`
	status  string `td:"name='errors',   w= 60,   sep='|'"`
	count   uint64
	sum     uint64
	valid   bool
}

func (ti tableItem) GenCells() *[]string {
	if ti.valid {
		return &[]string{
			ti.url,
			fmt.Sprintf("%4.d", ti.lastVal),
			fmt.Sprintf("%4.d", ti.min),
			fmt.Sprintf("%4.d", ti.max),
			fmt.Sprintf("%8.2f", float64(ti.sum)/float64(ti.count)),
			ti.status,
		}
	} else {
		return &[]string{ti.url, "?", "?", "?", "?", "connecting..."}
	}
}

// type tablit struct {
// }

// func (ti tablit) GenCells() *[]string { return &[]string{"", ""} }

func main() {
	// dd := tablit{}
	// var td []t.Data = []t.Data{dd, dd}
	// fmt.Println("%T", td)

	table := t.New()
	var tableData []t.Data = []t.Data{
		tableItem{url: "http://www.onet.pl"},
		tableItem{url: "http://www.gazeta.pl"},
		tableItem{url: "http://www.wsj.com"},
		tableItem{url: "http://www.google.com"},
		tableItem{url: "http://nonexistent.com"},
		tableItem{url: "http://www.icm.edu.pl"},
	}

	table.Init(tableData, true)
	table.Print()

	// func (ti *tableItem) GenRowData(row int) []string {
	// 	return []string{
	// 		fmt.Sprintf("%3d", row),
	// 		fmt.Sprintf("%s", tableData[row].url)
	// 	}
	// }
	// 	ctx, cancel := context.WithTimeout(context.Background(), measurePeriod*time.Second)
	// 	defer func() {
	// 		w.Log <- "Context timed out. Calling cancel()."
	// 		cancel()
	// 	}()

	// 	w.Init(len(urls) + 4)

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

}
