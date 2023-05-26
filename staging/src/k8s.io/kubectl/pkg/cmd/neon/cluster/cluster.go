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
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"

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

var (
	clusterLong = templates.LongDesc(i18n.T(`
	    Use the subcommands to deploy and manage NEONKUBE clusters.`))

	clusterExample = templates.Examples(i18n.T(`
		# Performs some cluster checks on the current NEONKUBE cluster,
		# typically performed by maintainers while testing
		neon cluster check

		# Lists the avaiable NONKUBE cluster dashboards and displays
		# the kubernetes and grafana dashboards
		neon cluster dashboard
		neon cluster dashboard kubernetes
		neon cluster dashboard grafana

		# Removes the current NEONKUBE cluster
		neon cluster delete

		# Deploys a NEONKUBE cluster from a cluster definition YAML file
		neon cluster deploy my-cluster.yaml
		
		# Returns information about the health of the current NEONKUBE cluster
		neon cluster health

		# Returns information about the current NEONKUBE cluster
		neon cluster info

		# Determines whether the current NEONKUBE cluster is locked, with exit
		# codes: 0=locked, 1=request failed, 2=unlocked.  Dangerous operations
		# like cluster delete are disabled by default for locked clusters.
		neon cluster islocked

		# Locks the current NEONKUBE cluster
		neon cluster lock

		# Pauses the current NEONKUBE cluster.  This isn't available for all
		# hosting environments.
		neon cluster pause

		# Prepares the underlying infrastructure for a NEONKUBE cluster and then
		# setup the cluster.  This approach is usually used only by maintainers
		# because the [neon cluster deploy] command performs both steps.
		neon cluster prepare my-cluster.yaml
		neon cluster setup root@CLUSTERNAME
		
		# Used to retrieve or update the purpose for the current NEONKUBE cluster.
		# We currently define these purposes: development, production, stage, test
		# and unspecified.  The first command returns the current NEONKUBE cluster
		# purpose and the second changes the purpose to [production].
		neon cluster purpose
		neon cluster purpose production

		# Used to reset the state of the current NEONKUBE cluster to close to its
		# "factory new" condition.  This typically used for returning a cluster to
		# a known state for testing.  The cluster must be unlocked.
		neon cluster reset

		# Restart a paused or stopped current NEONKUBE cluster
		neon cluster start

		# Stops the current NEONKUBE cluster by shutting down all of its node machines.
		# The cluster must be unlocked.
		neon cluster stop
		
		# Unlocks the current NEONKUBE cluster
		neon cluster unlock

		# Verifies that a cluster definition YAML file is valid
		neon cluster validate my-cluster.yaml`))
)

// NewCmdNeonCluster returns a Command instance for NEON-CLI 'cluster' sub commands.
func NewCmdNeonCluster(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cluster SUBCOMMAND",
		Short:   i18n.T("Deploy and manage NEONKUBE clusters."),
		Long:    clusterLong,
		Example: clusterExample,
		Run:     cmdutil.DefaultSubCommandRun(streams.Out),
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
