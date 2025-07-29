// +build linux darwin

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
	"syscall"
)

// execInheritStreams executes the program whose path is specified, passing
// the arguments passed.  The standard input, output, and error streams
// for the subprocess are wired up to the current (parent) process.  This
// also handles SIGINT and SIGTERM signals by killing the command subprocess 
// and any processes it creates.
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
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Start a go routine that listens for SIGINT and SIGTERM
	// signals and kills the command process and any subprocesses
	// it creates.

	cmd.Start()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)
	signal.Notify(signalChannel, syscall.SIGTERM)

	go func() {

		// Wait for a signal.
		<-signalChannel

		// Kill the subprocess and any subprocesses it creates.
		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err == nil {
			syscall.Kill(-pgid, syscall.SIGKILL)
		}
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
