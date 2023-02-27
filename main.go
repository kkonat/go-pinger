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

const howLong = 60

func main() {

	var (
		data []pinger = []pinger{
			{url: "www.onet.pl"},
			{url: "www.gazeta.pl"},
			{url: "www.wsj.com"},
			{url: "www.google.com"},
			{url: "nonexistent.com"},
			{url: "www.icm.edu.pl"},
		}
		table = t.New(data, true)

		ctx, cancel = context.WithCancel(context.Background())
	)
	defer cancel()

	abort := make(chan struct{})

	go func() {
		detect := make(chan os.Signal, 1)
		signal.Notify(detect, syscall.SIGTERM, os.Interrupt)
		<-detect
		close(abort)
	}()

	w.ClearScreen()
	w.SaveCursor()
	w.HideCursor()

	w.StartLog(ctx, table.GetWidth(), 5)

	go runMeters(ctx, data)

	go runDisplay(ctx, table)

	select {
	case <-abort:
	case <-time.After(time.Duration(howLong) * time.Second):
	}
	fmt.Println("Bye")
	w.ShowCursor()
}

const FPS = 5

func runDisplay(ctx context.Context, table *t.Table[pinger]) {

	var (
		progress = []string{"    ", ".   ", "..  ", "... ", "...."}
		ticker   = time.NewTicker(1000 * time.Millisecond / FPS).C // frames per second

		tick = 0
	)

loop:
	for {
		select {
		case <-ticker:
			w.RestorCursor()
			fmt.Println("Pinger running", progress[tick%len(progress)])
			table.Print()
			w.PrintLog()
			tick++
		case <-ctx.Done():
			break loop
		}
	}
}
