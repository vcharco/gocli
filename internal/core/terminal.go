package gocli

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	gv "github.com/vcharco/gocli/internal/core/validation"
	gt "github.com/vcharco/gocli/internal/types"
	gu "github.com/vcharco/gocli/internal/utils"
	"golang.org/x/term"
)

type Terminal struct {
	Prompt          string
	PromptColor     gu.Color
	Commands        []gt.Command
	BypassCharacter string
	CtrlKeys        []byte
	cursorPos       int
	commandHistory  *commandHistory
}

func (t *Terminal) Get(data ...string) TerminalResponse {

	if t.commandHistory == nil {
		t.commandHistory = &commandHistory{Commands: []string{}, CurrentIndex: 0, Cache: "", IsCacheActive: false}
	}

	t.commandHistory.resetIndex()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return getTerminalResponse("", map[string]interface{}{}, "", ExecutionError, 0, err, nil)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	var userInput string
	t.cursorPos = 0

	t.printPrompt()
	t.cleanNextLine()

	// Append the incoming data to the userInput
	if len(data) > 0 {
		joinedData := strings.Join(data, " ")
		if len(joinedData) > 0 {
			t.replaceLine(&userInput, data[0])
		}
	}

	for {
		buf := make([]byte, 6)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return getTerminalResponse("", map[string]interface{}{}, "", ExecutionError, 0, err, oldState)
		}

		input := buf[0]

		// Enter
		if len(userInput) > 0 && (input == 10 || input == 13) {

			// Bypass command to OS
			if len(t.BypassCharacter) > 0 && strings.HasPrefix(userInput, t.BypassCharacter) {
				t.commandHistory.append(userInput)
				rt := getTerminalResponse("", map[string]interface{}{}, userInput[len(t.BypassCharacter):], OsCmd, 0, nil, oldState)
				gu.ExecCmd(userInput[len(t.BypassCharacter):])
				return rt
			}

			// Print help
			userInput = strings.Trim(userInput, " ")

			if strings.HasSuffix(userInput, "?") {
				userInput = userInput[:len(userInput)-1]
				command, err := gv.GetClosestCommand(t.Commands, userInput)
				if err != nil {
					return getTerminalResponse(command.Name, map[string]interface{}{}, userInput, CmdError, 0, err, oldState)
				}
				tr := getTerminalResponse(command.Name, map[string]interface{}{}, userInput, CmdHelp, 0, nil, oldState)
				t.printHelp(command)
				return tr
			}

			// Validate command
			command, params, err := gv.ValidateCommand(t.Commands, userInput)

			if err != nil {
				return getTerminalResponse("", map[string]interface{}{}, userInput, ParamError, 0, err, oldState)
			}

			// Log in history, format line and return
			t.commandHistory.append(userInput)
			re := regexp.MustCompile(`^\S+`)
			t.replaceLine(&userInput, re.ReplaceAllString(userInput, command.Name))

			return getTerminalResponse(command.Name, params, userInput, Cmd, 0, nil, oldState)
		}

		// Check overriden CTRL+KEY
		ctrlKey := t.checkSpecialKeys(input, oldState)
		if ctrlKey != 0 {
			return getTerminalResponse("", map[string]interface{}{"key": string(ctrlKey)}, userInput, CtrlKey, ctrlKey, nil, oldState)
		}

		// Autocomplete TAB
		if input == 9 {
			bestMatch, found := gu.BestMatch(userInput, gt.GetCommandNames(t.Commands))
			if userInput == bestMatch {
				t.printAutocompleteSuggestions(userInput)
			} else {
				if found {
					t.replaceLine(&userInput, bestMatch+" ")
				} else {
					t.replaceLine(&userInput, bestMatch)
				}
			}
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
			// LEFT
			if buf[2] == 68 {
				if t.cursorPos > 0 {
					t.cursorPos--
					fmt.Print("\033[1D")
				}
			}
			// RIGHT
			if buf[2] == 67 {
				if t.cursorPos < len(userInput) {
					t.cursorPos++
					fmt.Print("\033[1C")
				}
			}
			// UP
			if buf[2] == 65 {
				str, err := t.commandHistory.getPrev(userInput)
				if err == nil {
					t.replaceLine(&userInput, str)
				}
			}
			// DOWN
			if buf[2] == 66 {
				str, err := t.commandHistory.getNext()
				if err == nil {
					t.replaceLine(&userInput, str)
				}
			}

			if len(buf) >= 6 {
				// SHIFT + ARROWS
				if buf[2] == 49 && buf[3] == 59 && buf[4] == 50 {
					// SHIFT + LEFT
					if buf[5] == 68 {
						continue
					}
					// SHIFT + RIGHT
					if buf[5] == 67 {
						continue
					}
					// SHIFT + UP
					if buf[5] == 65 {
						continue
					}
					// SHIFT + DOWN
					if buf[5] == 66 {
						continue
					}
				}

				// ALT (OPTION) + ARROWS
				if buf[2] == 49 && buf[3] == 59 && buf[4] == 51 {
					// ALT + LEFT | RIGHT | UP | DOWN
					if buf[5] == 68 || buf[5] == 67 || buf[5] == 65 || buf[5] == 66 {
						continue
					}
				}
			}

		}

		if input >= 32 && input < 127 {
			userInput = userInput[:t.cursorPos] + string(input) + userInput[t.cursorPos:]
			t.cursorPos++
			if t.cursorPos < len(userInput) {
				t.cleanLine()
				fmt.Print(userInput)
			} else {
				fmt.Print(string(input))
			}
		}

		t.cleanNextLine()
		t.moveCursorToPos(t.cursorPos)
	}

}
