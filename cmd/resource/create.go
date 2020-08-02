package resource

import (
	"fmt"

	"github.com/crossplane/hiveworld/pkg/resource"
	"github.com/crossplane/terraform-provider-runtime/pkg/api"
	"github.com/crossplane/terraform-provider-runtime/pkg/client"
	"github.com/crossplane/terraform-provider-runtime/pkg/registry"
)

// ReadResource will read a resource description
func CreateResource(resourceReadPath string, provider *client.Provider, r *registry.Registry) error {
	rd, err := resource.ResourceDataFromFile(resourceReadPath)
	if err != nil {
		return err
	}
	res, err := rd.ManagedResource(r)
	if err != nil {
		return err
	}
	gvk := rd.GVK
	newRes, err := api.Create(provider, r, res)
	if err != nil {
		return err
	}
	asYAML, err := r.GetYAMLEncodeFunc(gvk)
	if err != nil {
		return err
	}
	yamlBytes, err := asYAML(newRes)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(yamlBytes))

	return nil
}
