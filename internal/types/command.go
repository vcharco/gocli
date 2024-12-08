package goclitypes

import "sort"

type Command struct {
	Name        string
	Description string
	Hidden      bool
	Params      []Param
}

func SortCommands(candidates []Command) {
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Name < candidates[j].Name
	})
}
