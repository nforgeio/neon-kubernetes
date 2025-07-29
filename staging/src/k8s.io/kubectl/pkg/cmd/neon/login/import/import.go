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

package neon_login_import

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	neon_utility "k8s.io/kubectl/pkg/cmd/neon"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	importLong = templates.LongDesc(i18n.T(`
		Imports a NEONKUBE cluster context.`))

	importExample = templates.Examples(i18n.T(`
		# Import a NEONKUBE context and then log into it
		neon login import my-context.yaml
		
		# Import a NEONKUBE context but don't log into it
		neon login import --no-login my-context.yaml
		
		# Import a NEONKUBE context and don't prompt for permission to overwrite an existing context
		neon login import --force my-context.yaml`))
)

type flags struct {
	force   bool
	noLogin bool
}

// NewCmdNeonLoginImport returns a Command instance for NEON-CLI 'login import' sub command
func NewCmdNeonLoginImport(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {

	flags := flags{}

	cmd := &cobra.Command{
		Use:     "import PATH",
		Short:   i18n.T("Imports a NEONKUBE cluster context"),
		Long:    importLong,
		Example: importExample,
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				neon_utility.CommandError("PATH argument is required")
			}

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "login")
			neonCliArgs = append(neonCliArgs, "import")
			neonCliArgs = append(neonCliArgs, args[0])

			if flags.force {
				neonCliArgs = append(neonCliArgs, "--force")
			}
			if flags.noLogin {
				neonCliArgs = append(neonCliArgs, "--no-login")
			}

			neon_utility.ExecNeonCli(neonCliArgs)
		},
	}

	cmd.Flags().BoolVarP(&flags.force, "force", "", false,
		i18n.T("Don't prompt for permission to replace an existing context"))

	cmd.Flags().BoolVarP(&flags.noLogin, "no-login", "", false,
		i18n.T("Don't login to the new context"))

	return cmd
}
