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

var wg sync.WaitGroup

func main() {

	var (
		data []pinger = []pinger{
			{url: "http://www.onet.pl", min: 1, max: 2, avg: 3.14, valid: true},
			{url: "http://www.gazeta.pl"},
			{url: "http://www.wsj.com"},
			{url: "http://www.google.com"},
			{url: "http://nonexistent.com"},
			{url: "http://www.icm.edu.pl"},
		}
		table = t.New(data, true)

		ctx, cancel = context.WithCancel(context.Background())
	)

	go func() {
		detect := make(chan os.Signal, 1)
		signal.Notify(detect, syscall.SIGTERM, os.Interrupt)
		<-detect
		fmt.Println(" Ctrl+C")
		cancel()
	}()

	w.ClearScreen()
	w.InitLog(ctx, table.MaxWidth(), 5)
	w.SaveCursor()
	w.HideCursor()
	defer w.ShowCursor()

	//go runMeters(data)
	go RunDisplay(ctx, table)

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
	wg.Wait()
	fmt.Println("bye")
}

const FPS = 5

// TODO: implment graceful shutdown
// https://www.rudderstack.com/blog/implementing-graceful-shutdown-in-go/
// https://justbartek.ca/p/golang-context-wg-go-routines/
func RunDisplay(ctx context.Context, table *t.Table[pinger]) {

	var (
		progress = []string{"    ", ".   ", "..  ", "... ", "...."}
		ticker   = time.NewTicker(1000 * time.Millisecond / FPS).C // frames per second
		clock    = time.NewTicker(time.Second).C                   // once per second

		tick = 0
	)

	for {
		select {
		case <-ticker:
			wg.Add(1)
			w.RestorCursor()
			fmt.Println("Pinger running", progress[tick%len(progress)])
			table.Print()
			w.PrintLog()
			tick++
			wg.Done()
		case <-ctx.Done():

			return
		case <-clock:
			w.Log <- fmt.Sprint("Performing important stuff ", tick/FPS)
		}
	}
}
