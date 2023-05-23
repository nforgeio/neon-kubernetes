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

package neon_helm

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// NewCmdNeonHelm returns a Command instance for NEON-CLI 'helm' sub command
func NewCmdNeonHelm(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "",
		DisableFlagsInUseLine: true,
		Short:                 "",
		Long:                  "",
		Example:               "",
		Run:                   cmdutil.DefaultSubCommandRun(streams.Out),
	}

	return cmd
}
