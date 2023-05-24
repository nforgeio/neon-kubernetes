/*
Copyright 2023 NEONFORGE, LLC.

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

package neon_cluster

import (
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	neon_cluster_check "k8s.io/kubectl/pkg/cmd/neon/cluster/check"
	neon_cluster_dashboard "k8s.io/kubectl/pkg/cmd/neon/cluster/dashboard"
	neon_cluster_delete "k8s.io/kubectl/pkg/cmd/neon/cluster/delete"
	neon_cluster_deploy "k8s.io/kubectl/pkg/cmd/neon/cluster/deploy"
	neon_cluster_health "k8s.io/kubectl/pkg/cmd/neon/cluster/health"
	neon_cluster_info "k8s.io/kubectl/pkg/cmd/neon/cluster/info"
	neon_cluster_islocked "k8s.io/kubectl/pkg/cmd/neon/cluster/islocked"
	neon_cluster_lock "k8s.io/kubectl/pkg/cmd/neon/cluster/lock"
	neon_cluster_pause "k8s.io/kubectl/pkg/cmd/neon/cluster/pause"
	neon_cluster_prepare "k8s.io/kubectl/pkg/cmd/neon/cluster/prepare"
	neon_cluster_purpose "k8s.io/kubectl/pkg/cmd/neon/cluster/purpose"
	neon_cluster_reset "k8s.io/kubectl/pkg/cmd/neon/cluster/reset"
	neon_cluster_setup "k8s.io/kubectl/pkg/cmd/neon/cluster/setup"
	neon_cluster_start "k8s.io/kubectl/pkg/cmd/neon/cluster/start"
	neon_cluster_stop "k8s.io/kubectl/pkg/cmd/neon/cluster/stop"
	neon_cluster_validate "k8s.io/kubectl/pkg/cmd/neon/cluster/validate"
)

var commandExample string = `
# Provision a 
# neon cluster prepare my-clusterdef.yaml 
# neon cluster setup root@my-cluster
`

// NewCmdNeonCluster returns a Command instance for NEON-CLI 'cluster' sub commands.
func NewCmdNeonCluster(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "cluster SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Deploy and manage NEONKUBE clusters.",
		Long:                  "Deploy and manage NEONKUBE clusters.",
		Example:               "",
		Run:                   cmdutil.DefaultSubCommandRun(streams.Out),
	}
	// subcommands
	cmd.AddCommand(neon_cluster_check.NewCmdNeonClusterCheck(f, streams))
	cmd.AddCommand(neon_cluster_dashboard.NewCmdNeonClusterDashboard(f, streams))
	cmd.AddCommand(neon_cluster_delete.NewCmdNeonClusterDelete(f, streams))
	cmd.AddCommand(neon_cluster_deploy.NewCmdNeonClusterDeploy(f, streams))
	cmd.AddCommand(neon_cluster_health.NewCmdNeonClusterHealth(f, streams))
	cmd.AddCommand(neon_cluster_info.NewCmdNeonClusterInfo(f, streams))
	cmd.AddCommand(neon_cluster_islocked.NewCmdNeonClusterIsLocked(f, streams))
	cmd.AddCommand(neon_cluster_lock.NewCmdNeonClusterLock(f, streams))
	cmd.AddCommand(neon_cluster_pause.NewCmdNeonClusterPause(f, streams))
	cmd.AddCommand(neon_cluster_prepare.NewCmdNeonClusterPrepare(f, streams))
	cmd.AddCommand(neon_cluster_purpose.NewCmdNeonClusterPurpose(f, streams))
	cmd.AddCommand(neon_cluster_reset.NewCmdNeonClusterReset(f, streams))
	cmd.AddCommand(neon_cluster_setup.NewCmdNeonClusterSetup(f, streams))
	cmd.AddCommand(neon_cluster_start.NewCmdNeonClusterStart(f, streams))
	cmd.AddCommand(neon_cluster_stop.NewCmdNeonClusterStop(f, streams))
	cmd.AddCommand(neon_cluster_validate.NewCmdNeonClusterValidate(f, streams))

	return cmd
}
