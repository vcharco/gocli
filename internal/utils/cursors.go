package gocliutils

type Cursor string

type CursorColor string

const (
	CursorRed      CursorColor = "\033[31m"
	CursorGreen    CursorColor = "\033[32m"
	CursorYellow   CursorColor = "\033[33m"
	CursorBlue     CursorColor = "\033[34m"
	CursorMagenta  CursorColor = "\033[35m"
	CursorCyan     CursorColor = "\033[36m"
	CursorWhite    CursorColor = "\033[37m"
	CursorBlack    CursorColor = "\033[30m"
	CursorGray     CursorColor = "\033[90m"
	CursorDarkGray CursorColor = "\033[38;5;235m"
	CursorBrightRed     CursorColor = "\033[91m"
	CursorBrightGreen   CursorColor = "\033[92m"
	CursorBrightYellow  CursorColor = "\033[93m"
	CursorBrightBlue    CursorColor = "\033[94m"
	CursorBrightMagenta CursorColor = "\033[95m"
	CursorBrightCyan    CursorColor = "\033[96m"
	CursorBrightWhite   CursorColor = "\033[97m"
	CursorLightBlue    CursorColor = "\033[38;5;153m"
	CursorLightGreen   CursorColor = "\033[38;5;120m"
	CursorLightYellow  CursorColor = "\033[38;5;228m"
	CursorLightRed     CursorColor = "\033[38;5;203m"
	CursorLightMagenta CursorColor = "\033[38;5;213m"
	CursorLightCyan    CursorColor = "\033[38;5;159m"
)

const (
	CursorBlock     = "\033[2 q" // Cursor bloque
	CursorBar       = "\033[6 q" // Cursor barra
	CursorUnderline = "\033[4 q" // Cursor subrayado
	CursorReset     = "\033[0m"  // Restablecer cursor
)
