package registry

import (
	"fmt"

	xpresource "github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/hashicorp/terraform/providers"
	"github.com/zclconf/go-cty/cty"
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
)

type ResourceUnmarshalFunc func([]byte) (xpresource.Managed, error)
type CtyEncodeFunc func(xpresource.Managed, *providers.Schema) (cty.Value, error)
type CtyDecodeFunc func(xpresource.Managed, cty.Value, *providers.Schema) (xpresource.Managed, error)
type YAMLEncodeFunc func(xpresource.Managed) ([]byte, error)

var (
	resourceRepresenterMap = make(map[k8schema.GroupVersionKind]ResourceUnmarshalFunc)
	ctyEncodeFuncMap       = make(map[k8schema.GroupVersionKind]CtyEncodeFunc)
	ctyDecodeFuncMap       = make(map[k8schema.GroupVersionKind]CtyDecodeFunc)
	terraformNameToGVK     = make(map[string]k8schema.GroupVersionKind)
	gvkToTerraformName     = make(map[k8schema.GroupVersionKind]string)
	yamlEncodeFuncMap      = make(map[k8schema.GroupVersionKind]YAMLEncodeFunc)
)

func RegisterYAMLEncodeFunc(gvk k8schema.GroupVersionKind, f YAMLEncodeFunc) {
	if gvk.String() == "" {
		panic("RegisterYAMLEncodeFunc called with uninitialized GroupVersionKind")
	}
	if f == nil {
		panic(fmt.Sprintf("Cannot initialize RegisterYAMLEncodeFunc called with nil value for gvk=%s", gvk.String()))
	}
	yamlEncodeFuncMap[gvk] = f
}

func RegisterResourceUnmarshalFunc(gvk k8schema.GroupVersionKind, f ResourceUnmarshalFunc) {
	if gvk.String() == "" {
		panic("RegisterResourceUnmarshalFunc called with uninitialized GroupVersionKind")
	}
	if f == nil {
		panic(fmt.Sprintf("Cannot initialize RegisterResourceUnmarshalFunc called with nil value for gvk=%s", gvk.String()))
	}
	resourceRepresenterMap[gvk] = f
}

func RegisterCtyEncodeFunc(gvk k8schema.GroupVersionKind, f CtyEncodeFunc) {
	if gvk.String() == "" {
		panic("RegisterCtyEncodeFunc called with uninitialized GroupVersionKind")
	}
	if f == nil {
		panic(fmt.Sprintf("Cannot initialize: RegisterCtyEncodeFunc called with nil value for gvk=%s", gvk.String()))
	}
	ctyEncodeFuncMap[gvk] = f
}

func RegisterCtyDecodeFunc(gvk k8schema.GroupVersionKind, f CtyDecodeFunc) {
	if gvk.String() == "" {
		panic("RegisterCtyDecodeFunc called with uninitialized GroupVersionKind")
	}
	if f == nil {
		panic(fmt.Sprintf("Cannot initialize: RegisterCtyDecodeFunc called with nil value for gvk=%s", gvk.String()))
	}
	ctyDecodeFuncMap[gvk] = f
}

func RegisterTerraformNameMapping(tfname string, gvk k8schema.GroupVersionKind) {
	if gvk.String() == "" {
		panic("RegisterTerraformNameMapping called with uninitialized GroupVersionKind")
	}
	if tfname == "" {
		panic("RegisterTerraformNameMapping called with uninitialized tfname")
	}
	terraformNameToGVK[tfname] = gvk
	gvkToTerraformName[gvk] = tfname
}

func GetYAMLEncodeFunc(gvk k8schema.GroupVersionKind) (YAMLEncodeFunc, error) {
	f, ok := yamlEncodeFuncMap[gvk]
	if !ok {
		return nil, fmt.Errorf("Could not find a yaml encoder function for GVK=%s", gvk.String())
	}
	return f, nil
}

func GetCtyEncoder(gvk k8schema.GroupVersionKind) (CtyEncodeFunc, error) {
	f, ok := ctyEncodeFuncMap[gvk]
	if !ok {
		return nil, fmt.Errorf("Could not find a cty encoder function for GVK=%s", gvk.String())
	}
	return f, nil
}

func GetCtyDecoder(gvk k8schema.GroupVersionKind) (CtyDecodeFunc, error) {
	f, ok := ctyDecodeFuncMap[gvk]
	if !ok {
		return nil, fmt.Errorf("Could not find a cty decoder function for GVK=%s", gvk.String())
	}
	return f, nil
}

func GetResourceUnmarshalFunc(gvk k8schema.GroupVersionKind) (ResourceUnmarshalFunc, error) {
	rep, ok := resourceRepresenterMap[gvk]
	if !ok {
		return nil, fmt.Errorf("Could not find a resource representer for GVK=%s", gvk.String())
	}
	return rep, nil
}

func GetGVKForTerraformName(name string) (k8schema.GroupVersionKind, error) {
	gvk, ok := terraformNameToGVK[name]
	if !ok {
		return gvk, fmt.Errorf("Could not find GVK for Terraform resource name=%s", name)
	}
	return gvk, nil
}

func GetTerraformNameForGVK(gvk k8schema.GroupVersionKind) (string, error) {
	name, ok := gvkToTerraformName[gvk]
	if !ok {
		return "", fmt.Errorf("Could not find GVK for Terraform resource name=%s", name)
	}
	return name, nil
}
