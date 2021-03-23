package generated

import (
	"github.com/crossplane-contrib/terraform-runtime/pkg/plugin"
)

var resourceImplementations = make([]*plugin.Implementation, 0)
var providerInit *plugin.ProviderInit

// Index provides a plugin.Index for the generated provider
// note that the value of resourceImplementations is populated
// at runtime in implementations.go. This is to enable the separation of
// the build into multiple stages.
func Index() (*plugin.Index, error) {
	idxr := plugin.NewIndexer()
	for _, impl := range resourceImplementations {
		err := idxr.Overlay(impl)
		if err != nil {
			return nil, err
		}
	}
	return idxr.BuildIndex()
}

// Index provides a plugin.ProviderInit for the generated provider.
// Note that the value of providerInit is populated
// at runtime in provider.go. This is to enable the separation of
// the build into multiple stages.
func ProviderInit() *plugin.ProviderInit {
	return providerInit
}
