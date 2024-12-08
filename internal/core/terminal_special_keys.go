package gocli

import (
	"os"

	gt "github.com/vcharco/gocli/internal/types"
	gu "github.com/vcharco/gocli/internal/utils"
	"golang.org/x/term"
)

func (t *Terminal) checkSpecialKeys(input byte, userInput *string, oldState *term.State) byte {
	resp := t.checkOverridenCtrl(input)

	if resp == 0 {
		switch input {
		// Exit cli
		case gt.Ctrl_X:
			term.Restore(int(os.Stdin.Fd()), oldState)
			t.FnExitProgram()
		// Copy
		case gt.Ctrl_C:
			if t.startSelection < 1 {
				t.CopyToClipboard(*userInput)
			} else if t.startSelection < t.cursorPos {
				t.CopyToClipboard((*userInput)[t.startSelection:t.cursorPos])
			} else if t.cursorPos < t.startSelection {
				t.CopyToClipboard((*userInput)[t.cursorPos:t.startSelection])
			}
		// Paste
		case gt.Ctrl_V:
			t.PasteClipboard(userInput)
		// Clear screen
		case gt.Ctrl_L:
			t.FnClearScreen()
			t.printPrompt()
		// Move cursor at the beginning of the line
		case gt.Ctrl_A:
			t.cursorPos = 0
		// Move cursor at the end of the line
		case gt.Ctrl_E:
			t.cursorPos = len(*userInput)
		}
	}

	return resp
}

func (t *Terminal) checkOverridenCtrl(input byte) byte {
	for _, b := range t.CtrlKeys {
		if input == b {
			return b
		}
	}
	return 0
}

func (t *Terminal) CopyToClipboard(userInput string) {
	gu.SetClipboard(userInput)
}

func (t *Terminal) PasteClipboard(userInput *string) {
	if clipboard, err := gu.GetClipboardContent(); err == nil {
		t.replaceLine(userInput, *userInput+clipboard)
	}
}
