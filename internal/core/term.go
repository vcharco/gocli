package gocli

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

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

type TerminalResponseType int

const (
	Cmd TerminalResponseType = iota
	OsCmd
	CtrlKey
	CmdHelp
	CmdError
	ParamError
	ExecutionError
)

type TerminalResponse struct {
	Command  string
	Params   map[string]interface{}
	RawInput string
	Type     TerminalResponseType
	CtrlKey  byte
	Error    error
}

func (t *Terminal) Get(data ...string) TerminalResponse {

	if t.commandHistory == nil {
		t.commandHistory = &commandHistory{Commands: []string{}, CurrentIndex: 0, Cache: "", IsCacheActive: false}
	}

	t.commandHistory.resetIndex()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return t.getTerminalResponse("", map[string]interface{}{}, "", ExecutionError, 0, err, nil)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	var userInput string
	t.cursorPos = 0

	t.printPrompt()
	t.cleanNextLine()

	if len(data) > 0 && len(data[0]) > 0 {
		t.replaceLine(&userInput, data[0])
	}

	for {
		buf := make([]byte, 6)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return t.getTerminalResponse("", map[string]interface{}{}, "", ExecutionError, 0, err, oldState)
		}

		input := buf[0]

		// Enter
		if len(userInput) > 0 && (input == 10 || input == 13) {

			// Bypass command to OS
			if len(t.BypassCharacter) > 0 && strings.HasPrefix(userInput, t.BypassCharacter) {
				t.commandHistory.append(userInput)
				rt := t.getTerminalResponse("", map[string]interface{}{}, userInput[len(t.BypassCharacter):], OsCmd, 0, nil, oldState)
				gu.ExecCmd(userInput[len(t.BypassCharacter):])
				return rt
			}

			// Print help
			userInput = strings.Trim(userInput, " ")

			if strings.HasSuffix(userInput, "?") {
				userInput = userInput[:len(userInput)-1]
				command, err := GetClosestCommand(t.Commands, userInput)
				if err != nil {
					return t.getTerminalResponse(command.Name, map[string]interface{}{}, userInput, CmdError, 0, err, oldState)
				}
				tr := t.getTerminalResponse(command.Name, map[string]interface{}{}, userInput, CmdHelp, 0, nil, oldState)
				t.printHelp(command)
				return tr
			}

			// Validate command
			command, params, err := ValidateCommand(t.Commands, userInput)

			if err != nil {
				return t.getTerminalResponse("", map[string]interface{}{}, userInput, ParamError, 0, err, oldState)
			}

			// Log in history, format line and return
			t.commandHistory.append(userInput)
			re := regexp.MustCompile(`^\S+`)
			t.replaceLine(&userInput, re.ReplaceAllString(userInput, command.Name))

			return t.getTerminalResponse(command.Name, params, userInput, Cmd, 0, nil, oldState)
		}

		// Check overriden CTRL+KEY
		ctrlKey := t.checkOverridenCtrl(input)
		if ctrlKey != 0 {
			return t.getTerminalResponse("", map[string]interface{}{"key": string(ctrlKey)}, userInput, CtrlKey, ctrlKey, nil, oldState)
		}

		// Exit CTRL+C
		if input == 3 {
			term.Restore(int(os.Stdin.Fd()), oldState)
			fmt.Println()
			os.Exit(0)
		}

		// Autocomplete TAB
		if input == 9 {
			bestMatch, found := gu.BestMatch(userInput, t.Commands)
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

		// Clean Screen CTRL+L
		if input == 12 {
			fmt.Print("\033[H\033[2J")
			t.printPrompt()
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
					// ALT + LEFT
					if buf[5] == 68 {
						continue
					}
					// ALT + RIGHT
					if buf[5] == 67 {
						continue
					}
					// ALT + UP
					if buf[5] == 65 {
						continue
					}
					// ALT + DOWN
					if buf[5] == 66 {
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
	gt.SortCommands(t.Commands)
	for _, candidate := range t.Commands {
		if strings.HasPrefix(candidate.Name, prefix) && candidate.Name != prefix && !candidate.Hidden {
			result = append(result, candidate.Name)
		}
	}
	return result
}

func (t *Terminal) printHelp(command gt.Command) {

	gt.SortParams(command.Params)

	if len(command.Description) > 0 {
		fmt.Printf("\n%v\n", command.Description)
	}

	var commandFlags []gt.Param
	var commandParams []gt.Param
	var defaultParam gt.Param

	for _, param := range command.Params {
		if param.Modifier&gt.DEFAULT != 0 {
			defaultParam = param
			continue
		}
		if param.Type == gt.None {
			commandFlags = append(commandFlags, param)
			continue
		}
		commandParams = append(commandParams, param)
	}

	usageLine := fmt.Sprintf("\nUsage: %v", command.Name)

	if len(commandFlags) > 0 {
		usageLine += " [FLAGS]"
	}

	if len(commandParams) > 0 {
		usageLine += " [PARAMS]"
	}

	if len(defaultParam.Name) > 0 {
		if defaultParam.Modifier&gt.REQUIRED != 0 {
			usageLine += fmt.Sprintf(" [<%v>]", GetValidationTypeName(defaultParam.Type))
		} else {
			usageLine += fmt.Sprintf(" <%v>", GetValidationTypeName(defaultParam.Type))
		}
	}

	fmt.Println(usageLine)

	if len(defaultParam.Description) > 0 {
		fmt.Printf("\nDEFAULT PARAM: %v\n", defaultParam.Description)
	}

	if len(commandFlags) > 0 {
		fmt.Printf("\nFLAGS:\n")
	}

	for _, param := range commandFlags {
		reqText := "OPTIONAL"
		if param.Modifier&gt.REQUIRED != 0 {
			reqText = "REQUIRED"
		}
		fmt.Printf("  %v: (%v) %v\n", param.Name, reqText, param.Description)
	}

	if len(commandParams) > 0 {
		fmt.Printf("\nPARAMS:\n")
	}

	for _, param := range commandParams {
		reqText := "OPTIONAL"
		if param.Modifier&gt.REQUIRED != 0 {
			reqText = "REQUIRED"
		}
		fmt.Printf("  %v <%v>: (%v) %v\n", param.Name, GetValidationTypeName(param.Type), reqText, param.Description)
	}

	fmt.Println()
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

func (t *Terminal) checkOverridenCtrl(input byte) byte {
	for _, b := range t.CtrlKeys {
		if input == b {
			return b
		}
	}
	return 0
}

func (t *Terminal) getTerminalResponse(command string, params map[string]interface{}, rawInput string, responseType TerminalResponseType, ctrlKey byte, err error, oldState *term.State) TerminalResponse {
	if oldState != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}
	fmt.Println()
	return TerminalResponse{Command: command, Params: params, RawInput: rawInput, Type: responseType, CtrlKey: ctrlKey, Error: err}
}

func (t *Terminal) printPrompt() {
	fmt.Printf(gu.Colorize(t.PromptColor, "%v"), t.Prompt)
}

func (t *Terminal) PrintHistory(limit int) {
	t.commandHistory.print(limit)
}

func (t *Terminal) ClearHistory() {
	t.commandHistory.clear()
}

func (t *Terminal) CountHistory() int {
	return t.commandHistory.count()
}

func (t *Terminal) GetHistoryAt(index int) (string, error) {
	return t.commandHistory.getAt(index)
}

func (t *Terminal) GetHistory(index int) []string {
	return t.commandHistory.getAll()
}

func (tr *TerminalResponse) GetParam(name string, defaultValue interface{}) interface{} {
	if value, exists := tr.Params[name]; exists {
		if val, ok := value.(string); ok && len(val) == 0 {
			return true
		}
		return value
	}
	return defaultValue
}
