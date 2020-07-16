package google

import (
	"github.com/crossplane/hiveworld/pkg/registry"
)

func Register() {
	for _, entry := range RegistryEntries {
		registry.RegisterCtyEncodeFunc(entry.GVK, entry.EncodeCtyCallback)
		registry.RegisterCtyDecodeFunc(entry.GVK, entry.DecodeCtyCallback)
		registry.RegisterResourceUnmarshalFunc(entry.GVK, entry.UnmarshalResourceCallback)
		registry.RegisterTerraformNameMapping(entry.TerraformResourceName, entry.GVK)
		registry.RegisterYAMLEncodeFunc(entry.GVK, entry.YamlEncodeCallback)
	}
}
