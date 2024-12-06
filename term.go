package gocli

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	g "github.com/vcharco/gocli/internal"

	"golang.org/x/term"
)

type Terminal struct {
	Prompt                string
	Options               []string
	ExitMessage           string
	InvalidCommandMessage string
	BypassCharacter       string
	AllowInvalidCommands  bool
	cursorPos             int
	commandHistory        *g.CommandHistory
}

func (t *Terminal) Get() (string, error) {

	if t.commandHistory == nil {
		t.commandHistory = &g.CommandHistory{Commands: []string{}, CurrentIndex: 0, Cache: "", IsCacheActive: false}
	}

	t.commandHistory.ResetIndex()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("error switching to raw mode: %v", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	var userInput string
	t.cursorPos = 0

	fmt.Printf(g.Colorize(g.Cyan, "%v"), t.Prompt)
	t.cleanNextLine()

	for {
		buf := make([]byte, 3)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return "", fmt.Errorf("error reading input: %v", err)
		}

		input := buf[0]

		// Enter
		if len(userInput) > 0 && (input == 10 || input == 13) {
			if len(t.BypassCharacter) == 1 && strings.HasPrefix(userInput, t.BypassCharacter) {
				term.Restore(int(os.Stdin.Fd()), oldState)
				t.commandHistory.Append(userInput)
				g.ExecCmd(userInput[1:])
				return "", nil
			}

			bestMatch := g.BestMatch(userInput, t.Options)
			t.commandHistory.Append(bestMatch)

			if !t.AllowInvalidCommands && !t.IsCommandExists(bestMatch) {
				term.Restore(int(os.Stdin.Fd()), oldState)
				invalidCommandMessage := t.InvalidCommandMessage
				if len(invalidCommandMessage) == 0 {
					invalidCommandMessage = "Invalid command"
				}
				fmt.Println()
				fmt.Println(g.Colorize(g.Red, t.InvalidCommandMessage))
				return "", nil
			}

			t.replaceLine(&userInput, bestMatch)
			t.cleanNextLineAndStay()

			return userInput, nil
		}

		// Exit CTRL+C
		if input == 3 {
			term.Restore(int(os.Stdin.Fd()), oldState)
			fmt.Println()
			if len(t.ExitMessage) > 0 {
				fmt.Println()
				fmt.Print(g.Colorize(g.Green, t.ExitMessage))
				fmt.Print("\n\n")
			}
			os.Exit(0)
		}

		// Autocomplete TAB
		if input == 9 {
			bestMatch := g.BestMatch(userInput, t.Options)
			if userInput == bestMatch {
				t.printAutocompleteSuggestions(userInput)
			} else {
				t.replaceLine(&userInput, bestMatch)
			}
			continue
		}

		// Clean Screen CTRL+L
		if input == 12 {
			fmt.Print("\033[H\033[2J")
			fmt.Printf(g.Colorize(g.Cyan, "%v"), t.Prompt)
			continue
		}

		// Backspace
		if input == 127 {
			if t.cursorPos > 0 {
				userInput = userInput[:t.cursorPos-1] + userInput[t.cursorPos:]
				t.cursorPos--
				t.cleanLine()
				fmt.Print(userInput)
			}
		}

		// Arrows
		if input == 27 && len(buf) >= 3 && buf[1] == 91 {
			if buf[2] == 68 {
				if t.cursorPos > 0 {
					t.cursorPos--
					fmt.Print("\033[1D")
				}
			}
			if buf[2] == 67 {
				if t.cursorPos < len(userInput) {
					t.cursorPos++
					fmt.Print("\033[1C")
				}
			}
			if buf[2] == 65 {
				str, err := t.commandHistory.GetPrev(userInput)
				if err == nil {
					t.replaceLine(&userInput, str)
				}
			}
			if buf[2] == 66 {
				str, err := t.commandHistory.GetNext()
				if err == nil {
					t.replaceLine(&userInput, str)
				}
			}
		}

		if input >= 32 && input < 127 {
			userInput = userInput[:t.cursorPos] + string(input) + userInput[t.cursorPos:]
			t.cursorPos++
			fmt.Print(string(input))
		}

		t.cleanNextLine()
		t.moveCursorToPos(t.cursorPos)
	}

}

func (t *Terminal) replaceLine(userInput *string, text string) {
	t.cleanLine()
	fmt.Print(text)
	t.moveCursorToPos(len(text))
	t.cursorPos = len(text)
	*userInput = text
}

func (t *Terminal) printAutocompleteSuggestions(userInput string) {
	t.cleanNextLineAndStay()
	fmt.Print(strings.Join(t.filterStrings(userInput), "  "))
	fmt.Print("\033[K")
	fmt.Print("\033[1A")
	t.moveCursorToPos(t.cursorPos)
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

func (t *Terminal) cleanNextLine() {
	t.cleanNextLineAndStay()
	fmt.Print("\033[1A")
	t.moveCursorToPos(t.cursorPos)
}

func (t *Terminal) cleanNextLineAndStay() {
	fmt.Println()
	t.moveCursorToPosIgnorePrompt(1)
	fmt.Print("\033[K")
}

func (t *Terminal) moveCursorToPos(pos int) {
	fmt.Printf("\033[%dG", pos+len(t.Prompt)+1)
}

func (t *Terminal) moveCursorToPosIgnorePrompt(pos int) {
	fmt.Printf("\033[%dG", pos)
}

func (t *Terminal) PrintHistory(limit int) {
	t.commandHistory.PrintHistory(limit)
}

func (t *Terminal) CountHistory() int {
	return t.commandHistory.Count()
}

func (t *Terminal) ClearHistory() {
	t.commandHistory.Clear()
}

func (t *Terminal) IsCommandExists(str string) bool {
	for _, v := range t.Options {
		if v == str || strings.HasPrefix(str, v+" ") {
			return true
		}
	}
	return false
}
