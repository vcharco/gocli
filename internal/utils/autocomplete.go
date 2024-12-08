package gocliutils

import (
	"strings"

	gt "github.com/vcharco/gocli/internal/types"
)

func BestMatch(userOption string, options []gt.Command) (string, bool) {
	if userOption == "" {
		return "", false
	}

	var filteredOptions []string
	for _, option := range options {
		if strings.HasPrefix(option.Name, userOption) {
			filteredOptions = append(filteredOptions, option.Name)
		}
	}

	if len(filteredOptions) == 1 {
		return filteredOptions[0], true
	}

	if len(filteredOptions) > 1 {
		return findCommonPrefix(filteredOptions), false
	}

	return userOption, false
}

func findCommonPrefix(stringsList []string) string {
	if len(stringsList) == 0 {
		return ""
	}

	strModel := stringsList[0]
	prefix := ""

	for {
		if len(prefix) >= len(strModel) {
			return prefix
		}

		newPrefix := strModel[:len(prefix)+1]

		match := true
		for _, str := range stringsList {
			if !strings.HasPrefix(str, newPrefix) {
				match = false
				break
			}
		}

		if !match {
			return prefix
		}

		prefix = newPrefix
	}
}
