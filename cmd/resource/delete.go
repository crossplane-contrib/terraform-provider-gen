package resource

import (
	"fmt"

	"github.com/crossplane/hiveworld/pkg/resource"
	"github.com/crossplane/terraform-provider-runtime/pkg/api"
	"github.com/crossplane/terraform-provider-runtime/pkg/client"
	"github.com/crossplane/terraform-provider-runtime/pkg/registry"
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
	err = api.Delete(provider, r, res)
	if err != nil {
		return err
	}
	fmt.Println("Resource successfully deleted!")
	return nil
}
