package gocli

import (
	"fmt"
	"os"
)

func (t *Terminal) FnExitProgram() {
	fmt.Println()
	os.Exit(0)
}

func (t *Terminal) FnClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (t *Terminal) FnDeleteLastLine() {
	fmt.Print("\033[A")
	fmt.Printf("\033[%dG", 1)
	fmt.Print("\033[K")
}
