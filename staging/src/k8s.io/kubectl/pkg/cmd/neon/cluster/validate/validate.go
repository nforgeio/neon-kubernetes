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

package neon_cluster_validate

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	neon_utility "k8s.io/kubectl/pkg/cmd/neon"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	validateLong = templates.LongDesc(i18n.T(`
		Validates a NEONKUBE cluster definition YAML file.`))

	validateExample = templates.Examples(i18n.T(`
		# Validate a NEONKUBE cluster definition file
		neon cluster validate my-cluster.yaml`))
)

// NewCmdNeonClusterValidate returns a Command instance for NEON-CLI 'cluster validate' sub command
func NewCmdNeonClusterValidate(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validate CLUSTERDEF",
		Short:   i18n.T("Validates a NEONKUBE cluster definition YAML file"),
		Long:    validateLong,
		Example: validateExample,
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				neon_utility.CommandError("CLUSTERDEF argument is required")
			}

			neonCliArgs := make([]string, 0)
			neonCliArgs = append(neonCliArgs, "cluster")
			neonCliArgs = append(neonCliArgs, "validate")
			neonCliArgs = append(neonCliArgs, args[0])

			neon_utility.ExecNeonCli(neonCliArgs)
		},
	}

	return cmd
}
