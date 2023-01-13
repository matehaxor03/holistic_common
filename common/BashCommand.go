package common

import (
	"fmt"
	"os/exec"
	"sync"
	"bufio"
	"io/ioutil"
	"os"
	"time"
)

type BashCommand struct {
	ExecuteUnsafeCommand func(command string, stdout_callback *func(message string), stderr_callback *func(message error)) (*[]string, []error)
	ExecuteUnsafeCommandBackground func(command string)
	ExecuteUnsafeCommandUsingFiles func(command string, command_data string) ([]string, []error)
	ExecuteUnsafeCommandSimple func(command string)
}

func NewBashCommand() *BashCommand {
	const maxCapacity = 10*1024*1024  
	directory_parts := GetDataDirectory()
	directory := "/" 
	for index, directory_part := range directory_parts {
		directory += directory_part
		if index < len(directory_parts) - 1 {
			directory += "/"
		}
	}

	execute_unsafe_command_simple := func(command string) {
		cmd := exec.Command("bash", "-c", command)
		cmd.Start()
		cmd.Wait()
	}

	cleanup_files := func(input_file string, stdout_file string, std_err_file string) {
		os.Remove(input_file)
		os.Remove(stdout_file)
		os.Remove(std_err_file)	
	}
	
	x := BashCommand{
		ExecuteUnsafeCommand: func(command string, stdout_callback *func(message string), stderr_callback *func(message error)) (*[]string, []error) {
			var errors []error
			stdout_array := []string{}
			stderr_array := []error{}

			cmd := exec.Command("bash", "-c", command)

			var errors_found = false
			var wg_stdout sync.WaitGroup
			var wg_stderr sync.WaitGroup

			if stdout_callback != nil && stderr_callback != nil {
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
				
				wg_stdout.Add(1)
				stdout_scanner := bufio.NewScanner(cmd_stdout_reader)
				stdout_scanner_buffer := make([]byte, maxCapacity)
				stdout_scanner.Buffer(stdout_scanner_buffer, maxCapacity)
				stdout_scanner.Split(bufio.ScanLines)
				go func() {
					for stdout_scanner.Scan() {
						text := string(stdout_scanner.Text())
						stdout_array = append(stdout_array, text)
						if stdout_callback != nil {
							go (*stdout_callback)(text)
						}
					}
					wg_stdout.Done()
				}()

				wg_stderr.Add(1)
				stderr_scanner := bufio.NewScanner(cmd_stderr_reader)
				stderr_scanner_buffer := make([]byte, maxCapacity)
				stderr_scanner.Buffer(stderr_scanner_buffer, maxCapacity)
				stderr_scanner.Split(bufio.ScanLines)
				go func() {
					for stderr_scanner.Scan() {
						errors_found = true
						text := string(stderr_scanner.Text())
						temp_error := fmt.Errorf(text)
						stderr_array = append(stderr_array, fmt.Errorf(text))
						if stderr_callback != nil {
							go (*stderr_callback)(temp_error)
						}
					}
					wg_stderr.Done()
				}()
			}

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

			if stdout_callback != nil {
				wg_stdout.Wait()
			}

			if stderr_callback != nil {
				if errors_found {
					wg_stderr.Wait()
					errors = append(errors, stderr_array...)
				}
			}
			

			if len(errors) > 0 {
				return &stdout_array, errors
			}

			return &stdout_array, nil
		},
		ExecuteUnsafeCommandBackground: func(command string) {
			cmd := exec.Command("bash", "-c", command + " &")
			cmd.Start()
		},
		ExecuteUnsafeCommandSimple: func(command string) {
			execute_unsafe_command_simple(command)
		},
		ExecuteUnsafeCommandUsingFiles: func(command string, command_data string) ([]string, []error) {
			var errors []error
			uuid, _ := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
			time_now := time.Now().UnixNano()
			filename := directory + "/" + fmt.Sprintf("%v%s.sql", time_now, string(uuid))
			filename_stdout := directory + "/" + fmt.Sprintf("%v%s-stdout.sql", time_now, string(uuid))
			filename_stderr := directory + "/" + fmt.Sprintf("%v%s-stderr.sql", time_now, string(uuid))
			defer cleanup_files(filename, filename_stdout, filename_stderr)

			ioutil.WriteFile(filename, []byte(command_data), 0600)
			full_command := command + " < " + filename +  " > " + filename_stdout + " 2> " + filename_stderr + " | touch " + filename_stdout + " && touch " + filename_stderr
			execute_unsafe_command_simple(full_command)

			var stdout_lines []string
			if _, opening_stdout_error := os.Stat(filename_stdout); opening_stdout_error == nil {
				file_stdout, file_stdout_errors := os.Open(filename_stdout)
				if file_stdout_errors != nil {
					errors = append(errors, file_stdout_errors)
				} else {
					defer file_stdout.Close()
					stdout_scanner := bufio.NewScanner(file_stdout)
					stdout_scanner_buffer := make([]byte, maxCapacity)
					stdout_scanner.Buffer(stdout_scanner_buffer, maxCapacity)
					stdout_scanner.Split(bufio.ScanLines)
					for stdout_scanner.Scan() {
						current_text := stdout_scanner.Text()
						if current_text != "" {
							stdout_lines = append(stdout_lines, current_text)
						}
					}
				}
			}

			if _, opening_stderr_error := os.Stat(filename_stderr); opening_stderr_error == nil {
				file_stderr, file_stderr_errors := os.Open(filename_stderr)
				if file_stderr_errors != nil {
					errors = append(errors, file_stderr_errors)
				} else {
					defer file_stderr.Close()
					stderr_scanner := bufio.NewScanner(file_stderr)
					stderr_scanner_buffer := make([]byte, maxCapacity)
					stderr_scanner.Buffer(stderr_scanner_buffer, maxCapacity)
					stderr_scanner.Split(bufio.ScanLines)
					for stderr_scanner.Scan() {
						current_text := stderr_scanner.Text()
						if current_text != "" {
							errors = append(errors, fmt.Errorf("%s", current_text))
						}
					}
				}
			}

			if len(errors) > 0 {
				return stdout_lines, errors
			}

			return stdout_lines, nil
		},
	}

	return &x
}
