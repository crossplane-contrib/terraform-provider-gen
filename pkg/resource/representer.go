package resource

import (
	"io/ioutil"
	"strings"

	"github.com/crossplane/hiveworld/pkg/client"
	"gopkg.in/yaml.v2"
)

// YAMLByteRepresenter is a transformer that represents a resource as a yaml-encoded []byte
type YAMLByteRepresenter struct {
	Raw []byte
	gvk GVK
}

// Representer describes a type that can create different representations of the same
// underlying Resource.
type Representer interface {
	AsYAML() ([]byte, error)
	//AsGRPC([]byte) (string, error)
	//AsManagedResource([]byte) (xpresource.Managed, error)
}

// GVK respresents the outer "header" fields of a CustomResource
// to enable parsing so the GVK can be used to lookup corresponding types.
type GVK struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}

// Group returns the "group" portion of the `APIVersion` field
// eg returns "iam.gcp.crossplane.io" for APIVersion="iam.gcp.crossplane.io/v1alpha1"
func (g GVK) Group() string {
	parts := strings.Split(g.APIVersion, "/")
	return parts[0]
}

// Version returns the "version" portion of the `APIVersion` field
// eg returns "v1alpha1" for APIVersion="iam.gcp.crossplane.io/v1alpha1"
func (g GVK) Version() string {
	parts := strings.Split(g.APIVersion, "/")
	return parts[1]
}

// RepresenterFromYAMLFile uses GVK to partially parse
// the CustomResource in order to look up its type for complete unmarshalling.
// ParseResourceFromFile assumes path is fully-qualified and points at a file
// containing the yaml representation of a managed resource type that is
// understood by the provider.
func RepresenterFromYAMLFile(path string, p *client.Provider) (Representer, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewYAMLByteRepresenter(content)
}

// NewYAMLByteRepresenter creates a type satisfying the Representer interface
// which can understand a yaml []byte representation of a managed resource
func NewYAMLByteRepresenter(content []byte) (Representer, error) {
	y := &YAMLByteRepresenter{
		Raw: content,
	}
	gvk := GVK{}
	err := yaml.Unmarshal(y.Raw, &gvk)
	if err != nil {
		return nil, err
	}
	y.gvk = gvk

	return y, nil
}

// AsYAML returns a stringified representation of the raw file comments
func (y *YAMLByteRepresenter) AsYAML() ([]byte, error) {
	return y.Raw, nil
}
