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

package neon_login_list

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	neon_utility "k8s.io/kubectl/pkg/cmd/neon"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	listLong = templates.LongDesc(i18n.T(`
		Lists the NEONKUBE contexts.`))

	listExample = templates.Examples(i18n.T(`
		# List the NEONKUBE contexts
		neon login list
`))
)

type flags struct {
	outputFormat string
}

// NewCmdNeonLoginList returns a Command instance for NEON-CLI 'login list' sub command
func NewCmdNeonLoginList(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {

	flags := flags{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   i18n.T("Lists the NEONKUBE contexts"),
		Long:    listLong,
		Example: listExample,
		Run: func(cmd *cobra.Command, args []string) {

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "login")
			neonCliArgs = append(neonCliArgs, "list")

			if flags.outputFormat != "" {
				neonCliArgs = append(neonCliArgs, "--output="+flags.outputFormat)
			}

			neon_utility.ExecNeonCli(neonCliArgs)
		},
	}

	cmd.Flags().StringVarP(&flags.outputFormat, "output", "o", "",
		i18n.T("specifies the output format (json|yaml)"))

	return cmd
}
