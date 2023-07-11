// Copyright 2015 CoreOS, Inc.
// Copyright 2015 The Go Authors.
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
	"crypto/rand"
	"fmt"

	"github.com/coreos/coreos-assembler/mantle/platform"
	"github.com/coreos/coreos-assembler/mantle/platform/conf"
)

type cluster struct {
	*platform.BaseCluster
	flight *flight
}

func (ac *cluster) vmname() string {
	b := make([]byte, 5)
	rand.Read(b)
	return fmt.Sprintf("%s-%x", ac.Name()[0:13], b)
}

func (kc *cluster) NewMachine(userdata *conf.UserData) (platform.Machine, error) {
	return kc.NewMachineWithOptions(userdata, platform.MachineOptions{})
}

func (kc *cluster) NewMachineWithOptions(userdata *conf.UserData, options platform.MachineOptions) (platform.Machine, error) {
	instance := kc.flight.api.NewMinimalVMI(kc.vmname())

	mach := &machine{
		cluster: kc,
		mach:    instance,
	}

	return mach, nil
}

func (kc *cluster) Destroy() {
	kc.BaseCluster.Destroy()
	kc.flight.DelCluster(kc)
}
