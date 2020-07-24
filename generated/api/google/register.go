package google

import (
	"github.com/crossplane/provider-terraform-plugin/pkg/registry"
)

const ProviderReferenceName string = "google"

func Register(r *registry.Registry) {
	for _, entry := range RegistryEntries {
		r.Register(entry)
	}

	r.RegisterProvider(ProviderEntry)
}
