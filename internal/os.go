package gocliutils

import (
	"fmt"
	"os/exec"
	"runtime"
)

func ExecCmd(command string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-Command", command)
	case "linux", "darwin":
		cmd = exec.Command("bash", "-c", command)
	default:
		fmt.Println("OS not supported")
		return
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error executing command: %v\n", err)
		return
	}

	fmt.Println(string(output))
}
