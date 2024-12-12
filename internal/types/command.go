package goclitypes

import "sort"

type Command struct {
	Name              string
	Description       string
	Hidden            bool
	Params            []Param
}

func SortCommands(candidates []Command) {
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Name < candidates[j].Name
	})
}

func GetCommandNames(commands []Command) []string {
	cmdNames := make([]string, 0, len(commands))
	for _, cmd := range commands {
		cmdNames = append(cmdNames, cmd.Name)
	}
	return cmdNames
}
