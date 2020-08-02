package google

import (
	"github.com/crossplane/terraform-provider-runtime/pkg/registry"
)

const ProviderReferenceName string = "google"

func Register(r *registry.Registry) {
	for _, entry := range RegistryEntries {
		r.Register(entry)
	}

	r.RegisterProvider(ProviderEntry)
}
