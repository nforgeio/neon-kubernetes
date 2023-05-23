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

	neonclustercheck "k8s.io/kubectl/pkg/cmd/neon/cluster/check"
	neonclusterdashboard "k8s.io/kubectl/pkg/cmd/neon/cluster/dashboard"
	neonclusterdelete "k8s.io/kubectl/pkg/cmd/neon/cluster/delete"
	neonclusterdeploy "k8s.io/kubectl/pkg/cmd/neon/cluster/deploy"
	neonclusterhealth "k8s.io/kubectl/pkg/cmd/neon/cluster/health"
	neonclusterinfo "k8s.io/kubectl/pkg/cmd/neon/cluster/info"
	neonclusterislocked "k8s.io/kubectl/pkg/cmd/neon/cluster/islocked"
	neonclusterpause "k8s.io/kubectl/pkg/cmd/neon/cluster/pause"
	neonclusterprepare "k8s.io/kubectl/pkg/cmd/neon/cluster/prepare"
	neonclusterpurpose "k8s.io/kubectl/pkg/cmd/neon/cluster/purpose"
	neonclusterreset "k8s.io/kubectl/pkg/cmd/neon/cluster/reset"
	neonclustersetup "k8s.io/kubectl/pkg/cmd/neon/cluster/setup"
	neonclusterstart "k8s.io/kubectl/pkg/cmd/neon/cluster/start"
	neonclusterstop "k8s.io/kubectl/pkg/cmd/neon/cluster/stop"
	neonclustervalidate "k8s.io/kubectl/pkg/cmd/neon/cluster/validate"
)

var commandExample string = `
# Provision a 
# neon cluster prepare my-clusterdef.yaml 
# neon cluster setup root@my-cluster
`

// NewCmdNeonCluster returns a Command instance for NEON-CLI 'cluster' commands.
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
	cmd.AddCommand(neonclustercheck.NewCmdNeonClusterCheck(f, streams))
	cmd.AddCommand(neonclusterdashboard.NewCmdNeonClusterDashboard(f, streams))
	cmd.AddCommand(neonclusterdelete.NewCmdNeonClusterDelete(f, streams))
	cmd.AddCommand(neonclusterdeploy.NewCmdNeonClusterDeploy(f, streams))
	cmd.AddCommand(neonclusterhealth.NewCmdNeonClusterHealth(f, streams))
	cmd.AddCommand(neonclusterinfo.NewCmdNeonClusterInfo(f, streams))
	cmd.AddCommand(neonclusterislocked.NewCmdNeonClusterIsLocked(f, streams))
	cmd.AddCommand(neonclusterpause.NewCmdNeonClusterPause(f, streams))
	cmd.AddCommand(neonclusterprepare.NewCmdNeonClusterPrepare(f, streams))
	cmd.AddCommand(neonclusterpurpose.NewCmdNeonClusterPurpose(f, streams))
	cmd.AddCommand(neonclusterreset.NewCmdNeonClusterReset(f, streams))
	cmd.AddCommand(neonclustersetup.NewCmdNeonClusterSetup(f, streams))
	cmd.AddCommand(neonclusterstart.NewCmdNeonClusterStart(f, streams))
	cmd.AddCommand(neonclusterstop.NewCmdNeonClusterStop(f, streams))
	cmd.AddCommand(neonclustervalidate.NewCmdNeonClusterValidate(f, streams))

	return cmd
}
