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

package neon_cluster_check

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	neon_utility "k8s.io/kubectl/pkg/cmd/neon"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	checkLong = templates.LongDesc(i18n.T(`
		Performs some checks against the current NEOKUBE cluster  This is 
		typically used by maintainer when debugging cluster deployment.`))

	checkExample = templates.Examples(i18n.T(`
		# Perform all checks
		neon cluster check
		neon cluster check --all

		# Perform selected checks
		neon cluster check --container-images --details`))
)

type flags struct {
	all             bool
	containerImages bool
	priorityClass   bool
	resources       bool
	details         bool
}

// NewCmdNeonClusterCheck returns a Command instance for NEON-CLI 'cluster check' sub command
func NewCmdNeonClusterCheck(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {

	flags := flags{}

	cmd := &cobra.Command{
		Use:     "check",
		Short:   i18n.T("Peforms some checks against the current NEONKUBE cluster"),
		Long:    checkLong,
		Example: checkExample,
		Run: func(cmd *cobra.Command, args []string) {

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "cluster")
			neonCliArgs = append(neonCliArgs, "check")

			if flags.all {
				neonCliArgs = append(neonCliArgs, "--all")
			}
			if flags.containerImages {
				neonCliArgs = append(neonCliArgs, "--container-images")
			}
			if flags.priorityClass {
				neonCliArgs = append(neonCliArgs, "--priority-class")
			}
			if flags.resources {
				neonCliArgs = append(neonCliArgs, "--resources")
			}
			if flags.details {
				neonCliArgs = append(neonCliArgs, "--details")
			}

			neon_utility.ExecNeonCli(neonCliArgs)
		},
	}

	cmd.Flags().BoolVarP(&flags.all, "all", "", false,
		i18n.T("Performs all checks (this is implied when no other options are present)"))
	cmd.Flags().BoolVarP(&flags.containerImages, "container-images", "", false,
		i18n.T("Verifies that all container images running on the cluster are included in the cluster manifest"))
	cmd.Flags().BoolVarP(&flags.priorityClass, "priority-class", "", false,
		i18n.T("Verifies that all running pods have a non-zero PriorityClass"))
	cmd.Flags().BoolVarP(&flags.resources, "resources", "", false,
		i18n.T("Verifies that all pod containers specify resource request and limits"))
	cmd.Flags().BoolVarP(&flags.details, "details", "", false,
		i18n.T("VIncludes additional information for some of the checks even when there are no errors"))

	return cmd
}
