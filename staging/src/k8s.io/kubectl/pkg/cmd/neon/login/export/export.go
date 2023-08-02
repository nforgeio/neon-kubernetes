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

package neon_login_export

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	neon_utility "k8s.io/kubectl/pkg/cmd/neon"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	exportLong = templates.LongDesc(i18n.T(`
		Exports the current context or a context by name to STDOUT or or file.`))

	exportExample = templates.Examples(i18n.T(`
		# Export the current NEONKUBE context to STDOUT
		neon login export
		
		# Export the current NEONKUBE context to a file
		neon login export file.yaml
		
		# Export a specific NEONKUBE context to a file
		neon login export --context=root@MY-CLUSTER file.yaml`))
)

type flags struct {
	context string
}

// NewCmdNeonLoginExport returns a Command instance for NEON-CLI 'login export' sub command
func NewCmdNeonLoginExport(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {

	flags := flags{}

	cmd := &cobra.Command{
		Use:     "export [PATH]",
		Short:   i18n.T("Exports the current context or a context by name to STDOUT or a file"),
		Long:    exportLong,
		Example: exportExample,
		Run: func(cmd *cobra.Command, args []string) {

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "login")
			neonCliArgs = append(neonCliArgs, "export")

			if flags.context != "" {
				neonCliArgs = append(neonCliArgs, "--context="+flags.context)
			}

			if len(args) > 0 {
				neonCliArgs = append(neonCliArgs, args[0])
			}

			neon_utility.ExecNeonCli(neonCliArgs)
		},
	}

	cmd.Flags().StringVarP(&flags.context, "context", "", "",
		i18n.T("Identifies the context to export by name"))

	return cmd
}
