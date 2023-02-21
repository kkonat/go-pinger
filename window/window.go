package window

import (
	"fmt"
	"sync"
	"time"
)

var wm sync.Mutex
var consoleLine int = 0

var Log chan string

func Init(ConsoleLine int) {
	consoleLine = ConsoleLine
	Log = make(chan string)
	go LogPrinter()
}

func ClearScreen() {
	fmt.Printf("\033[2J")
}

func PrintLine(line int, text string) {
	fmt.Print(goTo(line), text)
}

func LogPrinter() {
	for line := range Log {
		timestamp := time.Now().Format("2006-01-02 15:04:05 ")
		fmt.Print(goTo(consoleLine), timestamp, line)
		consoleLine++

	}
}

func goTo(y int) string {
	return fmt.Sprintf("\033[%d;0H", y)
}
