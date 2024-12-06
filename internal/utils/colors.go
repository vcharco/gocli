package gocliutils

type Color int

const (
    Red Color = iota + 31
    Green
    Yellow
    Blue
    Magenta
    Cyan
    White
    Reset = 0
)

var colorCodes = map[Color]string{
    Red:     "\033[31m",
    Green:   "\033[32m",
    Yellow:  "\033[33m",
    Blue:    "\033[34m",
    Magenta: "\033[35m",
    Cyan:    "\033[36m",
    White:   "\033[37m",
    Reset:   "\033[0m",
}

func Colorize(c Color, text string) string {
    return colorCodes[c] + text + colorCodes[Reset]
}
