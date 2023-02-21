package window

import (
	"fmt"
	"log"
)

var consoleLine int = 0

func Init(ConsoleLine int) {
	consoleLine = ConsoleLine
}

func Gotoxy(x, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

func ClearScreen() {
	fmt.Printf("\033[2J")
}

func Log(line string) {
	Gotoxy(0, consoleLine)
	log.Print(line)
	consoleLine++
}
