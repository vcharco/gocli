package gocli

import (
	gg "github.com/vcharco/gocli/internal/core"
	gt "github.com/vcharco/gocli/internal/types"
	gu "github.com/vcharco/gocli/internal/utils"
)

type Terminal = gg.Terminal
type TerminalResponse = gg.TerminalResponse
type TerminalStyles = gg.TerminalStyles
type TerminalResponseType = gg.TerminalResponseType
type Command = gt.Command
type Param = gt.Param
type ParamModifier = gt.ParamModifier
type ParamType = gt.ParamType

const (
	Date        = gt.Date
	Domain      = gt.Domain
	Email       = gt.Email
	Ipv4        = gt.Ipv4
	Ipv6        = gt.Ipv6
	Number      = gt.Number
	FloatNumber = gt.FloatNumber
	Phone       = gt.Phone
	Text        = gt.Text
	Time        = gt.Time
	Url         = gt.Url
	UUID        = gt.UUID
)

const (
	Cmd            = gg.Cmd
	OsCmd          = gg.OsCmd
	CmdHelp        = gg.CmdHelp
	CmdError       = gg.CmdError
	CtrlKey        = gg.CtrlKey
	ParamError     = gg.ParamError
	ExecutionError = gg.ExecutionError
)

const (
	Reset           = gu.Reset
	BgReset         = gu.BgReset
	Red             = gu.Red
	Green           = gu.Green
	Yellow          = gu.Yellow
	Blue            = gu.Blue
	Magenta         = gu.Magenta
	Cyan            = gu.Cyan
	White           = gu.White
	Black           = gu.Black
	Gray            = gu.Gray
	DarkGray        = gu.DarkGray
	BrightRed       = gu.BrightRed
	BrightGreen     = gu.BrightGreen
	BrightYellow    = gu.BrightYellow
	BrightBlue      = gu.BrightBlue
	BrightMagenta   = gu.BrightMagenta
	BrightCyan      = gu.BrightCyan
	BrightWhite     = gu.BrightWhite
	LightBlue       = gu.LightBlue
	LightGreen      = gu.LightGreen
	LightYellow     = gu.LightYellow
	LightRed        = gu.LightRed
	LightMagenta    = gu.LightMagenta
	LightCyan       = gu.LightCyan
	LightGray       = gu.LightGray
	BgTransparent   = gu.BgTransparent
	BgRed           = gu.BgRed
	BgGreen         = gu.BgGreen
	BgYellow        = gu.BgYellow
	BgBlue          = gu.BgBlue
	BgMagenta       = gu.BgMagenta
	BgCyan          = gu.BgCyan
	BgWhite         = gu.BgWhite
	BgBlack         = gu.BgBlack
	BgBrightRed     = gu.BgBrightRed
	BgBrightGreen   = gu.BgBrightGreen
	BgBrightYellow  = gu.BgBrightYellow
	BgBrightBlue    = gu.BgBrightBlue
	BgBrightMagenta = gu.BgBrightMagenta
	BgBrightCyan    = gu.BgBrightCyan
	BgBrightWhite   = gu.BgBrightWhite
	BgLightBlue     = gu.BgLightBlue
	BgLightGreen    = gu.BgLightGreen
	BgLightYellow   = gu.BgLightYellow
	BgLightRed      = gu.BgLightRed
	BgLightMagenta  = gu.BgLightMagenta
	BgLightCyan     = gu.BgLightCyan
	BgLightGray     = gu.BgLightGray
	BgGray          = gu.BgGray
	BgDarkGray      = gu.BgDarkGray
)

const (
	CursorBlock     = gu.CursorBlock
	CursorBar       = gu.CursorBar
	CursorUnderline = gu.CursorUnderline
)

const (
	DEFAULT  = gt.DEFAULT
	REQUIRED = gt.REQUIRED
)

const (
	Ctrl_A = gt.Ctrl_A
	Ctrl_B = gt.Ctrl_B
	Ctrl_C = gt.Ctrl_C
	Ctrl_D = gt.Ctrl_D
	Ctrl_E = gt.Ctrl_E
	Ctrl_F = gt.Ctrl_F
	Ctrl_G = gt.Ctrl_G
	Ctrl_H = gt.Ctrl_H
	Ctrl_I = gt.Ctrl_I
	Ctrl_J = gt.Ctrl_J
	Ctrl_K = gt.Ctrl_K
	Ctrl_L = gt.Ctrl_L
	Ctrl_M = gt.Ctrl_M
	Ctrl_N = gt.Ctrl_N
	Ctrl_O = gt.Ctrl_O
	Ctrl_P = gt.Ctrl_P
	Ctrl_Q = gt.Ctrl_Q
	Ctrl_R = gt.Ctrl_R
	Ctrl_S = gt.Ctrl_S
	Ctrl_T = gt.Ctrl_T
	Ctrl_U = gt.Ctrl_U
	Ctrl_V = gt.Ctrl_V
	Ctrl_W = gt.Ctrl_W
	Ctrl_X = gt.Ctrl_X
	Ctrl_Y = gt.Ctrl_Y
	Ctrl_Z = gt.Ctrl_Z
	Escape = gt.Escape
)
