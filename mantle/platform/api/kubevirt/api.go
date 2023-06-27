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
	"context"
	"net/http"
	"path"

	"github.com/coreos/pkg/capnslog"

	"github.com/coreos/coreos-assembler/mantle/platform"
	"github.com/onsi/gomega/ghttp"
	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/api"
	"kubevirt.io/client-go/kubecli"
)

var (
	plog = capnslog.NewPackageLogger("github.com/coreos/coreos-assembler/mantle", "platform/api/kubevirt")
)

type Options struct {
	ConfigPath string
	*platform.Options
}

type ServiceListen struct {
	Name        string
	BindAddress string
	Port        int
}

type API struct {
	VMI VirtualMachineInstance
}

type VirtualMachineInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// VirtualMachineInstance Spec contains the VirtualMachineInstance specification.
	Spec VirtualMachineInstanceSpec `json:"spec" valid:"required"`
}

func New(opts *Options) (*API, error) {
	var server *ghttp.Server
	basePath := "/apis/kubevirt.io/v1/namespaces/default/virtualmachineinstances"
	proxyPath := "/proxy/path"

	app := &API{}
	client, err := kubecli.GetKubevirtClientFromFlags(server.URL()+proxyPath, "")

	vmi := api.NewMinimalVMI("testvm")
	server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("POST", path.Join(proxyPath, basePath)),
		ghttp.RespondWithJSONEncoded(http.StatusCreated, vmi),
	))
	createdVMI, err := client.VirtualMachineInstance(k8sv1.NamespaceDefault).Create(context.Background(), vmi)
	app.VMI = createdVMI

	return app, nil
}

func NewMinimalVMI(name string) *v1.VirtualMachineInstance {
	return NewMinimalVMIWithNS(k8sv1.NamespaceDefault, name)
}

// This is meant for testing
func NewMinimalVMIWithNS(namespace, name string) *v1.VirtualMachineInstance {
	vmi := v1.NewVMIReferenceFromNameWithNS(namespace, name)
	vmi.Spec = v1.VirtualMachineInstanceSpec{Domain: v1.DomainSpec{}}
	vmi.Spec.Domain.Resources.Requests = k8sv1.ResourceList{
		k8sv1.ResourceMemory: resource.MustParse("8192Ki"),
	}
	vmi.TypeMeta = metav1.TypeMeta{
		APIVersion: v1.GroupVersion.String(),
		Kind:       "VirtualMachineInstance",
	}
	return vmi
}

// VirtualMachineInstanceSpec is a description of a VirtualMachineInstance.
type VirtualMachineInstanceSpec struct {
	// Specification of the desired behavior of the VirtualMachineInstance on the host.
	Domain DomainSpec `json:"domain"`
	// List of volumes that can be mounted by disks belonging to the vmi.
	Volumes []Volume `json:"volumes,omitempty"`
}

// DomainSpec represents the actual conversion to libvirt XML. The fields must be
// tagged, and they must correspond to the libvirt domain as described in
// https://libvirt.org/formatdomain.html.
type DomainSpec struct {
	Devices  Devices   `xml:"devices"`
	Resource *Resource `xml:"resource,omitempty"`
}

type Devices struct {
	Disks []Disk `json:"disks,omitempty"`
	// Whether to have random number generator from host
	// +optional
	Rng *Rng `json:"rng,omitempty"`
}

type Resource struct {
	Partition string `xml:"partition"`
}

type Memory struct {
	Value uint64 `xml:",chardata"`
	Unit  string `xml:"unit,attr"`
}

type Disk struct {
	// Name is the device name
	Name string `json:"name"`
}

// Rng represents the random device passed from host
type Rng struct {
}

// Volume represents a named volume in a vmi.
type Volume struct {
	// Volume's name.
	// Must be a DNS_LABEL and unique within the vmi.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name"`
	// VolumeSource represents the location and type of the mounted volume.
	// Defaults to Disk, if no type is specified.
	VolumeSource `json:",inline"`
}

// Represents the source of a volume to mount.
// Only one of its members may be specified.
type VolumeSource struct {
	// SecretVolumeSource represents a reference to a secret data in the same namespace.
	// More info: https://kubernetes.io/docs/concepts/configuration/secret/
	// +optional
	Secret *SecretVolumeSource `json:"secret,omitempty"`
}

// SecretVolumeSource adapts a Secret into a volume.
type SecretVolumeSource struct {
	// Name of the secret in the pod's namespace to use.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
	SecretName string `json:"secretName,omitempty"`
	// Specify whether the Secret or it's keys must be defined
	// +optional
	Optional *bool `json:"optional,omitempty"`
	// The volume label of the resulting disk inside the VMI.
	// Different bootstrapping mechanisms require different values.
	// Typical values are "cidata" (cloud-init), "config-2" (cloud-init) or "OEMDRV" (kickstart).
	// +optional
	VolumeLabel string `json:"volumeLabel,omitempty"`
}
