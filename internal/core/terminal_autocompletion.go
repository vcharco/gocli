package gocli

import (
	"fmt"
	"os"
	"strings"

	gt "github.com/vcharco/gocli/internal/types"
	"golang.org/x/term"
)

func (t *Terminal) printAutocompleteSuggestions(userInput string) {
	t.CleanNextLines(t.autoCompletionLines)
	t.cleanNextLineAndStay()
	adjusted, lines := t.GetAdjustedLine(t.filterCommands(userInput), "    ")
	if lines < 1 {
		lines = 1
	}
	t.autoCompletionLines = lines
	fmt.Printf("%v%v", t.Styles.ForegroundSuggestions, adjusted)
	for i := 0; i < lines; i++ {
		fmt.Print("\033[1A")
	}
	t.moveCursorToPos(t.cursorPos)
}

func (t *Terminal) filterCommands(prefix string) []string {
	var result []string
	gt.SortCommands(t.Commands)
	for _, candidate := range t.Commands {
		if strings.HasPrefix(candidate.Name, prefix) && candidate.Name != prefix && !candidate.Hidden {
			result = append(result, candidate.Name)
		}
	}
	return result
}

func (t *Terminal) GetAdjustedLine(items []string, separator string) (string, int) {
	maxLen, _, _ := term.GetSize(int(os.Stdout.Fd()))

	if maxLen <= 0 {
		return "", 0
	}

	adjusted := ""
	currentLine := ""
	lineCount := 1

	for _, item := range items {
		if len(item)+len(separator) > maxLen {
			return "", 0
		}

		if len(currentLine)+len(separator)+len(item) >= maxLen {
			adjusted += currentLine + "\n\033[G"
			lineCount++
			currentLine = item
		} else {
			if len(currentLine) != 0 {
				currentLine += separator + item
			} else {
				currentLine += item
			}
		}
	}

	adjusted += currentLine

	return adjusted, lineCount
}
