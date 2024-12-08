package gocliutils

import (
	"fmt"
	"os/exec"
	"runtime"
)

func GetClipboardContent() (string, error) {
	var cmd *exec.Cmd
	var err error
	var errMsg string

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbpaste")
		errMsg = "Please, install pbpaste."
	case "linux":
		cmd = exec.Command("xclip", "-o")
		errMsg = "Please, install xclip (sudo apt install xclip)"
	case "windows":
		cmd = exec.Command("powershell", "-Command", "Get-Clipboard")
		errMsg = "Please, install powershell"
	default:
		return "", fmt.Errorf("read clipboard not supported on your operative system")
	}

	content, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("cannot get the clipboard. %v", errMsg), nil
	}

	return string(content), nil
}

func SetClipboard(content string) error {
	var cmd *exec.Cmd
	var errMsg string
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
		errMsg = "Please, install pbcopy."
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
		errMsg = "Please, install xclip (sudo apt install xclip)"
	case "windows":
		cmd = exec.Command("clip")
		errMsg = "Please, install clip"
	default:
		return fmt.Errorf("write clipboard not supported on your operative system")
	}

	pipe, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("cannot set the clipboard: %v", errMsg)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cannot set the clipboard: %v", errMsg)
	}

	_, err = pipe.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("cannot set the clipboard: %v", errMsg)
	}

	pipe.Close()
	return cmd.Wait()
}
