package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	t "pinger/table"
	w "pinger/window"
	"sync"
	"syscall"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

const howLong = 5

func main() {

	var (
		data []t.Data = []t.Data{
			tableItem{url: "http://www.onet.pl"},
			tableItem{url: "http://www.gazeta.pl"},
			tableItem{url: "http://www.wsj.com"},
			tableItem{url: "http://www.google.com"},
			tableItem{url: "http://nonexistent.com"},
			tableItem{url: "http://www.icm.edu.pl"},
		}
		table = t.New().Init(data, true)

		// closing     = make(chan struct{})
		wg          sync.WaitGroup
		ctx, cancel = context.WithCancel(context.Background())
	)

	go func() {
		detect := make(chan os.Signal, 1)
		signal.Notify(detect, syscall.SIGTERM, os.Interrupt)
		<-detect
		fmt.Println(" Ctrl+C")
		cancel()
		// close(closing)
	}()

	w.ClearScreen()
	w.InitLog(ctx, &wg, table.Width, 5)
	w.SaveCursor()
	w.HideCursor()
	defer w.ShowCursor()

	//go runMeters(data)
	go RunDisplay(ctx, &wg, table)

loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("closing")
			break loop
		case <-time.After(time.Duration(howLong) * time.Second):
			fmt.Println("timeout")
			cancel()

			break loop
		}
	}

	fmt.Println("waiting for all goroutinest to finish...")
	wg.Wait()
	fmt.Println("bye")
}

const FPS = 5

// TODO: implment graceful shutdown
// https://www.rudderstack.com/blog/implementing-graceful-shutdown-in-go/
// https://justbartek.ca/p/golang-context-wg-go-routines/
func RunDisplay(ctx context.Context, wg *sync.WaitGroup, table *t.Table) {
	wg.Add(1)
	defer func() {
		fmt.Println("RunDisplay: done")
		wg.Done()
	}()

	var (
		progress = []string{"    ", ".   ", "..  ", "... ", "...."}
		ticker   = time.NewTicker(1000 * time.Millisecond / FPS).C // frames per second
		clock    = time.NewTicker(time.Second).C                   // once per second

		tick = 0
	)

	for {
		select {
		case <-ticker:
			w.RestorCursor()
			fmt.Println("Pinger running", progress[tick%len(progress)])
			table.Print()
			w.PrintLog()
			tick++
		case <-ctx.Done():
			fmt.Println("RunDisplay: ctx.Done")
			return
		case <-clock:
			w.Log <- fmt.Sprint("Performing important stuff ", tick/FPS)
		}
	}
}
