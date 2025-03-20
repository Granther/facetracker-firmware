package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// Takes a command as a single string, first word is the command name, the remaining are the args
func ExecCommand(command string) error {
	commandArr := strings.Split(command, " ")
	cmd := exec.Command(commandArr[0], commandArr[1:]...)
	_, err := cmd.CombinedOutput()
	return err
}
