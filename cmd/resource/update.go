package resource

import (
	"fmt"

	"github.com/crossplane/hiveworld/pkg/resource"
	"github.com/crossplane/provider-terraform-plugin/pkg/api"
	"github.com/crossplane/provider-terraform-plugin/pkg/client"
	"github.com/crossplane/provider-terraform-plugin/pkg/registry"
)

// UpdateResource will read a resource description from disk and sync with the provider
func UpdateResource(resourceReadPath string, provider *client.Provider, r *registry.Registry) error {
	rd, err := resource.ResourceDataFromFile(resourceReadPath)
	if err != nil {
		return err
	}
	res, err := rd.ManagedResource(r)
	if err != nil {
		return err
	}
	gvk := rd.GVK
	newRes, err := api.Update(provider, r, res)
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
