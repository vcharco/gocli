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
	gv "github.com/vcharco/gocli/internal/validation"
	"golang.org/x/term"
)

type Terminal struct {
	Styles              TerminalStyles
	Commands            []gt.Command
	BypassCharacter     string
	CtrlKeys            []byte
	cursorPos           int
	startSelection      int
	commandHistory      *commandHistory
	autoCompletionLines int
}

type TerminalStyles struct {
	Prompt                 string
	PromptColor            gu.Color
	ForegroundColor        gu.Color
	ForegroundSuggestions  gu.Color
	BackgroundColor        gu.BgColor
	SelForegroundColor     gu.Color
	SelBackgroundColor     gu.BgColor
	HelpTextForeground     gu.Color
	HelpTitlesForeground   gu.Color
	HelpCommandForeground  gu.Color
	HelpParamsForeground   gu.Color
	HelpRequiredForeground gu.Color
	HelpLineColor          gu.Color
	Cursor                 gu.Cursor
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

			// Log command in the history
			t.commandHistory.append(userInput)

			if err != nil {
				return getTerminalResponse("", map[string]interface{}{}, userInput, ParamError, 0, err, oldState)
			}

			// Format line and return
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

		// Clean current input line and all allocated by the suggestions
		t.CleanNextLines(t.autoCompletionLines)
		t.autoCompletionLines = 1
		t.CleanCurrentLine()
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
