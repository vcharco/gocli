package gocli

import (
	gg "github.com/vcharco/gocli/internal/core"
	gt "github.com/vcharco/gocli/internal/types"
	gu "github.com/vcharco/gocli/internal/utils"
)

type Terminal = gg.Terminal
type TerminalResponse = gg.TerminalResponse
type TerminalResponseType = gg.TerminalResponseType
type Candidate = gt.Candidate
type CandidateOption = gt.CandidateOption
type CandidateOptionModifier = gt.CandidateOptionModifier
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
	CtrlKey        = gg.CtrlKey
	ParamError     = gg.ParamError
	ExecutionError = gg.ExecutionError
)

const (
	Red     = gu.Red
	Green   = gu.Green
	Yellow  = gu.Yellow
	Blue    = gu.Blue
	Magenta = gu.Magenta
	Cyan    = gu.Cyan
	White   = gu.White
	Reset   = gu.Reset
)

const (
	DEFAULT  = gt.DEFAULT
	REQUIRED = gt.REQUIRED
)

const (
	Ctrl_A byte = 1
	Ctrl_B byte = 2
	Ctrl_C byte = 3
	Ctrl_D byte = 4
	Ctrl_E byte = 5
	Ctrl_F byte = 6
	Ctrl_G byte = 7
	Ctrl_H byte = 8
	Ctrl_I byte = 9
	Ctrl_J byte = 10
	Ctrl_K byte = 11
	Ctrl_L byte = 12
	Ctrl_M byte = 13
	Ctrl_N byte = 14
	Ctrl_O byte = 15
	Ctrl_P byte = 16
	Ctrl_Q byte = 17
	Ctrl_R byte = 18
	Ctrl_S byte = 19
	Ctrl_T byte = 20
	Ctrl_U byte = 21
	Ctrl_V byte = 22
	Ctrl_W byte = 23
	Ctrl_X byte = 24
	Ctrl_Y byte = 25
	Ctrl_Z byte = 26
	Escape byte = 27
)
