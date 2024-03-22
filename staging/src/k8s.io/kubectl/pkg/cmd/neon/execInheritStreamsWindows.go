// +build windows

/*
Copyright 2023 NEONFORGE LLC.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package neon_utility

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// execInheritStreams executes the program whose path is specified, passing
// the arguments passedarguments.  The standard input, output, and error streams
// for the subprocess are wired up to the current (parent) process.  This also
// handles emulated SIGINT signals by killing the command subprocess and any
// processes it creates.
//
// IMPORTANT:
//
// This function does not return.  The current process exits, returning the
// exitcode returned by the subprocess or it panics when the executable
// could not be launched.
func execInheritStreams(path string, args []string) {

	cmd := exec.Command(path, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the process and then create a ProcessExitGroup/JobObject and 
	// add the process to it, and then start a go routine that listens for
	// an emulated SIGINT signal that kills the command process and any of
	// its subprocesses.

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	cmd.Start()

	processExitGroup, _ := CreateProcessExitGroup()
	processExitGroup.AddProcess(cmd.Process)

	go func() {

		// Wait for a signal.
		<-signalChannel

		// Kill the subprocess and any subprocesses it creates.
		processExitGroup.Dispose()
	}()

	err := cmd.Wait()

	// Assume [exitcode=0] when there's no error.

	if err == nil {
		os.Exit(0)
	}

	// We're going to special case "exit status" errors here by extracting
	// the exit code from the error and terminating the current process with
	// that code.
	//
	// We'll panic for all other errors, like when the executable file doesn't
	// exist or when it isn't a valid executable.

	if strings.HasPrefix(err.Error(), "exit status") {

		var exitError *exec.ExitError

		errors.As(err, &exitError)
		os.Exit(exitError.ExitCode())
	} else {
		panic(err)
	}
}

type process struct {
	Pid    int
	Handle uintptr
}

type ProcessExitGroup windows.Handle

func CreateProcessExitGroup() (ProcessExitGroup, error) {

	handle, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return 0, err
	}

	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	if _, err := windows.SetInformationJobObject(
		handle,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info))); err != nil {
		return 0, err
	}

	return ProcessExitGroup(handle), nil
}

func (group ProcessExitGroup) Dispose() error {

	return windows.CloseHandle(windows.Handle(group))
}

func (group ProcessExitGroup) AddProcess(p *os.Process) error {

	return windows.AssignProcessToJobObject(
		windows.Handle(group),
		windows.Handle((*process)(unsafe.Pointer(p)).Handle))
}
