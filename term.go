package gocli

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	g "gocli/internal"

	"golang.org/x/term"
)

type Terminal struct {
	Prompt          string
	Options         []string
	HistoryId       string
	ExitMessage     string
	BypassCharacter string
}

func (t *Terminal) Get() (string, error) {

	if len(t.HistoryId) == 0 {
		t.HistoryId = "default"
	}
	history := g.GetCommandHistory(t.HistoryId)
	history.ResetIndex()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("error switching to raw mode: %v", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	var userInput string
	cursorPos := 0

	t.Prompt = fmt.Sprintf(g.Colorize(g.Cyan, "%v"), t.Prompt)
	fmt.Print(t.Prompt)

	for {
		buf := make([]byte, 3)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return "", fmt.Errorf("error reading input: %v", err)
		}

		input := buf[0]

		// Enter
		if len(userInput) > 0 && (input == 10 || input == 13) {
			bestMatch := g.BestMatch(userInput, t.Options)
			t.replaceLine(&cursorPos, &userInput, bestMatch)
			history.Append(userInput)
			t.cleanNextLineAndStay()
			if len(t.BypassCharacter) == 1 && strings.HasPrefix(userInput, t.BypassCharacter) {
				g.ExecCmd(userInput)
				return "", nil
			}
			return userInput, nil
		}

		// Exit CTRL+C
		if input == 3 {
			t.cleanNextLineAndStay()
			fmt.Println(g.Colorize(g.Green, t.ExitMessage))
			os.Exit(0)
		}

		// Autocomplete TAB
		if input == 9 {
			bestMatch := g.BestMatch(userInput, t.Options)
			if userInput == bestMatch {
				t.printAutocompleteSuggestions(userInput, cursorPos)
			} else {
				t.replaceLine(&cursorPos, &userInput, bestMatch)
			}
			continue
		}

		// Clean Screen CTRL+L
		if input == 12 {
			fmt.Print("\033[H\033[2J")
			fmt.Print(t.Prompt)
			continue
		}

		// Backspace
		if input == 127 {
			if cursorPos > 0 {
				userInput = userInput[:cursorPos-1] + userInput[cursorPos:]
				cursorPos--
				t.cleanLine()
				fmt.Print(userInput)
			}
		}

		// Arrows
		if input == 27 && len(buf) >= 3 && buf[1] == 91 {
			if buf[2] == 68 {
				if cursorPos > 0 {
					cursorPos--
					fmt.Print("\033[1D")
				}
			}
			if buf[2] == 67 {
				if cursorPos < len(userInput) {
					cursorPos++
					fmt.Print("\033[1C")
				}
			}
			if buf[2] == 65 {
				str, err := history.GetPrev(userInput)
				if err == nil {
					t.replaceLine(&cursorPos, &userInput, str)
				}
			}
			if buf[2] == 66 {
				str, err := history.GetNext()
				if err == nil {
					t.replaceLine(&cursorPos, &userInput, str)
				}
			}
		}

		if input >= 32 && input < 127 {
			userInput = userInput[:cursorPos] + string(input) + userInput[cursorPos:]
			cursorPos++
			fmt.Print(string(input))
		}

		t.cleanNextLine(cursorPos)
		t.moveCursorToPos(cursorPos)
	}

}

func (t *Terminal) replaceLine(pos *int, userInput *string, text string) {
	t.cleanLine()
	fmt.Print(text)
	t.moveCursorToPos(len(text))
	*pos = len(text)
	*userInput = text
}

func (t *Terminal) printAutocompleteSuggestions(userInput string, pos int) {
	t.cleanNextLineAndStay()
	fmt.Print(strings.Join(t.filterStrings(userInput), "  "))
	fmt.Print("\033[K")
	fmt.Print("\033[1A")
	t.moveCursorToPos(pos)
}

func (t *Terminal) filterStrings(prefix string) []string {
	var result []string
	for _, str := range t.Options {
		if strings.HasPrefix(str, prefix) && str != prefix {
			result = append(result, str)
		}
	}
	return result
}

func (t *Terminal) cleanLine() {
	t.moveCursorToPos(0)
	fmt.Print("\033[K")
}

func (t *Terminal) cleanNextLine(pos int) {
	t.cleanNextLineAndStay()
	fmt.Print("\033[1A")
	t.moveCursorToPos(pos)
}

func (t *Terminal) cleanNextLineAndStay() {
	fmt.Println()
	t.moveCursorToPosIgnorePrompt(1)
	fmt.Print("\033[K")
}

func (t *Terminal) moveCursorToPos(pos int) {
	fmt.Printf("\033[%dG", pos+len(t.Prompt)+1-9)
}

func (t *Terminal) moveCursorToPosIgnorePrompt(pos int) {
	fmt.Printf("\033[%dG", pos)
}

func (t *Terminal) PrintHistory(limit int) {
	g.GetCommandHistory(t.HistoryId).PrintHistory(limit)
}

func (t *Terminal) CountHistory() int {
	return g.GetCommandHistory(t.HistoryId).Count()
}

func (t *Terminal) ClearHistory() {
	g.GetCommandHistory(t.HistoryId).Clear()
}
