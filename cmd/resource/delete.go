package resource

import (
	"fmt"

	"github.com/crossplane/hiveworld/pkg/api"
	"github.com/crossplane/hiveworld/pkg/client"
	"github.com/crossplane/hiveworld/pkg/resource"
	"github.com/crossplane/provider-terraform-plugin/pkg/registry"
)

// ReadResource will read a resource description
func DeleteResource(resourceReadPath string, provider *client.Provider, r *registry.Registry) error {
	rd, err := resource.ResourceDataFromFile(resourceReadPath)
	if err != nil {
		return err
	}
	res, err := rd.ManagedResource(r)
	if err != nil {
		return err
	}
	gvk := rd.GVK
	err = api.Delete(provider, r, res, gvk)
	if err != nil {
		return err
	}
	fmt.Println("Resource successfully deleted!")
	return nil
}
