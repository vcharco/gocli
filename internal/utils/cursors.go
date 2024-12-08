package gocliutils

type Cursor string

const (
	CursorBlock     = "\033[2 q" // Cursor bloque
	CursorBar       = "\033[6 q" // Cursor barra
	CursorUnderline = "\033[4 q" // Cursor subrayado
	CursorReset     = "\033[0m"  // Restablecer cursor
)
