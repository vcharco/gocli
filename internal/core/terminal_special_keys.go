package gocli

import (
	"os"

	"golang.org/x/term"
)

func (t *Terminal) checkSpecialKeys(input byte, oldState *term.State) byte {
	resp := t.checkOverridenCtrl(input)

	if resp == 0 {
		switch input {
		case 3:
			term.Restore(int(os.Stdin.Fd()), oldState)
			t.FnExitProgram()
		case 12:
			t.FnClearScreen()
			t.printPrompt()
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
