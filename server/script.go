package server

import (
	"fmt"
	"os/exec"
	"runtime"
)

func executeScript(body []byte) ([]byte, error) {
	var command string
	switch runtime.GOOS {
	case "linux", "darwin":
		command = "/bin/sh"
	default:
		return nil, fmt.Errorf("agent running on an unknown os %s", runtime.GOOS)
	}

	cmd := exec.Command(command)

	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		defer in.Close()
		in.Write(body)
	}()

	return cmd.CombinedOutput()
}
