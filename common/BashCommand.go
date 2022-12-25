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
			stdout_reader := bufio.NewReader(cmd_stdout_reader)

			go func() {
				stdout_bytes := []byte{}
				for {
					line, isPrefix, stdout_reader_error := stdout_reader.ReadLine()
					if stdout_reader_error != nil {
						errors = append(errors, stdout_reader_error)
						break
					}
					stdout_bytes = append(stdout_bytes, line...)
					if !isPrefix {
						text := strings.TrimSpace(string(stdout_bytes))
						if stdout_callback != nil {
							(*stdout_callback)(text)
						}
						
						if len(text) > 0 {
							stdout_bytes = []byte{}
						}
					}
				}
				
			/*for stdout_scanner.Scan() {
					text := strings.TrimSpace(stdout_scanner.Text())
					stdout_array = append(stdout_array, text)
					if stdout_callback != nil {
						(*stdout_callback)(text)
					}
				}*/
				wg_stdout.Done()
			}()

			var wg_stderr sync.WaitGroup
			wg_stderr.Add(1)
			stderr_reader := bufio.NewReader(cmd_stderr_reader)

			go func() {
				stderr_bytes := []byte{}
				for {
					line, isPrefix, stderr_reader_error := stderr_reader.ReadLine()
					if stderr_reader_error != nil {
						errors = append(errors, stderr_reader_error)
						break
					}
					stderr_bytes = append(stderr_bytes, line...)
					if !isPrefix {
						text := strings.TrimSpace(string(stderr_bytes))
						
						if stderr_callback != nil {
							(*stderr_callback)(fmt.Errorf(text))
						}
						
						if len(text) > 0 {
							stderr_bytes = []byte{}
						}
					}
				}

				 /* 
				for stderr_scanner.Scan() {
					

					text := strings.TrimSpace(stderr_scanner.Text())
					temp_error := fmt.Errorf(text)
					stderr_array = append(stderr_array, fmt.Errorf(text))
					if stderr_callback != nil {
						(*stderr_callback)(temp_error)
					}
				}*/
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
			wg_stderr.Wait()
			
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
