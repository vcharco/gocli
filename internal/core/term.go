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
	Options         []gt.Candidate
	BypassCharacter string
	CursorPos       int
	CommandHistory  *CommandHistory
	CtrlKeys        []byte
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
	Options  map[string]string
	RawInput string
	Type     TerminalResponseType
	CtrlKey  byte
	Error    error
}

func (t *Terminal) Get(data ...string) TerminalResponse {

	if t.CommandHistory == nil {
		t.CommandHistory = &CommandHistory{Commands: []string{}, CurrentIndex: 0, Cache: "", IsCacheActive: false}
	}

	t.CommandHistory.ResetIndex()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return t.getTerminalResponse("", map[string]string{}, "", ExecutionError, 0, err, nil)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	var userInput string
	t.CursorPos = 0

	t.printPrompt()
	t.cleanNextLine()

	if len(data) > 0 && len(data[0]) > 0 {
		t.replaceLine(&userInput, data[0])
	}

	for {
		buf := make([]byte, 6)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return t.getTerminalResponse("", map[string]string{}, "", ExecutionError, 0, err, oldState)
		}

		input := buf[0]

		// Enter
		if len(userInput) > 0 && (input == 10 || input == 13) {

			// Bypass command to OS
			if len(t.BypassCharacter) > 0 && strings.HasPrefix(userInput, t.BypassCharacter) {
				t.CommandHistory.Append(userInput)
				rt := t.getTerminalResponse("", map[string]string{}, userInput[len(t.BypassCharacter):], OsCmd, 0, nil, oldState)
				gu.ExecCmd(userInput[len(t.BypassCharacter):])
				return rt
			}

			// Print help
			userInput = strings.Trim(userInput, " ")

			if strings.HasSuffix(userInput, "?") {
				userInput = userInput[:len(userInput)-1]
				command, err := GetClosestCommand(t.Options, userInput)
				if err != nil {
					return t.getTerminalResponse(command.Name, map[string]string{}, userInput, CmdError, 0, err, oldState)
				}
				tr := t.getTerminalResponse(command.Name, map[string]string{}, userInput, CmdHelp, 0, nil, oldState)
				t.printHelp(command)
				return tr
			}

			// Validate command
			command, params, err := ValidateCommand(t.Options, userInput)

			if err != nil {
				return t.getTerminalResponse("", map[string]string{}, userInput, ParamError, 0, err, oldState)
			}

			// Log in history, format line and return
			t.CommandHistory.Append(userInput)
			re := regexp.MustCompile(`^\S+`)
			t.replaceLine(&userInput, re.ReplaceAllString(userInput, command.Name))

			return t.getTerminalResponse(command.Name, params, userInput, Cmd, 0, nil, oldState)
		}

		// Check overriden CTRL+KEY
		ctrlKey := t.checkOverridenCtrl(input)
		if ctrlKey != 0 {
			return t.getTerminalResponse("", map[string]string{"key": string(ctrlKey)}, userInput, CtrlKey, ctrlKey, nil, oldState)
		}

		// Exit CTRL+C
		if input == 3 {
			term.Restore(int(os.Stdin.Fd()), oldState)
			fmt.Println()
			os.Exit(0)
		}

		// Autocomplete TAB
		if input == 9 {
			bestMatch, found := gu.BestMatch(userInput, t.Options)
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
			if t.CursorPos > 0 {
				userInput = userInput[:t.CursorPos-1] + userInput[t.CursorPos:]
				t.CursorPos--
				t.cleanLine()
				fmt.Print(userInput)
			}
		}

		// Arrows
		if input == 27 && len(buf) >= 3 && buf[1] == 91 {
			// LEFT
			if buf[2] == 68 {
				if t.CursorPos > 0 {
					t.CursorPos--
					fmt.Print("\033[1D")
				}
			}
			// RIGHT
			if buf[2] == 67 {
				if t.CursorPos < len(userInput) {
					t.CursorPos++
					fmt.Print("\033[1C")
				}
			}
			// UP
			if buf[2] == 65 {
				str, err := t.CommandHistory.GetPrev(userInput)
				if err == nil {
					t.replaceLine(&userInput, str)
				}
			}
			// DOWN
			if buf[2] == 66 {
				str, err := t.CommandHistory.GetNext()
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
			userInput = userInput[:t.CursorPos] + string(input) + userInput[t.CursorPos:]
			t.CursorPos++
			if t.CursorPos < len(userInput) {
				t.cleanLine()
				fmt.Print(userInput)
			} else {
				fmt.Print(string(input))
			}
		}

		t.cleanNextLine()
		t.moveCursorToPos(t.CursorPos)
	}

}

func (t *Terminal) replaceLine(userInput *string, text string) {
	t.cleanLine()
	fmt.Print(text)
	t.moveCursorToPos(len(text))
	t.CursorPos = len(text)
	*userInput = text
}

func (t *Terminal) printAutocompleteSuggestions(userInput string) {
	t.cleanNextLineAndStay()
	fmt.Print(strings.Join(t.filterStrings(userInput), "  "))
	fmt.Print("\033[K")
	fmt.Print("\033[1A")
	t.moveCursorToPos(t.CursorPos)
}

func (t *Terminal) filterStrings(prefix string) []string {
	var result []string
	gt.SortCandidates(t.Options)
	for _, candidate := range t.Options {
		if strings.HasPrefix(candidate.Name, prefix) && candidate.Name != prefix && !candidate.Hidden {
			result = append(result, candidate.Name)
		}
	}
	return result
}

func (t *Terminal) printHelp(command gt.Candidate) {

	gt.SortCandidateOptions(command.Options)

	if len(command.Description) > 0 {
		fmt.Printf("\n%v\n", command.Description)
	}

	var commandFlags []gt.CandidateOption
	var commandParams []gt.CandidateOption
	var defaultParam gt.CandidateOption

	for _, option := range command.Options {
		if option.Modifier&gt.DEFAULT != 0 {
			defaultParam = option
			continue
		}
		if option.Type == gt.None {
			commandFlags = append(commandFlags, option)
			continue
		}
		commandParams = append(commandParams, option)
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

	for _, option := range commandFlags {
		reqText := "OPTIONAL"
		if option.Modifier&gt.REQUIRED != 0 {
			reqText = "REQUIRED"
		}
		fmt.Printf("  %v: (%v) %v\n", option.Name, reqText, option.Description)
	}

	if len(commandParams) > 0 {
		fmt.Printf("\nPARAMS:\n")
	}

	for _, option := range commandParams {
		reqText := "OPTIONAL"
		if option.Modifier&gt.REQUIRED != 0 {
			reqText = "REQUIRED"
		}
		fmt.Printf("  %v <%v>: (%v) %v\n", option.Name, GetValidationTypeName(option.Type), reqText, option.Description)
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
	t.moveCursorToPos(t.CursorPos)
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

func (t *Terminal) getTerminalResponse(command string, options map[string]string, rawInput string, responseType TerminalResponseType, ctrlKey byte, err error, oldState *term.State) TerminalResponse {
	if oldState != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}
	fmt.Println()
	return TerminalResponse{Command: command, Options: options, RawInput: rawInput, Type: responseType, CtrlKey: ctrlKey, Error: err}
}

func (t *Terminal) printPrompt() {
	fmt.Printf(gu.Colorize(t.PromptColor, "%v"), t.Prompt)
}
