package main

import (
	"os/exec"
)

// ExecuteCommand - function which runs the shell command
func ExecuteCommand(command, dir string, args []string) (output string, err error) {
	var bOutput []byte
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	bOutput, err = cmd.Output()
	output = string(bOutput)

	return
}
