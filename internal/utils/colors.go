package gocliutils

type Color string
type BgColor string

const (
	// Reset colors
	Reset Color = "\033[0m"

	// Regular text colors
	Red      Color = "\033[31m"
	Green    Color = "\033[32m"
	Yellow   Color = "\033[33m"
	Blue     Color = "\033[34m"
	Magenta  Color = "\033[35m"
	Cyan     Color = "\033[36m"
	White    Color = "\033[37m"
	Black    Color = "\033[30m"
	Gray     Color = "\033[90m"
	DarkGray Color = "\033[38;5;235m"

	// Bright text colors
	BrightRed     Color = "\033[91m"
	BrightGreen   Color = "\033[92m"
	BrightYellow  Color = "\033[93m"
	BrightBlue    Color = "\033[94m"
	BrightMagenta Color = "\033[95m"
	BrightCyan    Color = "\033[96m"
	BrightWhite   Color = "\033[97m"

	// Regular background colors
	BgRed      BgColor = "\033[41m"
	BgGreen    BgColor = "\033[42m"
	BgYellow   BgColor = "\033[43m"
	BgBlue     BgColor = "\033[44m"
	BgMagenta  BgColor = "\033[45m"
	BgCyan     BgColor = "\033[46m"
	BgWhite    BgColor = "\033[47m"
	BgBlack    BgColor = "\033[40m"
	BgGray     BgColor = "\033[48;5;235m"
	BgDarkGray BgColor = "\033[48;5;236m"

	// Bright background colors
	BgBrightRed     BgColor = "\033[101m"
	BgBrightGreen   BgColor = "\033[102m"
	BgBrightYellow  BgColor = "\033[103m"
	BgBrightBlue    BgColor = "\033[104m"
	BgBrightMagenta BgColor = "\033[105m"
	BgBrightCyan    BgColor = "\033[106m"
	BgBrightWhite   BgColor = "\033[107m"

	// Extended background colors (256 colors)
	BgLightBlue    BgColor = "\033[48;5;153m"
	BgLightGreen   BgColor = "\033[48;5;120m"
	BgLightYellow  BgColor = "\033[48;5;228m"
	BgLightRed     BgColor = "\033[48;5;203m"
	BgLightMagenta BgColor = "\033[48;5;213m"
	BgLightCyan    BgColor = "\033[48;5;159m"
)

func ColorizeForeground(c Color, text string) string {
	return string(c) + text + string(Reset)
}

func ColorizeBackground(c Color, text string) string {
	return string(c) + text + string(Reset)
}

func ColorizeBoth(c Color, bc BgColor, text string) string {
	return string(bc) + string(c) + text + string(Reset)
}
