// Copyright 2018 Red Hat
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
	"os"

	"github.com/coreos/coreos-assembler/mantle/platform"
	"github.com/coreos/pkg/capnslog"
	"github.com/spf13/pflag"
	kubecli "kubevirt.io/client-go/kubecli"
)

var (
	plog = capnslog.NewPackageLogger("github.com/coreos/coreos-assembler/mantle", "platform/api/kubevirt")
)

type Options struct {
	*platform.Options

	// Config file. Defaults to ~/.kube/config
	ConfigPath string
	Image      string
}

type ServiceListen struct {
	Name        string
	BindAddress string
	Port        int
}

type API struct {
	options        *Options
	kubeVirtClient kubecli.KubevirtClient
	namespace      string
}

// func (a *API) GetConsoleOutput(name string) (string, error) {
// 	out, err := a.compute.Instances.GetSerialPortOutput(a.options.Project, a.options.Zone, name).Do()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to retrieve console output for %q: %v", name, err)
// 	}
// 	return out.Contents, nil
// }

// func (a *API) TerminateInstance(name string) error {
// 	plog.Debugf("Terminating instance %q", name)

// 	_, err := a.compute.Instances.Delete(a.options.Project, a.options.Zone, name).Do()
// 	return err
// }

// New creates a new KubeVirt API wrapper.
// No validation is done that kubeconfig exist and before using the API a
// preflight check is recommended via api.PreflightCheck
func New(opts *Options) (*API, error) {
	if opts.ConfigPath == "" {
		opts.ConfigPath = "~/.kube/config"
	}
	fmt.Printf("The configpath is %s\n", opts.ConfigPath)

	os.Setenv("KUBECONFIG", opts.ConfigPath)
	// kubecli.DefaultClientConfig() prepares config using kubeconfig.
	// typically, you need to set env variable, KUBECONFIG=<path-to-kubeconfig>/.kubeconfig
	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})

	// // retrive default namespace.
	namespace, _, err := clientConfig.Namespace()
	if err != nil {
		return nil, fmt.Errorf("error in namespace : %v\n", err)
	}

	// get the kubevirt client, using which kubevirt resources can be managed.
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot obtain KubeVirt client: %v\n", err)
	}

	// var vmi *v1.VirtualMachineInstance
	// vmi, err = virtClient.VirtualMachineInstance(namespace).Create(context.Background(), vmi)
	// if err != nil {
	// 	return nil, err
	// }
	// // Fetch list of VMs & VMIs
	// vmList, err := virtClient.VirtualMachine(namespace).List(&metav1.ListOptions{})
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot obtain KubeVirt vm list: %v\n", err)
	// }
	// vmiList, err := virtClient.VirtualMachineInstance(namespace).List(context.Background(), &metav1.ListOptions{})
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot obtain KubeVirt vmi list: %v\n", err)
	// }

	// w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	// fmt.Fprintln(w, "Type\tName\tNamespace\tStatus")

	// for _, vm := range vmList.Items {
	// 	fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", vm.Kind, vm.Name, vm.Namespace, vm.Status.Ready)
	// }
	// for _, vmi := range vmiList.Items {
	// 	fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", vmi.Kind, vmi.Name, vmi.Namespace, vmi.Status.Phase)
	// }
	// w.Flush()

	api := &API{
		kubeVirtClient: virtClient,
		options:        opts,
		namespace:      namespace,
	}
	// client, _ := kubecli.GetKubevirtClientFromFlags(server.URL()+proxyPath, "")

	// vmi := api.NewMinimalVMI("testvm")
	// server.AppendHandlers(ghttp.CombineHandlers(
	// 	ghttp.VerifyRequest("POST", path.Join(proxyPath, basePath)),
	// 	ghttp.RespondWithJSONEncoded(http.StatusCreated, vmi),
	// ))
	// createdVMI, _ := client.VirtualMachineInstance(k8sv1.NamespaceDefault).Create(context.Background(), vmi)
	// app.VMI = createdVMI

	return api, nil
}
