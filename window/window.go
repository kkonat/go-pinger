package window

import (
	"fmt"
	"sync"
	"time"
)

var wm sync.Mutex

func ClearScreen()      { fmt.Printf("\033[2J") }
func HideCursor()       { fmt.Print("\033[?25l") }
func ShowCursor()       { fmt.Print("\033[?25h") }
func SaveCursor()       { fmt.Print("\033[s") }
func RestorCursor()     { fmt.Print("\033[u") }
func GoTo(y int) string { return fmt.Sprintf("\033[%d;0H", y) }

var logLines []string
var maxWidth, maxLines int
var Log chan string

func InitLog(mw, ml int) {
	maxLines, maxWidth = ml, mw
	logLines = make([]string, 0, maxLines)
	Log = make(chan string)

	go logger()
}

func PrintLog() {
	for _, l := range logLines {
		fmt.Println(l)
	}
}

func logger() {
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
		}

	}
}
