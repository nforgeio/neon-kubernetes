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

package neon_cluster_deploy

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
	deployLong = templates.LongDesc(i18n.T(`
		Deploys a NEONKUBE cluster, based on a cluster definition YAML file`))

	deployExample = templates.Examples(i18n.T(`
		# Deploys a NEONKUBE cluster
		neon cluster deploy my-cluster.yaml`))
)

type flags struct {
	check        bool
	force        bool
	maxParallel  int
	noTelemetry  bool
	packageCache string
	quiet        bool
	uploadCharts bool
	useStaged    string
	usePreview   bool
}

// NewCmdNeonClusterDeploy returns a Command instance for NEON-CLI 'cluster delete' sub command
func NewCmdNeonClusterDeploy(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {

	flags := flags{}

	cmd := &cobra.Command{
		Use:     "deploy CLUSTERDEF",
		Short:   i18n.T("Deploys a NEONKUBE cluster, based on a cluster definition"),
		Long:    deployLong,
		Example: deployExample,
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				neon_utility.CommandError("CLUSTERDEF argument is required")
			}

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "cluster")
			neonCliArgs = append(neonCliArgs, "deploy")
			neonCliArgs = append(neonCliArgs, args[0])

			if flags.check {
				neonCliArgs = append(neonCliArgs, "--check")
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
			if flags.packageCache != "" {
				neonCliArgs = append(neonCliArgs, "--package-cache="+flags.packageCache)
			}
			if flags.quiet {
				neonCliArgs = append(neonCliArgs, "--quiet")
			}
			if flags.uploadCharts {
				neonCliArgs = append(neonCliArgs, "--upload-charts")
			}
			if flags.usePreview {
				neonCliArgs = append(neonCliArgs, "--use-preview")
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
		i18n.T("Performs development related checks against the cluster after it's been deployed.  A non-zero exit code will be returned when this option is specified and one or more checks fail"))
	cmd.Flags().BoolVarP(&flags.force, "force", "", false,
		i18n.T("Don't prompt for permission to remove existing contexts that reference the target cluster"))
	cmd.Flags().IntVarP(&flags.maxParallel, "max-parallel", "", neon_utility.DefaultClusterDeployParallel,
		i18n.T("Specifies the maximum number of node related operations to perform in parallel"))
	cmd.Flags().BoolVarP(&flags.noTelemetry, "no-telemetry", "", false,
		i18n.T("Disables telemetry uploads for failed cluster deployment, overriding the NEONKUBE_DISABLE_TELEMETRY environment variable"))
	cmd.Flags().StringVarP(&flags.packageCache, "package-cache", "", "",
		i18n.T("Specifies one or more APT Package cache servers by hostname and port for use by the new cluster.  Specify multiple servers by separating the endpoints with commas"))
	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "", false,
		i18n.T("Only print the currently executing step rather than displaying detailed setup status"))
	cmd.Flags().BoolVarP(&flags.uploadCharts, "upload-charts", "", false,
		i18n.T("MAINTAINER ONLY: Upload Helm charts from your workstation rather than using the charts baked into the node image"))
	cmd.Flags().StringVarP(&flags.useStaged, "use-staged", "", "",
		i18n.T("MAINTAINER ONLY: Deploy a NEONKUBE cluster from an internal build, optionally specifiying a GitHub source branch"))
	cmd.Flag("use-staged").NoOptDefVal = neon_utility.NoFlagValue
	cmd.Flags().BoolVarP(&flags.usePreview, "use-preview", "", false,
		templates.LongDesc(i18n.T(`
			Uses the preview VM image from the Azure Marketplace to
			provision the cluster.  Note that a project maintainer
			must specifically grant access to each user, generally
			to work with users to solve NEONKUBE problems.`)))

	return cmd
}
