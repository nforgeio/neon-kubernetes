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

package neon_login

import (
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"

	neon_login_delete "k8s.io/kubectl/pkg/cmd/neon/login/delete"
	neon_login_export "k8s.io/kubectl/pkg/cmd/neon/login/export"
	neon_login_import "k8s.io/kubectl/pkg/cmd/neon/login/import"
	neon_login_list "k8s.io/kubectl/pkg/cmd/neon/login/list"
)

var (
	loginLong = templates.LongDesc(i18n.T(`
	    Use the subcommands to manage cluster logins.`))

	loginExample = templates.Examples(i18n.T(`
		# Removes the current NEONKUBE context from the kubeconfig
		neon login delete

		# Exports the current NEONKUBE cluster context to a file
		neon login export

		# Imports an exported NEONKUBE cluster context from a file
		neon login import

		# Lists the NEONKUBE cluster contexts from the kubeconfig
		neon login list`))
)

// NewCmdNeonCluster returns a Command instance for NEON-CLI 'cluster' sub commands.
func NewCmdNeonLogin(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "login SUBCOMMAND",
		Short:   i18n.T("Deploy and manage NEONKUBE clusters."),
		Long:    loginLong,
		Example: loginExample,
		Run:     cmdutil.DefaultSubCommandRun(streams.Out),
	}
	// subcommands
	cmd.AddCommand(neon_login_delete.NewCmdNeonLoginDelete(f, streams))
	cmd.AddCommand(neon_login_export.NewCmdNeonLoginExport(f, streams))
	cmd.AddCommand(neon_login_import.NewCmdNeonLoginImport(f, streams))
	cmd.AddCommand(neon_login_list.NewCmdNeonLoginList(f, streams))

	return cmd
}
