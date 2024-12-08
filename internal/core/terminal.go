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
	Styles          TerminalStyles
	Commands        []gt.Command
	BypassCharacter string
	CtrlKeys        []byte
	cursorPos       int
	startSelection  int
	commandHistory  *commandHistory
}

type TerminalStyles struct {
	Prompt                string
	PromptColor           gu.Color
	ForegroundColor       gu.Color
	ForegroundSuggestions gu.Color
	BackgroundColor       gu.BgColor
	SelForegroundColor    gu.Color
	SelBackgroundColor    gu.BgColor
	Cursor                gu.Cursor
}

func (t *Terminal) Get(data ...string) TerminalResponse {

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return getTerminalResponse("", map[string]interface{}{}, "", ExecutionError, 0, err, nil)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	t.init()

	var userInput string
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

		// Check special commands and overriden CTRL+KEY
		ctrlKey := t.checkSpecialKeys(input, &userInput, oldState)
		if ctrlKey != 0 {
			return getTerminalResponse("", nil, userInput, CtrlKey, ctrlKey, nil, oldState)
		}

		t.checkTextSelection(input, buf, &userInput)

		// Autocomplete TAB
		if input == 9 {
			bestMatch, _ := gu.BestMatch(userInput, gt.GetCommandNames(t.Commands))
			if userInput == bestMatch {
				t.printAutocompleteSuggestions(userInput)
				continue
			} else {
				userInput = bestMatch
				t.cursorPos = len(bestMatch)
			}
		}

		// Backspace
		if input == 127 {
			if t.cursorPos > 0 {
				userInput = userInput[:t.cursorPos-1] + userInput[t.cursorPos:]
				t.cursorPos--
			}
		}

		// Handle cursor movement and text selection
		if !t.handleCursorAndContinue(input, buf, &userInput) {
			continue
		}

		// Print characters
		if input >= 32 && input < 127 {
			userInput = userInput[:t.cursorPos] + string(input) + userInput[t.cursorPos:]
			t.cursorPos++
		}

		// Clean current input and write new one
		t.cleanNextLine()
		t.cleanLine()
		output := fmt.Sprint(gu.ColorizeBoth(t.Styles.ForegroundColor, t.Styles.BackgroundColor, userInput))

		// Apply highlight to selected text
		if highlighted, ok := t.highlightSelected(userInput); ok {
			output = highlighted
		}

		// Print the line
		fmt.Print(output)

		// Set the cursor position at the right place
		t.moveCursorToPos(t.cursorPos)
	}

}

func (t *Terminal) init() {
	t.cursorPos = 0
	t.startSelection = -1
	if len(t.Styles.Prompt) == 0 {
		t.Styles.Prompt = "gocli> "
	}
	if len(t.Styles.PromptColor) == 0 {
		t.Styles.PromptColor = gu.Blue
	}
	if len(t.Styles.ForegroundColor) == 0 {
		t.Styles.ForegroundColor = gu.White
	}
	if len(t.Styles.ForegroundSuggestions) == 0 {
		t.Styles.ForegroundSuggestions = gu.LightGray
	}
	if len(t.Styles.BackgroundColor) == 0 {
		t.Styles.BackgroundColor = gu.BgTransparent
	}
	if len(t.Styles.SelBackgroundColor) == 0 {
		t.Styles.SelBackgroundColor = gu.BgLightBlue
	}
	if len(t.Styles.SelForegroundColor) == 0 {
		t.Styles.SelForegroundColor = gu.Black
	}
	if t.commandHistory == nil {
		t.commandHistory = &commandHistory{Commands: []string{}, CurrentIndex: 0, Cache: "", IsCacheActive: false}
	}
	if t.Styles.Cursor == "" {
		t.Styles.Cursor = gu.CursorBlock
	}
	t.commandHistory.resetIndex()

	t.printPrompt()
	t.cleanNextLine()
}
