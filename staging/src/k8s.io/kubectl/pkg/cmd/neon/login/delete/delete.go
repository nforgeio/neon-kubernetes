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

package neon_login_delete

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	neon_utility "k8s.io/kubectl/pkg/cmd/neon"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	deleteLong = templates.LongDesc(i18n.T(`
		Removes the current NEONKUBE cluster context or a context by name.`))

	deleteExample = templates.Examples(i18n.T(`
		# Remove the current NEONKUBE cluster context
		neon login delete
		
		# Remove the root@mycluster context
		neon login delete root@mycluster

		# Remove root@mycluster context without prompting for permission
		neon login delete --force root@mycluster`))
)

type flags struct {
	force bool
}

// NewCmdNeonLoginDelete returns a Command instance for NEON-CLI 'login delete' sub command
func NewCmdNeonLoginDelete(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {

	flags := flags{}

	cmd := &cobra.Command{
		Use:     "delete [CONTEXTNAME]",
		Short:   i18n.T("Removes the current NEONKUBE cluster context or a context by name"),
		Long:    deleteLong,
		Example: deleteExample,
		Run: func(cmd *cobra.Command, args []string) {

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "login")
			neonCliArgs = append(neonCliArgs, "delete")

			if len(args) > 0 {
				neonCliArgs = append(neonCliArgs, args[0])
			}

			if flags.force {
				neonCliArgs = append(neonCliArgs, "--force")
			}

			neon_utility.ExecNeonCli(neonCliArgs)
		},
	}

	cmd.Flags().BoolVarP(&flags.force, "force", "", false,
		i18n.T("Don't prompt for permission and also ignore missing contexts"))

	return cmd
}
