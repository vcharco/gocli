package gocli

import (
	"fmt"
)

func (t *Terminal) replaceLine(userInput *string, text string) {
	t.CleanCurrentLine()
	fmt.Print(text)
	t.moveCursorToPos(len(text))
	t.cursorPos = len(text)
	*userInput = text
}

func (t *Terminal) CleanCurrentLine() {
	t.moveCursorToPos(0)
	fmt.Print("\033[K")
}

func (t *Terminal) cleanNextLineAndStay() {
	fmt.Println()
	fmt.Printf("\033[%dG", 1)
	fmt.Print("\033[K")
}

func (t *Terminal) CleanNextLines(lines int) {
	for i := 0; i < lines; i++ {
		fmt.Println()
		fmt.Printf("\033[%dG", 1)
		fmt.Print("\033[K")
	}
	for i := 0; i < lines; i++ {
		fmt.Print("\033[1A")
	}
	t.moveCursorToPos(t.cursorPos)
}

func (t *Terminal) moveCursorToPos(pos int) {
	fmt.Printf("\033[%dG", pos+len(t.Styles.Prompt)+1)
}

func (t *Terminal) printPrompt() {
	prompt := string(t.Styles.PromptColor) + t.Styles.Prompt // Prompt color
	prompt += string(t.Styles.Cursor)                        // Cursor type

	fmt.Print(prompt)
}
