package common

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type BashCommand struct {
	ExecuteUnsafeCommand func(command string) (*string, []error)
}

func NewBashCommand() *BashCommand {
	x := BashCommand{
		ExecuteUnsafeCommand: func(command string) (*string, []error) {
			var errors []error

			var stdout bytes.Buffer
			var stderr bytes.Buffer

			cmd := exec.Command("bash", "-c", command)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			command_err := cmd.Run()

			shell_output := stdout.String()
			shell_output_errs := stderr.String()

			if command_err != nil {
				errors = append(errors, command_err)
			}

			if strings.TrimSpace(shell_output_errs) != "" {
				errors = append(errors, fmt.Errorf("error: %s", fmt.Sprintf(strings.TrimSpace(shell_output_errs))))
			}

			if len(errors) > 0 {
				return &shell_output, errors
			}

			return &shell_output, nil
		},
	}

	return &x
}
