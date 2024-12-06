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

type Candidate struct {
	Name              string
	DefaultOptionType CandidateType
	Options           []CandidateOption
}

type CandidateOption struct {
	Name string
	Type CandidateType
}
