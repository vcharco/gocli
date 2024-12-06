package gocli

import (
	gg "github.com/vcharco/gocli/internal/core"
	gt "github.com/vcharco/gocli/internal/types"
)

type Terminal = gg.Terminal
type TerminalResponseType = gg.TerminalResponseType
type Candidate = gt.Candidate
type CandidateOption = gt.CandidateOption
type CandidateType = gt.CandidateType

const (
	Date   = gt.Date
	Domain = gt.Domain
	Email  = gt.Email
	Ipv4   = gt.Ipv4
	Ipv6   = gt.Ipv6
	Number = gt.Number
	Phone  = gt.Phone
	Text   = gt.Text
	Time   = gt.Time
	Url    = gt.Url
	UUID   = gt.UUID
)

const (
	Cmd            = gg.Cmd
	OsCmd          = gg.OsCmd
	CmdError       = gg.CmdError
	ParamError     = gg.ParamError
	ExecutionError = gg.ExecutionError
)
