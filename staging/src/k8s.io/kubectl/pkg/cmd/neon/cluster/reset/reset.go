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

package neon_cluster_reset

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	neon_utility "k8s.io/kubectl/pkg/cmd/neon"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	resetLong = templates.LongDesc(i18n.T(`
		Resets the current cluster to its factory new condition.
		
		This command works by removing all non-system namespaces including [default],
		along with any resources contained within.  The [default] namespace will be
		recreated afterwards, restoring it to its original empty condition.  You can
		specify namespaces to be retained via [--namespace-exclude], passing a comma
		separated list of namespaces.
		
		The command also resets Harbor, Minio, CRIO-O, Dex and the monitoring components
		to their defaults.  All components are reset by default, but you can wou can control
		whether some components are reset by passing one or more of the options.
		
		This command will not work on a locked clusters as a safety measure.

		All clusters besides NEONDESKTOP clusters are locked by default when they're
		deployed.  You can disable this by setting [IsLocked=false] in your cluster
		definition or by executing this command on your cluster:

    	neon cluster unlock`))

	resetExample = templates.Examples(i18n.T(`
		# Reset cluster with confirmation prompt
		neon cluster reset

		# Reset cluster without confirmation prompt
		neon cluster reset --force

		# Full cluster reset while retaining the ""foo"" and ""bar"" namespaces
		neon cluster reset --keep-namespaces=foo,bar

		# Full cluster reset excluding all non-system namespaces
		neon cluster reset --keep-namespaces=*

		# Reset Minio and Harbor as well as removing all non-system namespaces:
		neon cluster reset --minio --harbor`))
)

type flags struct {
	auth           bool
	crio           bool
	force          bool
	harbor         bool
	keepNamespaces string
	minio          bool
	monitoring     bool
}

// NewCmdNeonClusterReset returns a Command instance for NEON-CLI 'cluster reset' sub command
func NewCmdNeonClusterReset(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {

	flags := flags{}

	cmd := &cobra.Command{
		Use:     "reset",
		Short:   i18n.T("Resets the current NEONKUBE cluster to its factor new condition"),
		Long:    resetLong,
		Example: resetExample,
		Run: func(cmd *cobra.Command, args []string) {

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "cluster")
			neonCliArgs = append(neonCliArgs, "reset")

			if flags.auth {
				neonCliArgs = append(neonCliArgs, "--auth")
			}
			if flags.crio {
				neonCliArgs = append(neonCliArgs, "--crio")
			}
			if flags.force {
				neonCliArgs = append(neonCliArgs, "--force")
			}
			if flags.harbor {
				neonCliArgs = append(neonCliArgs, "--harbor")
			}
			if flags.keepNamespaces != "" {
				neonCliArgs = append(neonCliArgs, "--keep-namespaces="+flags.keepNamespaces)
			}
			if flags.minio {
				neonCliArgs = append(neonCliArgs, "--minio")
			}
			if flags.monitoring {
				neonCliArgs = append(neonCliArgs, "--monitoring")
			}

			neon_utility.ExecNeonCli(neonCliArgs)
		},
	}

	cmd.Flags().BoolVarP(&flags.auth, "auth", "", false,
		i18n.T("Resets authentication (Dex, Glauth)"))
	cmd.Flags().BoolVarP(&flags.crio, "crio", "", false,
		i18n.T("Resets referenced container registeries to the cluster definition specifications and removes any non-system container images"))
	cmd.Flags().BoolVarP(&flags.force, "force", "", false,
		i18n.T("Don't prompt for permission or require the cluster be unlocked before reset"))
	cmd.Flags().BoolVarP(&flags.harbor, "harbor", "", false,
		i18n.T("Resets Harbor components"))
	cmd.Flags().StringVarP(&flags.keepNamespaces, "keep-namespaces", "", "",
		i18n.T("comma separated list of non-system namespaces to be retained or \"*\" to retain all non-system namespaces"))
	cmd.Flags().BoolVarP(&flags.minio, "minio", "", false,
		i18n.T("Rresets Minio"))
	cmd.Flags().BoolVarP(&flags.monitoring, "monitoring", "", false,
		i18n.T("Clears monitoring data as well as non-system dashboards and alerts"))

	return cmd
}
