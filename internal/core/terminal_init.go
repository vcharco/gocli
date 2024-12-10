package gocli

import (
	gu "github.com/vcharco/gocli/internal/utils"
)

func (t *Terminal) init() {
	t.cursorPos = 0
	t.startSelection = -1
	t.autoCompletionLines = 1
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
	if len(t.Styles.HelpTextForeground) == 0 {
		t.Styles.HelpTextForeground = gu.LightGray
	}
	if len(t.Styles.HelpTitlesForeground) == 0 {
		t.Styles.HelpTitlesForeground = gu.Blue
	}
	if len(t.Styles.HelpRequiredForeground) == 0 {
		t.Styles.HelpRequiredForeground = gu.Red
	}
	if len(t.Styles.HelpCommandForeground) == 0 {
		t.Styles.HelpCommandForeground = gu.White
	}
	if len(t.Styles.HelpParamsForeground) == 0 {
		t.Styles.HelpParamsForeground = gu.Yellow
	}
	if len(t.Styles.HelpLineColor) == 0 {
		t.Styles.HelpLineColor = gu.Blue
	}
	if t.commandHistory == nil {
		t.commandHistory = &commandHistory{Commands: []string{}, CurrentIndex: 0, Cache: "", IsCacheActive: false}
	}
	if t.Styles.Cursor == "" {
		t.Styles.Cursor = gu.CursorBlock
	}
	t.commandHistory.resetIndex()

	t.printPrompt()
	t.CleanNextLines(1)
}
