package resource

import (
	"fmt"

	"github.com/crossplane/hiveworld/pkg/api"
	"github.com/crossplane/hiveworld/pkg/client"
	"github.com/crossplane/hiveworld/pkg/registry"
	"github.com/crossplane/hiveworld/pkg/resource"
)

// ReadResource will read a resource description
func CreateResource(resourceReadPath string, provider *client.Provider) error {
	rd, err := resource.ResourceDataFromFile(resourceReadPath)
	if err != nil {
		return err
	}
	res, err := rd.ManagedResource()
	if err != nil {
		return err
	}
	gvk := rd.GVK
	newRes, err := api.Create(provider, res, gvk)
	if err != nil {
		return err
	}
	asYAML, err := registry.GetYAMLEncodeFunc(gvk)
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
