package resource

import (
	"io/ioutil"
	"strings"

	xpresource "github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/provider-terraform-plugin/pkg/registry"
	"gopkg.in/yaml.v2"
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
)

type ResourceData struct {
	Data []byte
	GVK  k8schema.GroupVersionKind
}

func (rd *ResourceData) ManagedResource(r *registry.Registry) (xpresource.Managed, error) {
	rep, err := r.GetResourceUnmarshalFunc(rd.GVK)
	if err != nil {
		return nil, err
	}
	rep = registry.ResourceUnmarshalFunc(rep)
	return rep(rd.Data)
}

func ResourceDataFromFile(path string) (*ResourceData, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	gvk, err := ResourceDataGVK(content)
	return &ResourceData{Data: content, GVK: gvk}, err
}

func ResourceDataGVK(data []byte) (k8schema.GroupVersionKind, error) {
	gvk := k8schema.GroupVersionKind{}
	vk := versionKind{}
	err := yaml.Unmarshal(data, &vk)
	if err != nil {
		return gvk, err
	}
	gvk.Group = vk.Group()
	gvk.Version = vk.Version()
	gvk.Kind = vk.Kind

	return gvk, nil
}

// GVK respresents the outer "header" fields of a CustomResource
// to enable parsing so the GVK can be used to lookup corresponding types.
type versionKind struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}

// Group returns the "group" portion of the `APIVersion` field
// eg returns "iam.gcp.terraform-plugin.crossplane.io" for APIVersion="iam.gcp.terraform-plugin.crossplane.io/v1alpha1"
func (v versionKind) Group() string {
	parts := strings.Split(v.APIVersion, "/")
	return parts[0]
}

// Version returns the "version" portion of the `APIVersion` field
// eg returns "v1alpha1" for APIVersion="iam.gcp.terraform-plugin.crossplane.io/v1alpha1"
func (v versionKind) Version() string {
	parts := strings.Split(v.APIVersion, "/")
	return parts[1]
}

func (v versionKind) GVK() k8schema.GroupVersionKind {
	return k8schema.GroupVersionKind{
		Group:   v.Group(),
		Version: v.Version(),
		Kind:    v.Kind,
	}
}
