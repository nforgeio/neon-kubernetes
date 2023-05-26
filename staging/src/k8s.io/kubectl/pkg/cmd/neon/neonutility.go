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
	"path"
	"strings"
)

// NeonCliExec locates the [neon-cli] executable and then executes it, passing
// the specified arguments.  The current process will be terminated with a (-1)
// exit code if the executable couldn't be located.  The standard input, output,
// and error streams for the current process are redirected to the subprocess.
//
// This function does not return.  The current process exits, returning the
// exitcode returned by the subprocess.
func NeonCliExec(args []string) {

	ExecInheritStreams(getNeonCliPath(), args)
}

// HelmExec locates the [helm] executable and then executes it, passing the
// the specified arguments.  The current process will be terminated with a (-1)
// exit code if the executable couldn't be located.  The standard input, output,
// and error streams for the current process are redirected to the subprocess.
//
// This function does not return.  The current process exits, returning the
// exitcode returned by the subprocess.
func HelmExec(args []string) {

	// Locate the [neon-cli.exe] binary, handling two possible scenarios:
	//
	// neon-cli/neon-desktop is installed on the current machine:
	// ----------------------------------------------------------
	// In this scenario, the NEON_INSTALL_FOLDER environment variable will
	// be present and will reference the folder holding the application
	// binaries.  Tools like Helm will be located in the [tools] subfolder.
	//
	// neon-cli/neon-desktop is not installed:
	// ---------------------------------------
	// In this case, a maintainer is probably debugging or is otherwise using
	// [neon-cli] without it being formally installed.  We're going to execute
	// the [neon-cli toolpath helm] command which will attempt to locate the
	// Helm binary and try to download it when it's not found, returning its
	// full path as the command output.

	helmPath := ""
	neonInstallFolder := os.Getenv("NEON_INSTALL_FOLDER")

	if neonInstallFolder != "" {
		helmPath = path.Join(neonInstallFolder, "tools", "helm.exe")
	} else {
		// Note that we can't use [NeonCliExec()] here because that redirects
		// the standard streams and never returns.

		neonCliPath := getNeonCliPath()

		neonCliToolPathCmd := exec.Command(neonCliPath, "toolpath", "helm")
		neonCliToolPathCmd.Env = os.Environ()
		bytes, err := neonCliToolPathCmd.Output()
		if err != nil {
			panic(err)
		}

		helmPath = string(bytes[:])

		if helmPath == "" {
			panic(errors.New("cannot locate the helm binary"))
		}

		helmPath = strings.TrimSpace(string(bytes[:]))
	}

	ExecInheritStreams(helmPath, args)
}

// fileExists returns TRUE when a specified file exists.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// getNeonCliPath attempts to locate the [neon-cli] binary.
func getNeonCliPath() string {

	// Locate the [neon-cli.exe] binary, handling two possible scenarios:
	//
	// neon-cli/neon-desktop is installed on the current machine:
	// ----------------------------------------------------------
	// In this scenario, the NEON_INSTALL_FOLDER environment variable will
	// be present and will reference the folder holding the application
	// binaries, including [neon-cli], so we'll execute [neon-cli] from
	// there.
	//
	// neon-cli/neon-desktop is not installed:
	// ---------------------------------------
	// The NEON_INSTALL_FOLDER environment variable will not be present for
	// this case.  We're going to check if the NC_ROOT environment variable
	// pointing to the NEONCLOUD source repo exists and try to locate the
	// [neon-cli] binary there.  We're going to use the first binary found
	// in one of (searched in the listed order):
	//
	//		1. {NC_ROOT}/Build/neon-cli
	//		2. {NC_ROOT}/Tools/neon-cli/bin/Debug
	//		2. {NC_ROOT}/Tools/neon-cli/bin/Release
	//
	// NOTE: We'll need to update the hardcoded subfolder paths when we
	//       update .NET SDKs or we target another version of Windows.

	neonCliPath := ""
	neonInstallFolder := os.Getenv("NEON_INSTALL_FOLDER")

	if neonInstallFolder != "" {
		neonCliPath = path.Join(neonInstallFolder, "neoncli.exe")
	} else {
		ncRoot := os.Getenv("NC_ROOT")
		if ncRoot != "" {
			neonCliBuildPath := path.Join(ncRoot, "Build", "neon-cli", "neon-cli.exe")
			neonCliDebugPath := path.Join(ncRoot, "Tools", "neon-cli", "bin", "Debug", "net7.0-windows10.0.17763.0", "win10-x64", "neon-cli.exe")
			neonCliReleasePath := path.Join(ncRoot, "Tools", "neon-cli", "bin", "Release", "net7.0-windows10.0.17763.0", "win10-x64", "neon-cli.exe")

			if fileExists(neonCliBuildPath) {
				neonCliPath = neonCliBuildPath
			} else if fileExists(neonCliDebugPath) {
				neonCliPath = neonCliDebugPath
			} else if fileExists(neonCliReleasePath) {
				neonCliPath = neonCliReleasePath
			}
		}
	}

	if neonCliPath == "" || !fileExists(neonCliPath) {
		panic(errors.New("cannot locate the [neon-cli] binary"))
	} else {
		return neonCliPath
	}
}

// ExecInheritStreams executes the program whose path is specified, passing
// the specified arguments.  The standard input, output, and error streams
// for the subprocess are wired up to the current (parent) process.
//
// This function does not return.  The current process exits, returning the
// exitcode returned by the subprocess or it panics when the executable
// could not be launched.
func ExecInheritStreams(path string, args []string) {

	// Attempt to execute the command.

	cmd := exec.Command(path, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

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
