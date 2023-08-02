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

package neon_cluster_prepare

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
	prepareLong = templates.LongDesc(i18n.T(`
		MAINTAINERS ONLY: Provisions local and/or cloud infrastructure required to host a 
		NEONKUBE cluster.  This includes provisioning networks, load balancers, virtual
		machines, etc.  Once the infrastructure is ready, you'll use the [neon cluster setup ...]
		command to actually setup the cluster.
		
		NOTE: This is used by maintainers while debugging cluster setup.`))

	prepareExample = templates.Examples(i18n.T(`
	# Step 1: Prepare cluster infrastructure
	neon cluster prepare my-cluster.yaml
	
	# Step 2: Setup the cluster
	neon cluster setup root@CLUSTERNAME`))
)

type flags struct {
	baseImageName  string
	debug          bool
	disablePending bool
	insecure       bool
	maxParallel    int
	nodeImagePath  string
	nodeImageUri   string
	noTelemetry    bool
	packageCache   string
	quiet          bool
	unredacted     bool
	useStaged      string
}

// NewCmdNeonClusterPrepare returns a Command instance for NEON-CLI 'cluster prepare' sub command
func NewCmdNeonClusterPrepare(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {

	flags := flags{}

	cmd := &cobra.Command{
		Use:     "prepare CLUSTERDEF",
		Short:   i18n.T("MAINTAINERS ONLY: Provisions local and/or cloud infrastructure required to host a NEONKUBE cluster"),
		Long:    prepareLong,
		Example: prepareExample,
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				neon_utility.CommandError("CLUSTERDEF argument is required")
			}

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "cluster")
			neonCliArgs = append(neonCliArgs, "prepare")
			neonCliArgs = append(neonCliArgs, args[0])

			if flags.baseImageName != "" {
				neonCliArgs = append(neonCliArgs, "--base-image-name="+flags.baseImageName)
			}
			if flags.debug {
				neonCliArgs = append(neonCliArgs, "--debug")
			}
			if flags.disablePending {
				neonCliArgs = append(neonCliArgs, "--disable-pending")
			}
			if flags.insecure {
				neonCliArgs = append(neonCliArgs, "--insecure")
			}
			if flags.maxParallel != neon_utility.DefaultClusterDeployParallel {
				neonCliArgs = append(neonCliArgs, "--max-parallel="+strconv.Itoa(flags.maxParallel))
			}
			if flags.nodeImagePath != "" {
				neonCliArgs = append(neonCliArgs, "--node-image-path="+flags.nodeImagePath)
			}
			if flags.nodeImageUri != "" {
				neonCliArgs = append(neonCliArgs, "--node-image-uri="+flags.nodeImageUri)
			}
			if flags.noTelemetry {
				neonCliArgs = append(neonCliArgs, "--no-telemetry")
			}
			if flags.packageCache != "" {
				neonCliArgs = append(neonCliArgs, "--package-cache"+flags.packageCache)
			}
			if flags.quiet {
				neonCliArgs = append(neonCliArgs, "--quiet")
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

	cmd.Flags().StringVarP(&flags.baseImageName, "base-image-name", "", "",
		templates.LongDesc(i18n.T(`
			Specifies the base image name to use when operating
			in [--debug] mode.  This will be the name of the base
			image file as published to our public S3 bucket for
			the target hosting manager`)))

	cmd.Flags().BoolVarP(&flags.debug, "debug", "", false,
		templates.LongDesc(i18n.T(`
			Implements cluster setup from the base rather than the node image.
			This mode is useful while developing and debugging cluster setup.
			
			NOTE: This mode is not supported for cloud and bare-metal environments.`)))

	cmd.Flags().BoolVarP(&flags.disablePending, "disable-pending", "", false,
		templates.LongDesc(i18n.T(`
			Disable parallization of setup tasks across steps.
			This is generally intended for use while debugging
			cluster setup and may slow cluster prepare substantially.`)))

	cmd.Flags().BoolVarP(&flags.insecure, "insecure", "", false,
		templates.LongDesc(i18n.T(`
		MAINTAINER ONLY: Prevents the cluster node [sysadmin]
		account from being set to a secure password and also
		enables SSH password authentication.  Used for debugging.

		WARNING: NEVER USE FOR PRODUCTION CLUSTERS!`)))

	cmd.Flags().IntVarP(&flags.maxParallel, "max-parallel", "", neon_utility.DefaultClusterDeployParallel,
		templates.LongDesc(i18n.T(`
			Specifies the maximum number of node related operations
			to perform in parallel.`)))

	cmd.Flags().StringVarP(&flags.nodeImagePath, "node-image-path", "", "",
		templates.LongDesc(i18n.T(`
			Uses the node image at the PATH specified rather than
			downloading the node image.  This is useful for debugging 
			node image changes locally.`)))

	cmd.Flags().StringVarP(&flags.nodeImageUri, "node-image-uri", "", "",
		templates.LongDesc(i18n.T(`
			Overrides the default node image URI.  Note that this is ignored
			for [--debug] mode and when [--node-image-path] is specified.`)))

	cmd.Flags().BoolVarP(&flags.noTelemetry, "no-telemetry", "", false,
		templates.LongDesc(i18n.T(`
			Disables telemetry uploads for failed cluster deployment,
			overriding the NEONKUBE_DISABLE_TELEMETRY environment
			variable`)))

	cmd.Flags().StringVarP(&flags.packageCache, "package-cache", "", "",
		i18n.T("Specifies one or more APT Package cache servers by hostname and port for use by the new cluster.  Specify multiple servers by separating the endpoints with commas."))

	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "", false,
		i18n.T("Only print the currently executing step rather than displaying detailed setup status"))

	cmd.Flags().BoolVarP(&flags.unredacted, "unredacted", "", false,
		templates.LongDesc(i18n.T(`
			Runs commands with potential secrets without redacting logs.
			This is useful for debugging cluster setup issues.  Do not 
			use for production clusters.`)))

	cmd.Flags().StringVarP(&flags.useStaged, "use-staged", "", "",
		templates.LongDesc(i18n.T("Deploy a NEONKUBE cluster from an internal build, optionally specifiying a GitHub source branch")))
	cmd.Flag("use-staged").NoOptDefVal = neon_utility.NoFlagValue

	return cmd
}
