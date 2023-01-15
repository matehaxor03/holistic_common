package common

import (
	"fmt"
	"os/exec"
	"bufio"
	"io/ioutil"
	"os"
	"time"
	"sync"
	"container/list"
)

type BashCommand struct {
	ExecuteUnsafeCommandBackground func(command string)
	ExecuteUnsafeCommandUsingFiles func(command string, command_data string) ([]string, []error)
	ExecuteUnsafeCommandUsingFilesWithoutInputFile func(command string) ([]string, []error)
	ExecuteUnsafeCommandSimple func(command string)
}

func NewBashCommand() *BashCommand {
	const maxCapacity = 10*1024*1024  
	delete_files := list.New()
	lock := &sync.RWMutex{}
	status_lock := &sync.RWMutex{}
	file_lock := &sync.RWMutex{}
	var wg sync.WaitGroup
	wakeup_lock := &sync.Mutex{}
	status := "running"

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

	get_or_set_status := func(s string) string {
		status_lock.Lock()
		defer status_lock.Unlock()
		if s == "" {
			return status
		} else {
			status = s
			return ""
		}
	}

	wakeup_delete_file_processor := func() {
		wakeup_lock.Lock()
		defer wakeup_lock.Unlock()
		if get_or_set_status("") == "paused" {
			get_or_set_status("try again") 
			wg.Done()
		} else {
			get_or_set_status("try again") 
		}
	}

	get_or_set_files := func(absolute_path_filename *string, mode string) (*string, error) {
		file_lock.Lock()
		defer file_lock.Unlock()
		if mode == "push" {
			if absolute_path_filename == nil {
				return nil, fmt.Errorf("absolute_path_filename is nil")
			}
			delete_files.PushFront(*absolute_path_filename)
			wakeup_delete_file_processor()
			return nil, nil
		} else if mode == "pull" {
			message := delete_files.Front()
			if message == nil {
				return nil, nil
			}
			delete_files.Remove(message)
			return message.Value.(*string), nil
		} else {
			return nil, fmt.Errorf("mode is not supported %s", mode)
		}
	}

	cleanup_files := func(input_file string, stdout_file string, std_err_file string) {
		if input_file != "" {
			get_or_set_files(&input_file, "push")
		}
		get_or_set_files(&stdout_file, "push")
		get_or_set_files(&std_err_file, "push")

		/*if input_file != "" {
			os.Remove(input_file)
		}
		os.Remove(stdout_file)
		os.Remove(std_err_file)	*/
	}

	

	process_clean_up := func() {
		for {
			get_or_set_status("running")
			time.Sleep(1 * time.Nanosecond) 
			absolute_path_filename_to_delete, absolute_path_filename_to_delete_error := get_or_set_files(nil, "pull")
			if absolute_path_filename_to_delete_error != nil {
				fmt.Println(absolute_path_filename_to_delete_error)
			} else if absolute_path_filename_to_delete != nil {
				os.Remove(*absolute_path_filename_to_delete)
			} else if get_or_set_status("") == "running" {
				wg.Add(1)
				get_or_set_status("paused")
				wg.Wait()
				get_or_set_status("running")
			}
		}
	}
	
	x := BashCommand{
		ExecuteUnsafeCommandBackground: func(command string) {
			lock.Lock()
			defer lock.Unlock()
			cmd := exec.Command("bash", "-c", command + " &")
			cmd.Start()
		},
		ExecuteUnsafeCommandSimple: func(command string) {
			lock.Lock()
			defer lock.Unlock()
			execute_unsafe_command_simple(command)
		},
		ExecuteUnsafeCommandUsingFiles: func(command string, command_data string) ([]string, []error) {
			lock.Lock()
			defer lock.Unlock()
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
		ExecuteUnsafeCommandUsingFilesWithoutInputFile: func(command string) ([]string, []error) {
			lock.Lock()
			defer lock.Unlock()
			var errors []error
			uuid, _ := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
			time_now := time.Now().UnixNano()
			filename_stdout := directory + "/" + fmt.Sprintf("%v%s-stdout.sql", time_now, string(uuid))
			filename_stderr := directory + "/" + fmt.Sprintf("%v%s-stderr.sql", time_now, string(uuid))
			defer cleanup_files("", filename_stdout, filename_stderr)

			full_command := command +  " > " + filename_stdout + " 2> " + filename_stderr + " | touch " + filename_stdout + " && touch " + filename_stderr
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

	go process_clean_up()

	return &x
}
