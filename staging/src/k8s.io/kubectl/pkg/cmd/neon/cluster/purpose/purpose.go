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

package neon_cluster_purpose

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	neon_utility "k8s.io/kubectl/pkg/cmd/neon"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	purposeLong = templates.LongDesc(i18n.T(`
		Prints or sets information about the cluster purpose.  This
		is a high-level indication of how the cluster will be used.
		These purposes are currently supported:
		
		development production stage test unspecified`))

	purposeExample = templates.Examples(i18n.T(`
		# Print the purpose for the current NEONKUBE cluster
		neon cluster purpose
		
		# Set the current NEONKUBE's cluster purpose to: production
		neon cluster purpose production`))
)

// NewCmdNeonClusterPurpose returns a Command instance for NEON-CLI 'cluster purpose' sub command
func NewCmdNeonClusterPurpose(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "purpose [NEWPURPOSE]",
		Short:   i18n.T("Prints or sets information about the cluster purpose"),
		Long:    purposeLong,
		Example: purposeExample,
		Run: func(cmd *cobra.Command, args []string) {

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "cluster")
			neonCliArgs = append(neonCliArgs, "purpose")

			if len(args) > 0 {
				neonCliArgs = append(neonCliArgs, args[0])
			}

			neon_utility.ExecNeonCli(neonCliArgs)
		},
	}

	return cmd
}
