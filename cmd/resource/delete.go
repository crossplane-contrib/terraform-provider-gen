package resource

import (
	"fmt"

	"github.com/crossplane/hiveworld/pkg/api"
	"github.com/crossplane/hiveworld/pkg/client"
	"github.com/crossplane/hiveworld/pkg/resource"
)

// ReadResource will read a resource description
func DeleteResource(resourceReadPath string, provider *client.Provider) error {
	rd, err := resource.ResourceDataFromFile(resourceReadPath)
	if err != nil {
		return err
	}
	res, err := rd.ManagedResource()
	if err != nil {
		return err
	}
	gvk := rd.GVK
	err = api.Delete(provider, res, gvk)
	if err != nil {
		return err
	}
	fmt.Println("Resource successfully deleted!")
	return nil
}
