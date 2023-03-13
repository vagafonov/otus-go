package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	environ := os.Environ()
	_ = environ

	for key, value := range env {
		os.Setenv(key, value.Value)
	}

	result, err := exec.Command(cmd[0], cmd[1:]...).Output() //nolint:gosec
	if err != nil {
		return 0
	}

	fmt.Println(string(result))

	return 1
}
