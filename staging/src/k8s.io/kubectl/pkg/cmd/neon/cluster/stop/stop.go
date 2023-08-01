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

package neon_cluster_stop

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	neon_utility "k8s.io/kubectl/pkg/cmd/neon"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	stopLong = templates.LongDesc(i18n.T(`
		Stops the current NEONKUBE cluster.
		
		NOTE: This command requires that the cluster be unlocked by default as
        a safety measure.

		All clusters besides NEONDESKTOP clusters are locked by default when
		they're deployed.  You can disable this by setting [IsLocked=false]
		in your cluster definition or by executing this command:

    	neon cluster unlock`))

	stopExample = templates.Examples(i18n.T(`
		# Stop the current NEONKUBE cluster
		neon cluster stop
		
		# Stop the current NEONKUBE cluster without prompting for permission
		# if it is currently locked
		neon cluster stop --force
		
		# Stops the cluster immediately without waiting for a graceful shutdown. 
		# This may cause data loss.
		neon cluster stop --turnoff --force`))
)

type flags struct {
	force   bool
	turnoff bool
}

// NewCmdNeonClusterStopreturns a Command instance for NEON-CLI 'cluster stop' sub command
func NewCmdNeonClusterStop(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {

	flags := flags{}

	cmd := &cobra.Command{
		Use:                   "stop",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Stops the current NEONKUBE cluster."),
		Long:                  stopLong,
		Example:               stopExample,
		Run: func(cmd *cobra.Command, args []string) {

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "cluster")
			neonCliArgs = append(neonCliArgs, "stop")

			if flags.force {
				neonCliArgs = append(neonCliArgs, "--force")
			}
			if flags.turnoff {
				neonCliArgs = append(neonCliArgs, "--turnoff")
			}

			neon_utility.ExecNeonCli(neonCliArgs)
		},
	}

	cmd.Flags().BoolVarP(&flags.force, "force", "", false,
		i18n.T("Don't prompt for permission or require the the cluster be unlocked before stopping"))
	cmd.Flags().BoolVarP(&flags.turnoff, "turnoff", "", false,
		i18n.T("Turns the nodes off immediately without waiting for a graceful shutdown.  This may cause data loss."))

	return cmd
}
