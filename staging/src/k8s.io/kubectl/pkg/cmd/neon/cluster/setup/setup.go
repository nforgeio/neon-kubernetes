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

package neon_cluster_setup

import (
	"strconv"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	neon_utility "k8s.io/kubectl/pkg/cmd/neon"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	setupLong = templates.LongDesc(i18n.T(`
		MAINTAINERS ONLY: Sets up a NEONKUBE cluster as described in the cluster
		definition file.  This is the second part of deploying a cluster in two
		stages, where you first prepare the cluster to provision any virtual
		machines and network infrastructure and then you setup NEONKUBE.
		
		NOTE: This is used by maintainers while debugging cluster setup.`))

	setupExample = templates.Examples(i18n.T(`
		# Step 1: Prepare cluster infrastructure
		neon cluster prepare my-cluster.yaml
		
		# Step 2: Setup the cluster
		neon cluster setup root@CLUSTERNAME`))
)

type flags struct {
	check          bool
	debug          bool
	disablePending bool
	force          bool
	maxParallel    int
	noTelemetry    bool
	quiet          bool
	unredacted     bool
	uploadCharts   bool
	useStaged      string
}

// NewCmdNeonClusterSetup returns a Command instance for NEON-CLI 'cluster setup' sub command
func NewCmdNeonClusterSetup(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {

	flags := flags{}

	cmd := &cobra.Command{
		Use:     "setup CONTEXTNAME",
		Short:   i18n.T("MAINTAINERS ONLY: Sets up a NEONKUBE cluster on prepared infrastructure"),
		Long:    setupLong,
		Example: setupExample,
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				neon_utility.CommandError("CONTEXTNAME argument is required")
			}

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "cluster")
			neonCliArgs = append(neonCliArgs, "setup")
			neonCliArgs = append(neonCliArgs, args[0])

			if flags.check {
				neonCliArgs = append(neonCliArgs, "--check")
			}
			if flags.debug {
				neonCliArgs = append(neonCliArgs, "--debug")
			}
			if flags.disablePending {
				neonCliArgs = append(neonCliArgs, "--disable-pending")
			}
			if flags.force {
				neonCliArgs = append(neonCliArgs, "--force")
			}
			if flags.maxParallel != neon_utility.DefaultClusterDeployParallel {
				neonCliArgs = append(neonCliArgs, "--max-parallel="+strconv.Itoa(flags.maxParallel))
			}
			if flags.noTelemetry {
				neonCliArgs = append(neonCliArgs, "--no-telemetry")
			}
			if flags.quiet {
				neonCliArgs = append(neonCliArgs, "--quiet")
			}
			if flags.uploadCharts {
				neonCliArgs = append(neonCliArgs, "--upload-charts")
			}
			if flags.unredacted {
				neonCliArgs = append(neonCliArgs, "--unredacted")
			}
			if flags.useStaged != "" {
				if flags.useStaged == neon_utility.NoFlagValue {
					neonCliArgs = append(neonCliArgs, "--use-staged")
				} else {
					neonCliArgs = append(neonCliArgs, "--use-staged="+flags.useStaged)
				}
			}

			neon_utility.ExecNeonCli(neonCliArgs)
		},
	}

	cmd.Flags().BoolVarP(&flags.check, "check", "", false,
		templates.LongDesc(i18n.T(`
			Performs development related checks against the cluster
			after it's been setup.  Note that checking is disabled
			when [--debug] is specified.

			NOTE: A non-zero exit code will be returned when this
		    option is specified and one or more checks fail`)))

	cmd.Flags().BoolVarP(&flags.debug, "debug", "", false,
		templates.LongDesc(i18n.T(`
			Implements cluster setup from the base rather than the node image.
			This mode is useful while developing and debugging cluster setup.

			NOTE: This is not supported for cloud and bare-metal environments.`)))

	cmd.Flags().BoolVarP(&flags.disablePending, "disable-pending", "", false,
		templates.LongDesc(i18n.T(`
			Disable parallization of setup tasks across steps.
			This is generally intended for use while debugging
			cluster setup and may slow cluster setup substantially.`)))

	cmd.Flags().BoolVarP(&flags.force, "force", "", false,
		templates.LongDesc(i18n.T(`
			Don't prompt before removing existing contexts
			that reference the target cluster`)))

	cmd.Flags().IntVarP(&flags.maxParallel, "max-parallel", "", neon_utility.DefaultClusterDeployParallel,
		templates.LongDesc(i18n.T(`
			Specifies the maximum number of node related operations
			to perform in parallel.`)))

	cmd.Flags().BoolVarP(&flags.noTelemetry, "no-telemetry", "", false,
		templates.LongDesc(i18n.T(`
			Disables telemetry uploads for failed cluster deployment,
			overriding the NEONKUBE_DISABLE_TELEMETRY environment
			variable`)))

	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "", false,
		i18n.T("Only print the currently executing step rather than displaying detailed setup status"))

	cmd.Flags().BoolVarP(&flags.unredacted, "unredacted", "", false,
		templates.LongDesc(i18n.T(`
			Runs commands with potential secrets without redacting logs.
			This is useful for debugging cluster setup issues.  Do not 
			use for production clusters.`)))

	cmd.Flags().BoolVarP(&flags.uploadCharts, "upload-charts", "", false,
		templates.LongDesc(i18n.T(`
			Upload Helm charts from your workstation rather than using
			the charts baked into the node image`)))

	cmd.Flags().StringVarP(&flags.useStaged, "use-staged", "", "",
		templates.LongDesc(i18n.T(`
			Deploy a NEONKUBE cluster from an internal build, optionally specifiying a GitHub source branch"`)))
	cmd.Flag("use-staged").NoOptDefVal = neon_utility.NoFlagValue

	return cmd
}
