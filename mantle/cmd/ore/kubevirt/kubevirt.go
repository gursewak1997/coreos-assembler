// Copyright 2017 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubevirt

import (
	"fmt"

	"github.com/coreos/coreos-assembler/mantle/cli"
	"github.com/coreos/coreos-assembler/mantle/platform"
	"github.com/coreos/coreos-assembler/mantle/platform/api/kubevirt"
	"github.com/coreos/pkg/capnslog"
	"github.com/spf13/cobra"
)

var (
	plog = capnslog.NewPackageLogger("github.com/coreos/coreos-assembler/mantle", "ore/kubevirt")

	KubeVirt = &cobra.Command{
		Use:   "kubevirt [command]",
		Short: "kubevirt image utilities",
	}
	API     *kubevirt.API
	options kubevirt.Options
)

func init() {
	KubeVirt.PersistentFlags().StringVar(&options.ConfigPath, "kubeconfig", "~/.kube/config", "KubeVirt config file")
	cli.WrapPreRun(KubeVirt, preflightCheck)
}

func preflightCheck(cmd *cobra.Command, args []string) error {
	plog.Debugf("Running KubeVirt preflight check")
	api, err := kubevirt.New(&options)
	if err != nil {
		return fmt.Errorf("could not create KubeVirt client: %v", err)
	}
	// if err := api.PreflightCheck(); err != nil {
	// 	return fmt.Errorf("could not complete KubeVirt preflight check: %v", err)
	// }

	plog.Debugf("Preflight check success; we have liftoff")
	API = api
	return nil
}
