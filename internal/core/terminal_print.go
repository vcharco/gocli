package gocli

import "fmt"

func (t *Terminal) PrintInfo(text string, params ...any) {
	t.PrintText(fmt.Sprintf("\033[36m%v\033[0m", text), params...)
}

func (t *Terminal) PrintError(text string, params ...any) {
	t.PrintText(fmt.Sprintf("\033[31m%v\033[0m", text), params...)
}

func (t *Terminal) PrintSuccess(text string, params ...any) {
	t.PrintText(fmt.Sprintf("\033[32m%v\033[0m", text), params...)
}

func (t *Terminal) PrintWarning(text string, params ...any) {
	t.PrintText(fmt.Sprintf("\033[33m%v\033[0m", text), params...)
}

func (t *Terminal) PrintText(text string, params ...any) {
	if len(params) == 0 {
		fmt.Print(text)
	} else {
		fmt.Printf(text, params...)
	}
	fmt.Println()
}
