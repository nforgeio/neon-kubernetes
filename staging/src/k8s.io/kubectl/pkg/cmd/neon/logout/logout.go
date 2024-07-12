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

package neon_logout

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	neon_utility "k8s.io/kubectl/pkg/cmd/neon"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	logoutLong = templates.LongDesc(i18n.T(`
	Logs out of the current NEONKUBE context by selecting no context.`))

	logoutExample = templates.Examples(i18n.T(`
		# Logout of the current NEONKUBE cluster context
		neon logout`))
)

// NewCmdNeonLogout returns a Command instance for NEON-CLI 'logout' sub command
func NewCmdNeonLogout(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logout",
		Short:   i18n.T("Logs out of the current NEONKUBE context"),
		Long:    logoutLong,
		Example: logoutExample,
		Run: func(cmd *cobra.Command, args []string) {

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "logout")

			neon_utility.ExecNeonCli(neonCliArgs)
		},
	}

	return cmd
}
