package common

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"bufio"
)

type BashCommand struct {
	ExecuteUnsafeCommand func(command string, stdout_callback *func(message string), stderr_callback *func(message error)) (*[]string, []error)
}

func NewBashCommand() *BashCommand {
	x := BashCommand{
		ExecuteUnsafeCommand: func(command string, stdout_callback *func(message string), stderr_callback *func(message error)) (*[]string, []error) {
			var errors []error
			stdout_array := []string{}
			stderr_array := []error{}

			cmd := exec.Command("bash", "-c", command)

			cmd_stdout_reader, cmd_stdout_reader_error := cmd.StdoutPipe()
			if cmd_stdout_reader_error != nil {
				errors = append(errors, cmd_stdout_reader_error)
			}

			cmd_stderr_reader, cmd_stderr_reader_error := cmd.StderrPipe()
			if cmd_stderr_reader_error != nil {
				errors = append(errors, cmd_stderr_reader_error)
			}

			if len(errors) > 0 {
				return nil, errors
			}

		
			
			var wg_stdout sync.WaitGroup
			wg_stdout.Add(1)
			stdout_scanner := bufio.NewScanner(cmd_stdout_reader)
			go func() {
				wg_stdout.Done()
				for stdout_scanner.Scan() {
					wg_stdout.Add(1)
					text := strings.TrimSpace(stdout_scanner.Text())
					stdout_array = append(stdout_array, text)
					if stdout_callback != nil {
						(*stdout_callback)(text)
					}
					wg_stdout.Done()
				}
			}()

			var wg_stderr sync.WaitGroup
			wg_stderr.Add(1)
			stderr_scanner := bufio.NewScanner(cmd_stderr_reader)
			go func() {
				wg_stderr.Done()
				for stderr_scanner.Scan() {
					wg_stderr.Add(1)
					text := strings.TrimSpace(stderr_scanner.Text())
					temp_error := fmt.Errorf(text)
					stderr_array = append(stderr_array, fmt.Errorf(text))
					if stderr_callback != nil {
						(*stderr_callback)(temp_error)
					}
					wg_stderr.Done()
				}
			}()


			command_start_err := cmd.Start()
			if command_start_err != nil {
				errors = append(errors, command_start_err)
			}
			
			if len(errors) > 0 {
				return nil, errors
			}

			command_wait_err := cmd.Wait()
			if command_wait_err != nil {
				errors = append(errors, command_wait_err)
			}
			
			if len(errors) > 0 {
				return nil, errors
			}

			wg_stdout.Wait()
			wg_stderr.Wait()

			errors = append(errors, stderr_array...)

			if len(errors) > 0 {
				return &stdout_array, errors
			}

			return &stdout_array, nil
		},
	}

	return &x
}