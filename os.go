package main

import (
	"os/exec"
)

// ExecuteCommand - function which runs the shell command
func ExecuteCommand(command string) (output string, err error) {
	var bOutput []byte
	cmd := exec.Command(command)
	cmd.Dir = o.SimcDir
	bOutput, err = cmd.Output()
	output = string(bOutput)

	return
}
