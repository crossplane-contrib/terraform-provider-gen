package resource

import (
	"fmt"

	"github.com/crossplane/hiveworld/pkg/api"
	"github.com/crossplane/hiveworld/pkg/client"
	"github.com/crossplane/hiveworld/pkg/resource"
)

// ReadResource will read a resource description
func ReadResource(resourceReadPath string, provider *client.Provider) error {
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
}
