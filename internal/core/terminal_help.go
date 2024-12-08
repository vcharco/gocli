package gocli

import (
	"fmt"
	"strings"

	gv "github.com/vcharco/gocli/internal/core/validation"
	gt "github.com/vcharco/gocli/internal/types"
)

func (t *Terminal) printAutocompleteSuggestions(userInput string) {
	t.cleanNextLineAndStay()
	fmt.Print(t.Styles.ForegroundSuggestions)
	fmt.Print(strings.Join(t.filterCommands(userInput), "  "))
	fmt.Print("\033[K")
	fmt.Print("\033[1A")
	t.moveCursorToPos(t.cursorPos)
}

func (t *Terminal) filterCommands(prefix string) []string {
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
			usageLine += fmt.Sprintf(" [<%v>]", gv.GetValidationTypeName(defaultParam.Type))
		} else {
			usageLine += fmt.Sprintf(" <%v>", gv.GetValidationTypeName(defaultParam.Type))
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
		fmt.Printf("  %v <%v>: (%v) %v\n", param.Name, gv.GetValidationTypeName(param.Type), reqText, param.Description)
	}

	fmt.Println()
}
