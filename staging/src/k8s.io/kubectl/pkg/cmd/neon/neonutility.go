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
	"time"
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

type pathInfo struct {
	path      string
	timestamp time.Time
}

// appendPathInfo appends a pathInfo struct to the paths slice when the file
// at the specified path exists and returns the new slice.  If the file doesn't
// exist, the function returns the unmodified slice.
func appendPathInfo(paths []pathInfo, path string) []pathInfo {

	info, err := os.Stat(path)
	if os.IsNotExist(err) || info.IsDir() {
		return paths
	}

	return append(paths, pathInfo{path: path, timestamp: info.ModTime()})
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
	// this case.  We're going to look for the most recently built binary
	// that exists at these locations andf use that.
	//
	// $(NC_ROOT)/Build/neon-cli/neon-cli.exe
	// $(NC_ROOT)/Tools/neon-cli/bin/Debug/net7.0-windows10.0.17763.0/win10-x64/neon-cli.exe
	// $(NC_ROOT)/Tools/neon-cli/bin/Release/net7.0-windows10.0.17763.0/win10-x64/neon-cli.exe
	// $(NK_ROOT)/Tools/neon-cli/bin/Debug/net7.0-windows10.0.17763.0/win10-x64/neon-cli.exe
	// $(NK_ROOT)/Tools/neon-cli/bin/Debug/net7.0-windows10.0.17763.0/win10-x64/neon-cli.exe
	//
	// NOTE: We'll need to update the hardcoded subfolder paths when we
	//       update .NET SDKs or we target another version of Windows.

	const frameworkMoniker = "net7.0-windows10.0.17763.0"
	const architecture = "win10-x64"

	// Create a slice with information about the candidate executables.

	candidates := make([]pathInfo, 10)

	neonInstallFolder := os.Getenv("NEON_INSTALL_FOLDER")
	if neonInstallFolder != "" {
		candidates = appendPathInfo(candidates, path.Join(neonInstallFolder, "neoncli.exe"))
	}

	ncRoot := os.Getenv("NC_ROOT")
	if ncRoot != "" {

		candidates = appendPathInfo(candidates, path.Join(ncRoot, "Build", "neon-cli", "neon-cli.exe"))
		candidates = appendPathInfo(candidates, path.Join(ncRoot, "Tools", "neon-cli", "bin", "Debug", frameworkMoniker, architecture, "neon-cli.exe"))
		candidates = appendPathInfo(candidates, path.Join(ncRoot, "Tools", "neon-cli", "bin", "Release", frameworkMoniker, architecture, "neon-cli.exe"))
	}

	nkRoot := os.Getenv("NK_ROOT")
	if nkRoot != "" {

		candidates = appendPathInfo(candidates, path.Join(nkRoot, "Tools", "neon-cli", "bin", "Debug", frameworkMoniker, architecture, "neon-cli.exe"))
		candidates = appendPathInfo(candidates, path.Join(nkRoot, "Tools", "neon-cli", "bin", "Release", frameworkMoniker, architecture, "neon-cli.exe"))
	}

	if len(candidates) == 0 {
		panic(errors.New("cannot locate the [neon-cli] binary"))
	}

	// Look for the candidate executable with the most recent timestamp.

	latestCandidate := candidates[0]

	for _, candidate := range candidates {
		if candidate.timestamp.After(latestCandidate.timestamp) {
			latestCandidate = candidate
		}
	}

	return latestCandidate.path
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
