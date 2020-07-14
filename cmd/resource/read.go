package resource

import (
	"fmt"

	"github.com/crossplane/hiveworld/pkg/api"
	"github.com/crossplane/hiveworld/pkg/client"
	"github.com/crossplane/hiveworld/pkg/registry"
	"github.com/crossplane/hiveworld/pkg/resource"
)

// ReadResource will read a resource description
func ReadResource(resourceReadPath string, provider *client.Provider) error {
	rd, err := resource.ResourceDataFromFile(resourceReadPath)
	if err != nil {
		return err
	}
	res, err := rd.ManagedResource()
	if err != nil {
		return err
	}
	//meta.AddAnnotations(res, map[string]string{"crossplane.io/external-name": "testing"})
	gvk := rd.GVK
	newRes, err := api.Read(provider, res, gvk)
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
	/*
		representer, err := resource.RepresenterFromYAMLFile(resourceReadPath, provider)
		if err != nil {
			return err
		}

		updated, err := api.Read(provider, representer)
		if err != nil {
			return err
		}

		asYaml, err := updated.AsYAML()
		if err != nil {
			return err
		}
		fmt.Println(string(asYaml))

		return nil
	*/
}
