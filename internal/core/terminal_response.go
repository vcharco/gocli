package gocli

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

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

func (tr *TerminalResponse) GetParam(name string, defaultValue interface{}) interface{} {
	if value, exists := tr.Params[name]; exists {
		if val, ok := value.(string); ok && len(val) == 0 {
			return true
		}
		return value
	}
	return defaultValue
}

func getTerminalResponse(command string, params map[string]interface{}, rawInput string, responseType TerminalResponseType, ctrlKey byte, err error, oldState *term.State) TerminalResponse {
	if oldState != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}
	fmt.Println()
	return TerminalResponse{Command: command, Params: params, RawInput: rawInput, Type: responseType, CtrlKey: ctrlKey, Error: err}
}
