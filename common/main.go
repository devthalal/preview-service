package common

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func RunCmd(dir, execCmd string, cmdArgs ...string) error {
	cmd := exec.Command(execCmd, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if len(dir) > 0 {
		cmd.Dir = dir
	}

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Printf("Exit error: %s\n", exitErr)
			if exitErr.ExitCode() == 130 {
				return fmt.Errorf("error while running cmd: Control-C")
			}
		}
		return fmt.Errorf("error while running cmd: %s", err)
	}

	return nil
}
