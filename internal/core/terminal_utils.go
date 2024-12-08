package gocli

import (
	"fmt"

	gu "github.com/vcharco/gocli/internal/utils"
)

func (t *Terminal) replaceLine(userInput *string, text string) {
	t.cleanLine()
	fmt.Print(text)
	t.moveCursorToPos(len(text))
	t.cursorPos = len(text)
	*userInput = text
}

func (t *Terminal) cleanLine() {
	t.moveCursorToPos(0)
	fmt.Print("\033[K")
}

func (t *Terminal) cleanNextLine() {
	t.cleanNextLineAndStay()
	fmt.Print("\033[1A")
	t.moveCursorToPos(t.cursorPos)
}

func (t *Terminal) cleanNextLineAndStay() {
	fmt.Println()
	fmt.Printf("\033[%dG", 1)
	fmt.Print("\033[K")
}

func (t *Terminal) moveCursorToPos(pos int) {
	fmt.Printf("\033[%dG", pos+len(t.Prompt)+1)
}

func (t *Terminal) printPrompt() {
	fmt.Printf(gu.ColorizeForeground(t.Styles.PromptColor, "%v%v"), t.Styles.Cursor, t.Prompt)
}
