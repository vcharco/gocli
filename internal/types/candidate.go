package goclitypes

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
	Options     []CandidateOption
}

type CandidateOption struct {
	Name        string
	Description string
	Modifier    CandidateOptionModifier
	Type        CandidateType
}
