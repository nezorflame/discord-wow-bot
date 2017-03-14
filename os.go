package main

import (
	"os/exec"
)

// ExecuteCommand - function which runs the shell command
func ExecuteCommand(cmd string) (output string, err error) {
	var bOutput []byte

	bOutput, err = exec.Command(cmd).Output()
	output = string(bOutput)

	return
}
