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
	neon_utility "k8s.io/kubectl/pkg/cmd/neon"
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
	    Manage NEONKUBE cluster logins/contexts.
		
		The base command is used to select the current Kubernetes context, which identifies
		the cluster where subsequent [neon] tool commands will be applied.  Subcommands can
		be used to delete, export, import, or list Kubernetes contexts.
		
		The base login command accepts an optional CONTEXTREF parameter.  This can be a local
		context name like "root@mycluster" or a URI for the cluster itself like "https://CLUSTERDOMAIN".
		For context names, the command simply switches to the new context in the kubeconfig
		file.
		
		For cluster URIs, the command uses Single Sign On (SSO) to sign into to the cluster 
		itself (via a browser window) to obtain security information, adds a new context to 
		the kubeconfig  file and then makes it current.  SSO is a convienent way to log into
		a cluster you didn't deploy or wasn't deployed from the corrent workstation.

		You can also use the base login command to select a new default namespace for the
		current cluster.  This is a bit easier to use than the equivalent [kubectl] commands`))

	loginExample = templates.Examples(i18n.T(`
	    # Print the current context name and default namespace
		neon login
		neon login --output json
		neon login -o yaml

		# Select root@mycluster as the current context
		neon login root@mycluster

		# Select root@mycluster as the current context and also set "my-namespace" as the 
		# default namespace in one go
		neon login root@mycluster --namespace=my-namespace
		neon login root@mycluster -n=my-namespace

		# Use single sign on to authenticate with the cluster, creating a local context and
		# make it the current context.  Note that the presence of the "https://" scheme in
		# the argument is required.  Note that newly created clusters are identified by a 
		# UUID but users are free to configure their own DNS to make this nicer.
		neon login https://09e2-bc8a-216c-f17d.neoncluster.io

		# Alternative way to use use single sign.  This accepts the "--sso" option with just
		# the cluster domain but otherwise works the same.
		neon login --sso 09e2-bc8a-216c-f17d.neoncluster.io

		# Change the current context's default namespace
		neon login --namespace=my-namespace
		neon login -n=my-namespace

		# Removes the current NEONKUBE context from the kubeconfig
		neon login delete

		# Exports the current NEONKUBE cluster context to a file
		neon login export

		# Imports an exported NEONKUBE cluster context from a file
		neon login import

		# Lists the NEONKUBE cluster contexts from kubeconfig
		neon login list`))
)

type flags struct {
	namespace    string
	outputFormat string
	sso          string
}

// NewCmdNeonCluster returns a Command instance for NEON-CLI 'cluster' sub commands.
func NewCmdNeonLogin(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {

	flags := flags{}

	cmd := &cobra.Command{
		Use:     "login ([CONTEXTREF])|(SUBCOMMAND)",
		Short:   i18n.T("Manage NEONKUBE cluster logins/contexts"),
		Long:    loginLong,
		Example: loginExample,
		Run: func(cmd *cobra.Command, args []string) {

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "login")

			if flags.namespace != "" {
				neonCliArgs = append(neonCliArgs, "--namespace="+flags.namespace)
			}

			if flags.outputFormat != "" {
				neonCliArgs = append(neonCliArgs, "--output="+flags.outputFormat)
			}

			if flags.sso != "" {
				neonCliArgs = append(neonCliArgs, "--sso="+flags.sso)
			}

			if len(args) > 0 {
				neonCliArgs = append(neonCliArgs, args[0])
			}

			neon_utility.ExecNeonCli(neonCliArgs)
		},
	}

	cmd.Flags().StringVarP(&flags.namespace, "namespace", "n", "",
		i18n.T("Identifies the namespace to be configured as the default for the current context"))

	cmd.Flags().StringVarP(&flags.outputFormat, "output", "o", "",
		i18n.T("specifies the format used to print the current context and namespace (json|yaml)"))

	cmd.Flags().StringVarP(&flags.sso, "sso", "", "",
		i18n.T("Alternative way to specify the cluster domain for single sign on"))

	// subcommands
	cmd.AddCommand(neon_login_delete.NewCmdNeonLoginDelete(f, streams))
	cmd.AddCommand(neon_login_export.NewCmdNeonLoginExport(f, streams))
	cmd.AddCommand(neon_login_import.NewCmdNeonLoginImport(f, streams))
	cmd.AddCommand(neon_login_list.NewCmdNeonLoginList(f, streams))

	return cmd
}
