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

package neon_cluster_delete

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
		Permanently deletes the current NEONKUBE cluster.`))

	deleteExample = templates.Examples(i18n.T(`
		# Delete the current NEONKUBE cluster
		neon cluster delete
		
		# Delete the current NEONKUBE cluster without prompting for permission
		neon cluster delete --force
		
		# Delete the NEONKUBE cluster identified by kube context name
		neon cluster delete root@my-cluster`))
)

type flags struct {
	force bool
}

// NewCmdNeonClusterDelete returns a Command instance for NEON-CLI 'cluster delete' sub command
func NewCmdNeonClusterDelete(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {

	flags := flags{}

	cmd := &cobra.Command{
		Use:     "delete CONTEXT",
		Short:   i18n.T("Permenantly deletes a NEONKUBE cluster"),
		Long:    deleteLong,
		Example: deleteExample,
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				neon_utility.CommandError("CONTEXT argument is required")
			}

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "cluster")
			neonCliArgs = append(neonCliArgs, "delete")
			neonCliArgs = append(neonCliArgs, args[0])

			if flags.force {
				neonCliArgs = append(neonCliArgs, "--force")
			}

			neon_utility.ExecNeonCli(neonCliArgs)
		},
	}

	cmd.Flags().BoolVarP(&flags.force, "force", "", false,
		i18n.T("Don't prompt for permission or require the the cluster be unlocked before removal"))

	return cmd
}
