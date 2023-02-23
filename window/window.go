package window

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wm sync.Mutex

func ClearScreen()      { fmt.Printf("\033[2J\033[0;0H") } // cls and goto 0:0
func HideCursor()       { fmt.Print("\033[?25l") }
func ShowCursor()       { fmt.Print("\033[?25h") }
func SaveCursor()       { fmt.Print("\033[s") }
func RestorCursor()     { fmt.Print("\033[u") }
func GoTo(y int) string { return fmt.Sprintf("\033[%d;0H", y) }

var logLines []string
var maxWidth, maxLines int
var Log chan string

func InitLog(ctx context.Context, wg *sync.WaitGroup, mw, ml int) {

	maxLines, maxWidth = ml, mw
	logLines = make([]string, 0, maxLines)
	Log = make(chan string)

	go logger(ctx, wg)
}

func PrintLog() {
	for _, l := range logLines {
		fmt.Println(l)
	}
}

func logger(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer func() {
		wg.Done()
		fmt.Println("logger done")
	}()

	for {
		select {
		case line := <-Log:
			timestamp := time.Now().Format("2006-01-02 15:04:05 ")

			l := (timestamp + line)
			if len(l) > maxWidth {
				l = l[:maxWidth]
			}
			if len(logLines) == maxLines {
				logLines = logLines[1:maxLines]
			}
			logLines = append(logLines, l)
		case <-ctx.Done():
			fmt.Println("logger ctx.Done")
			return
		}
	}
}
