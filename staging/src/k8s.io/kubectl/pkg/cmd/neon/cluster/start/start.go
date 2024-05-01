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

package neon_cluster_start

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	neon_utility "k8s.io/kubectl/pkg/cmd/neon"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	startLong = templates.LongDesc(i18n.T(`
	Starts the current stopped or paused NEONKUBE cluster.`))

	startExample = templates.Examples(i18n.T(`
		# Start the current NEONKUBE cluster
		neon cluster start
`))
)

// NewCmdNeonClusterStart returns a Command instance for NEON-CLI 'cluster start' sub command
func NewCmdNeonClusterStart(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   i18n.T("Starts the current stopped or paused NEONKUBE cluster"),
		Long:    startLong,
		Example: startExample,
		Run: func(cmd *cobra.Command, args []string) {

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "cluster")
			neonCliArgs = append(neonCliArgs, "start")

			neon_utility.ExecNeonCli(neonCliArgs)
		},
	}

	return cmd
}
