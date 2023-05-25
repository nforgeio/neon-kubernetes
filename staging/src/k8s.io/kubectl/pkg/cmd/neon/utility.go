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
	// [neon-cli] binary there, first in the projects [bin/Debug] and then
	// in its [bin/Release] folder.
	//
	// NOTE: We'll need to update the hardcoded subfolder paths when we
	//       update .NET SDKs or we target another version of Windows.

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

	if neonInstallFolder != "" {
		helmPath = path.Join(neonInstallFolder, "tools", "helm.exe")
	} else {
		neonCliPath, err := getNeonCliPath()
		if err != nil {
			panic(err)
		}

		cmd := exec.Command(neonCliPath, "toolpath", "helm")
		cmd.Env = os.Environ()
		bytes, err := cmd.Output()
		if err != nil {
			panic(err)
		}

		if helmPath == "" {
			panic(errors.New("cannot locate the helm binary"))
		}

		helmPath = strings.TrimSpace(string(bytes[:]))
	}

	err := syscall.Exec(helmPath, args, os.Environ())
	if err != nil {
		fmt.Fprintf(os.Stderr, "*** ERROR: Cannot launch the [neon-cli] binary.\n")
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

	neonCliPath := ""
	neonInstallFolder := os.Getenv("NEON_INSTALL_FOLDER")

	if neonInstallFolder != "" {
		neonCliPath = path.Join(neonInstallFolder, "neoncli.exe")
	} else {
		ncRoot := os.Getenv("NC_ROOT")

		if ncRoot != "" {

			neonCliPath = path.Join(ncRoot, "Tools", "bin", "Debug", "net7.0-windows10.0.17763.0", "win10-x64", "neon-cli.exe")

			if !fileExists(neonCliPath) {

				neonCliPath = path.Join(ncRoot, "Tools", "bin", "Release", "net7.0-windows10.0.17763.0", "win10-x64", "neon-cli.exe")
			}
		}
	}

	if neonCliPath == "" || !fileExists(neonCliPath) {
		return "", errors.New("cannot locate the [neon-cli] binary")

	} else {
		return neonCliPath, nil
	}
}
