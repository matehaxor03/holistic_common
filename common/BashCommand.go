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

			const maxCapacity = 10*1024*1024  
			
			var wg_stdout sync.WaitGroup
			wg_stdout.Add(1)
			stdout_scanner := bufio.NewScanner(cmd_stdout_reader)
			stdout_scanner_buffer := make([]byte, maxCapacity)
			stdout_scanner.Buffer(stdout_scanner_buffer, maxCapacity)
			stdout_scanner.Split(bufio.ScanLines)
			go func() {
				for stdout_scanner.Scan() {
					text := strings.TrimSpace(string(stdout_scanner.Text()))
					stdout_array = append(stdout_array, text)
					if stdout_callback != nil {
						(*stdout_callback)(text)
					}
				}
				wg_stdout.Done()
			}()

			var wg_stderr sync.WaitGroup
			wg_stderr.Add(1)
			stderr_scanner := bufio.NewScanner(cmd_stderr_reader)
			stderr_scanner_buffer := make([]byte, maxCapacity)
			stderr_scanner.Buffer(stderr_scanner_buffer, maxCapacity)
			stderr_scanner.Split(bufio.ScanLines)
			var errors_found = false
			go func() {
				for stderr_scanner.Scan() {
					errors_found = true
					text := strings.TrimSpace(string(stderr_scanner.Text()))
					temp_error := fmt.Errorf(text)
					stderr_array = append(stderr_array, fmt.Errorf(text))
					if stderr_callback != nil {
						(*stderr_callback)(temp_error)
					}
				}
				wg_stderr.Done()
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

			wg_stdout.Wait()
			if errors_found {
				wg_stderr.Wait()
			}
		
			if len(errors) > 0 {
				return nil, errors
			}

			errors = append(errors, stderr_array...)

			if len(errors) > 0 {
				return &stdout_array, errors
			}

			return &stdout_array, nil
		},
	}

	return &x
}
