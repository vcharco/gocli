package goclitypes

import "sort"

type ParamType int

const (
	None ParamType = iota
	Date
	Domain
	Email
	Ipv4
	Ipv6
	Number
	FloatNumber
	Phone
	Text
	Time
	Url
	UUID
)

type ParamModifier int

const (
	DEFAULT ParamModifier = 1 << iota
	REQUIRED
)

type Param struct {
	Name        string
	Description string
	Modifier    ParamModifier
	Type        ParamType
}

func SortParams(candidates []Param) {
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Name < candidates[j].Name
	})
}
