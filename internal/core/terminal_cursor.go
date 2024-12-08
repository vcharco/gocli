package gocli

import (
	"fmt"

	gu "github.com/vcharco/gocli/internal/utils"
)

func (t *Terminal) handleCursorAndContinue(input byte, buf []byte, userInput *string) bool {

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
			if t.cursorPos < len(*userInput) {
				t.cursorPos++
				fmt.Print("\033[1C")
			}
		}
		// UP
		if buf[2] == 65 {
			str, err := t.commandHistory.getPrev(*userInput)
			if err == nil {
				t.replaceLine(userInput, str)
			}
		}
		// DOWN
		if buf[2] == 66 {
			str, err := t.commandHistory.getNext()
			if err == nil {
				t.replaceLine(userInput, str)
			}
		}

		if len(buf) >= 6 {
			// SHIFT + ARROWS
			if buf[2] == 49 && buf[3] == 59 && buf[4] == 50 {
				// UP | DOWN (future use)
				if buf[5] == 65 || buf[5] == 66 {
					return false
				}
			}

			// ALT (OPTION) + ARROWS
			if buf[2] == 49 && buf[3] == 59 && buf[4] == 51 {
				// ALT + LEFT | RIGHT | UP | DOWN (future use)
				if buf[5] == 68 || buf[5] == 67 || buf[5] == 65 || buf[5] == 66 {
					return false
				}
			}
		}
	}

	return true
}

// This must be executed after Clipboard validation, else Clipboard Copy (CTRL+C) always be empty
func (t *Terminal) checkTextSelection(input byte, buf []byte, userInput *string) {
	if input == 27 && len(buf) >= 3 && buf[1] == 91 {
		if len(buf) >= 6 {
			// SHIFT + ARROWS
			if buf[2] == 49 && buf[3] == 59 && buf[4] == 50 {
				// SHIFT + LEFT
				if buf[5] == 68 {
					if t.cursorPos > 0 {
						// If startSelection is -1, the selection has just begun, we set it
						if t.startSelection == -1 {
							t.startSelection = t.cursorPos
						}
						t.cursorPos--
					}
					return
				}
				// SHIFT + RIGHT
				if buf[5] == 67 {
					if t.cursorPos < len(*userInput) {
						// If startSelection is -1, the selection has just begun, we set it
						if t.startSelection == -1 {
							t.startSelection = t.cursorPos
						}
						t.cursorPos++
					}
					return
				}
			}
		}
		// If SHIFT+L/R was pressed, we cancel selection
		t.startSelection = -1
	}
}

func (t *Terminal) highlightSelected(userInput string) (string, bool) {
	if t.startSelection != -1 {
		init := t.startSelection
		end := t.cursorPos
		if init > end {
			init, end = end, init
		}
		colorizedSelection := gu.ColorizeBoth(t.Styles.SelForegroundColor, t.Styles.SelBackgroundColor, userInput[init:end])
		regularTextStart := gu.ColorizeBoth(t.Styles.ForegroundColor, t.Styles.BackgroundColor, userInput[:init])
		regularTextEnd := gu.ColorizeBoth(t.Styles.ForegroundColor, t.Styles.BackgroundColor, userInput[end:])
		return fmt.Sprintf("%v%v%v", regularTextStart, colorizedSelection, regularTextEnd), true
	}
	return userInput, false
}
