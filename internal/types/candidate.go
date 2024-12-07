package goclitypes

import "sort"

type CandidateType int

const (
	None CandidateType = iota
	Date
	Domain
	Email
	Ipv4
	Ipv6
	Number
	Phone
	Text
	Time
	Url
	UUID
)

type CandidateOptionModifier int

const (
	DEFAULT CandidateOptionModifier = 1 << iota
	REQUIRED
)

type Candidate struct {
	Name        string
	Description string
	Hidden      bool
	Options     []CandidateOption
}

type CandidateOption struct {
	Name        string
	Description string
	Modifier    CandidateOptionModifier
	Type        CandidateType
}

func SortCandidates(candidates []Candidate) {
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Name < candidates[j].Name
	})
}

func SortCandidateOptions(candidates []CandidateOption) {
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Name < candidates[j].Name
	})
}
