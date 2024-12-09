package gocli

import (
	"fmt"
	"math"
	"strings"

	gv "github.com/vcharco/gocli/internal/core/validation"
	gt "github.com/vcharco/gocli/internal/types"
	gu "github.com/vcharco/gocli/internal/utils"
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

	var commandFlags []gt.Param
	var commandParams []gt.Param
	var defaultParam gt.Param
	largestParamNameLen := 0
	largestFlagNameLen := 0

	for _, param := range command.Params {
		if param.Modifier&gt.DEFAULT != 0 {
			defaultParam = param
			continue
		}
		if param.Type == gt.None {
			commandFlags = append(commandFlags, param)
			paramLen := len(param.Name)
			if largestFlagNameLen < paramLen {
				largestFlagNameLen = paramLen
			}
			continue
		}
		commandParams = append(commandParams, param)
		paramLen := len(param.Name) + len(gv.GetValidationTypeName(param.Type))
		if largestParamNameLen < paramLen {
			largestParamNameLen = paramLen
		}
	}

	vs := gu.ColorizeForeground(t.Styles.HelpLineColor, "│")
	hs := gu.ColorizeForeground(t.Styles.HelpLineColor, "─")
	tc := gu.ColorizeForeground(t.Styles.HelpLineColor, "┬")
	tl := gu.ColorizeForeground(t.Styles.HelpLineColor, "┌")
	tr := gu.ColorizeForeground(t.Styles.HelpLineColor, "┐")
	bl := gu.ColorizeForeground(t.Styles.HelpLineColor, "└")
	br := gu.ColorizeForeground(t.Styles.HelpLineColor, "┘")
	lc := gu.ColorizeForeground(t.Styles.HelpLineColor, "├")

	var prefix string

	prefix = strings.Repeat(hs, int(math.Max(0, float64(len(command.Name)))))
	fmt.Printf("\n%v%v%v%v%v\n%v %v %v\n%v%v%v%v%v\n", tl, hs, hs, prefix, tr, vs, gu.ColorizeForeground(t.Styles.HelpCommandForeground, command.Name), vs, bl, hs, tc, prefix, br)

	if len(command.Description) > 0 {
		if len(defaultParam.Name) > 0 || len(commandFlags) > 0 || len(commandParams) > 0 {
			prefix = fmt.Sprintf("  %v\n  %v%v ", vs, lc, hs)
		} else {
			prefix = fmt.Sprintf("  %v\n  %v%v ", vs, bl, hs)
		}
		fmt.Printf("%v%v%v\n", prefix, gu.ColorizeForeground(t.Styles.HelpTitlesForeground, "DESCRIPTION  "), gu.ColorizeForeground(t.Styles.HelpTextForeground, command.Description))
	}

	if len(defaultParam.Name) > 0 || len(commandFlags) > 0 || len(commandParams) > 0 {
		prefix = fmt.Sprintf("  %v\n  %v%v ", vs, lc, hs)
	} else {
		prefix = fmt.Sprintf("  %v\n  %v%v ", vs, bl, hs)
	}

	usageLine := fmt.Sprintf("%v%v%v", prefix, gu.ColorizeForeground(t.Styles.HelpTitlesForeground, "USAGE  "), command.Name)
	usageLineValue := ""

	if len(commandFlags) > 0 {
		usageLineValue += " [FLAGS]"
	}

	if len(commandParams) > 0 {
		usageLineValue += " [PARAMS]"
	}

	if len(defaultParam.Name) > 0 {
		if defaultParam.Modifier&gt.REQUIRED != 0 {
			usageLineValue += fmt.Sprintf(" <%v>", gv.GetValidationTypeName(defaultParam.Type))
		} else {
			usageLineValue += fmt.Sprintf(" [<%v>]", gv.GetValidationTypeName(defaultParam.Type))
		}
	}

	fmt.Printf("%v\n", usageLine+gu.ColorizeForeground(t.Styles.HelpTextForeground, usageLineValue))

	if len(commandFlags) > 0 || len(commandParams) > 0 {
		prefix = fmt.Sprintf("  %v\n  %v%v ", vs, lc, hs)
	} else {
		prefix = fmt.Sprintf("  %v\n  %v%v ", vs, lc, bl)
	}

	if len(defaultParam.Description) > 0 {
		fmt.Printf("%v%v%v\n", prefix, gu.ColorizeForeground(t.Styles.HelpTitlesForeground, "DEFAULT PARAM  "), gu.ColorizeForeground(t.Styles.HelpTextForeground, defaultParam.Description))
	}

	if len(commandParams) > 0 {
		prefix = fmt.Sprintf("  %v\n  %v%v ", vs, lc, hs)
	} else {
		prefix = fmt.Sprintf("  %v\n  %v%v ", vs, lc, bl)
	}

	if len(commandFlags) > 0 {
		fmt.Printf("%v%v\n", prefix, gu.ColorizeForeground(t.Styles.HelpTitlesForeground, "FLAGS"))
	}

	for i, param := range commandFlags {
		if i < len(commandFlags)-1 {
			if len(commandParams) > 0 {
				prefix = fmt.Sprintf("  %v   %v\n  %v   %v%v ", vs, vs, vs, lc, hs)
			} else {
				prefix = fmt.Sprintf("      %v\n      %v%v ", vs, bl, hs)
			}
		} else {
			if len(commandParams) > 0 {
				prefix = fmt.Sprintf("  %v   %v\n  %v   %v%v ", vs, vs, vs, bl, hs)
			} else {
				prefix = fmt.Sprintf("      %v\n      %v%v ", vs, bl, hs)
			}
		}
		reqText := ""
		if param.Modifier&gt.REQUIRED != 0 {
			reqText = gu.ColorizeForeground(t.Styles.HelpRequiredForeground, " (REQUIRED)")
		}
		formattedParamName := fmt.Sprintf("%-*v", largestFlagNameLen, param.Name)
		fmt.Printf("%v %v %v %v\n", prefix, gu.ColorizeForeground(t.Styles.HelpParamsForeground, formattedParamName), reqText, gu.ColorizeForeground(t.Styles.HelpTextForeground, param.Description))
	}

	prefix = fmt.Sprintf("  %v\n  %v%v ", vs, bl, hs)

	if len(commandParams) > 0 {
		fmt.Printf("%v%v\n", prefix, gu.ColorizeForeground(t.Styles.HelpTitlesForeground, "PARAMS"))
	}

	for i, param := range commandParams {
		if i < len(commandParams)-1 {
			prefix = fmt.Sprintf("      %v\n      %v%v ", vs, lc, hs)
		} else {
			prefix = fmt.Sprintf("      %v\n      %v%v ", vs, bl, hs)
		}
		reqText := ""
		if param.Modifier&gt.REQUIRED != 0 {
			reqText = gu.ColorizeForeground(t.Styles.HelpRequiredForeground, " (REQUIRED)")
		}
		paramValue := fmt.Sprintf("%v <%v>", param.Name, gv.GetValidationTypeName(param.Type))
		paramValue = fmt.Sprintf("%-*v", largestParamNameLen+3, paramValue)
		fmt.Printf("%v %v %v %v\n", prefix, gu.ColorizeForeground(t.Styles.HelpParamsForeground, paramValue), reqText, gu.ColorizeForeground(t.Styles.HelpTextForeground, param.Description))
	}

	fmt.Println()
}
