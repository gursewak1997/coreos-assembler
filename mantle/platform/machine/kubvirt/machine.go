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
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/coreos/coreos-assembler/mantle/platform"
)

type machine struct {
	cluster *cluster
	// mach    *kubevirt.VirtualMachine
	name    string
	intIP   string
	extIP   string
	dir     string
	journal *platform.Journal
	console string
}

func (km *machine) ID() string {
	return km.name
}

func (km *machine) IP() string {
	return km.extIP
}

func (km *machine) PrivateIP() string {
	return km.intIP
}

func (km *machine) RuntimeConf() platform.RuntimeConfig {
	return km.cluster.RuntimeConf()
}

func (km *machine) SSHClient() (*ssh.Client, error) {
	return km.cluster.SSHClient(km.IP())
}

func (km *machine) PasswordSSHClient(user string, password string) (*ssh.Client, error) {
	return km.cluster.PasswordSSHClient(km.IP(), user, password)
}

func (km *machine) SSH(cmd string) ([]byte, []byte, error) {
	return km.cluster.SSH(km, cmd)
}

func (km *machine) IgnitionError() error {
	return nil
}

func (km *machine) Start() error {
	return platform.StartMachine(km, km.journal)
}

func (km *machine) Reboot() error {
	return platform.RebootMachine(km, km.journal)
}

func (km *machine) WaitForReboot(timeout time.Duration, oldBootId string) error {
	return platform.WaitForMachineReboot(km, km.journal, timeout, oldBootId)
}

func (km *machine) Destroy() {
	if err := km.saveConsole(); err != nil {
		plog.Errorf("Error saving console for instance %v: %v", km.ID(), err)
	}

	if err := km.cluster.flight.api.TerminateInstance(km.name); err != nil {
		plog.Errorf("Error terminating instance %v: %v", km.ID(), err)
	}

	if km.journal != nil {
		km.journal.Destroy()
	}

	km.cluster.DelMach(km)
}

func (km *machine) ConsoleOutput() string {
	return km.console
}

func (km *machine) saveConsole() error {
	var err error
	km.console, err = km.cluster.flight.api.GetConsoleOutput(km.name)
	if err != nil {
		return err
	}

	path := filepath.Join(km.dir, "console.txt")
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(km.console); err != nil {
		return err
	}

	return nil
}

func (km *machine) JournalOutput() string {
	if km.journal == nil {
		return ""
	}

	data, err := km.journal.Read()
	if err != nil {
		plog.Errorf("Reading journal for instance %v: %v", km.ID(), err)
	}
	return string(data)
}
