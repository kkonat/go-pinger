package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	t "pinger/table"
	w "pinger/window"
	"syscall"
	"time"
)

// var wg sync.WaitGroup
var flowControl chan string

func main() {
	w.ClearScreen()
	fmt.Print(w.GoTo(0))
	table := t.New()
	var data []t.Data = []t.Data{
		tableItem{url: "http://www.onet.pl"},
		tableItem{url: "http://www.gazeta.pl"},
		tableItem{url: "http://www.wsj.com"},
		tableItem{url: "http://www.google.com"},
		tableItem{url: "http://nonexistent.com"},
		tableItem{url: "http://www.icm.edu.pl"},
	}
	flowControl = make(chan string)
	table.Init(data, true)

	//go runMeters(data)
	w.SaveCursor()
	w.HideCursor()
	defer w.ShowCursor()

	go runDisplay(table, flowControl, 15)
	w.Log <- "Log starts here:"

	fmt.Println("Waiting")
	// wg.Wait()

	fmt.Println(<-flowControl)
	fmt.Println("finito")
}

const FPS = 5

// TODO: implment graceful shutdown
// https://www.rudderstack.com/blog/implementing-graceful-shutdown-in-go/
func runDisplay(table *t.Table, ch chan string, howlong int) {
	// wg.Add(1)
	// defer wg.Done()
	progress := []string{"    ", ".   ", "..  ", "... ", "...."}
	ticker := time.NewTicker(1000 * time.Millisecond / FPS).C // frames per second
	clock := time.NewTicker(time.Second).C                    // once per second

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(howlong)*time.Second)
	defer cancel()

	w.InitLog(table.Width, 5)

	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt, syscall.SIGTERM)

	tick := 0

	for {
		select {
		case <-ticker:
			w.RestorCursor()
			fmt.Println("Pinger running", progress[tick%len(progress)])
			table.Print()
			w.PrintLog()
			tick++
		case <-ctrlC:
			fmt.Println("Ctrl+C")
			flowControl <- "Ctrl+C"
			return
		case <-ctx.Done():
			fmt.Println("Done")
			flowControl <- "Done"
			return
		case <-clock:
			w.Log <- fmt.Sprint("Performing important stuff ", tick/FPS)
		}
	}
}
