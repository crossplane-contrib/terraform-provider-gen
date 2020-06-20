package api

import (
	"github.com/crossplane/hiveworld/pkg/client"
	"github.com/crossplane/hiveworld/pkg/resource"
)

// Read returns an up-to-date version of the resource
// TODO: If `id` is unset for a new resource, how do we figure out
// what value needs to be used as the id?
func Read(p *client.Provider, r resource.Representer) (resource.Representer, error) {
	asYAML, err := r.AsYAML()
	if err != nil {
		return nil, err
	}
	return resource.NewYAMLByteRepresenter(asYAML)
}
