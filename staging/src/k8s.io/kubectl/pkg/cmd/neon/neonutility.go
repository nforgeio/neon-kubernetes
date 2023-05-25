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
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

// NeonCliExec locates the [neon-cli] executable and then executes it, passing
// the specified arguments.  The current process will be terminated with a (-1)
// exit code if the executable couldn't be located.
//
// This function does not return.  It replaces the current process with the new
// [neon-cli] process.
func NeonCliExec(args []string) {

	neonCliPath, err := getNeonCliPath()
	if err != nil {
		panic(err)
	}

	err = syscall.Exec(neonCliPath, args, os.Environ())
	if err != nil {
		fmt.Fprintf(os.Stderr, "*** ERROR: Cannot launch the [neon-cli] binary.\n")
		os.Exit(-1)
	}
}

// HelmExec locates the [helm] executable and then executes it, passing the
// specified arguments.  The current process will be terminated with a (-1)
// exit code if the executable couldn't be located.
//
// This function does not return.  It replaces the current process with the new
// [helm] process.
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

	fmt.Fprintf(os.Stderr, "HelmExec: 0a: %s\n", neonInstallFolder)
	if neonInstallFolder != "" {
		helmPath = path.Join(neonInstallFolder, "tools", "helm.exe")
	} else {
		// Note that we can't use [NeonCliExec()] here because that replaces
		// the existing process and we need to run [neon-cli] as a subprocess
		// here.

		fmt.Fprintf(os.Stderr, "HelmExec: 0b\n")
		neonCliPath, err := getNeonCliPath()
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(os.Stderr, "HelmExec: 1\n")

		cmd := exec.Command(neonCliPath, "toolpath", "helm")
		fmt.Fprintf(os.Stderr, "HelmExec: 2\n")
		cmd.Env = os.Environ()
		fmt.Fprintf(os.Stderr, "HelmExec: 3\n")
		bytes, err := cmd.Output()
		fmt.Fprintf(os.Stderr, "HelmExec: 4: %s\n", neonCliPath)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(os.Stderr, "HelmExec: 5\n")
		if helmPath == "" {
			panic(errors.New("cannot locate the helm binary"))
		}

		fmt.Fprintf(os.Stderr, "HelmExec: 6\n")
		helmPath = strings.TrimSpace(string(bytes[:]))
		fmt.Fprintf(os.Stderr, "HelmExec: 7\n")
	}

	err := syscall.Exec(helmPath, args, os.Environ())
	if err != nil {
		fmt.Fprintf(os.Stderr, "HelmExec: 8\n")
		fmt.Fprintf(os.Stderr, "*** ERROR: Cannot launch the [helm] binary: %s\n", helmPath)
		os.Exit(-1)
	}
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
func getNeonCliPath() (string, error) {

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

	fmt.Fprintf(os.Stderr, "getNeonCliPath: 0\n")
	neonCliPath := ""
	neonInstallFolder := os.Getenv("NEON_INSTALL_FOLDER")

	fmt.Fprintf(os.Stderr, "getNeonCliPath: 1: %s\n", neonInstallFolder)
	if neonInstallFolder != "" {
		fmt.Fprintf(os.Stderr, "getNeonCliPath: 2\n")
		neonCliPath = path.Join(neonInstallFolder, "neoncli.exe")
	} else {
		fmt.Fprintf(os.Stderr, "getNeonCliPath: 3\n")
		ncRoot := os.Getenv("NC_ROOT")
		if ncRoot != "" {
			neonCliBuildPath := path.Join(ncRoot, "Build", "neon-cli", "neon-cli.exe")
			neonCliDebugPath := path.Join(ncRoot, "Tools", "neon-cli", "bin", "Debug", "net7.0-windows10.0.17763.0", "win10-x64", "neon-cli.exe")
			neonCliReleasePath := path.Join(ncRoot, "Tools", "neon-cli", "bin", "Release", "net7.0-windows10.0.17763.0", "win10-x64", "neon-cli.exe")

			fmt.Fprintf(os.Stderr, "getNeonCliPath: 4\n")
			if fileExists(neonCliBuildPath) {
				fmt.Fprintf(os.Stderr, "getNeonCliPath: 5\n")
				neonCliPath = neonCliBuildPath
			} else if fileExists(neonCliDebugPath) {
				fmt.Fprintf(os.Stderr, "getNeonCliPath: 6\n")
				neonCliPath = neonCliDebugPath
			} else if fileExists(neonCliReleasePath) {
				fmt.Fprintf(os.Stderr, "getNeonCliPath: 7\n")
				neonCliPath = neonCliReleasePath
			}
		}
	}
	fmt.Fprintf(os.Stderr, "getNeonCliPath: 8\n")

	if neonCliPath == "" || !fileExists(neonCliPath) {
		return "", errors.New("cannot locate the [neon-cli] binary")
	} else {
		return neonCliPath, nil
	}
}
