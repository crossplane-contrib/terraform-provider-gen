package google

import (
	iam "github.com/crossplane/hiveworld/generated/api/google/iam/v1alpha1"
	"github.com/crossplane/hiveworld/pkg/registry"
)

func Register() {
	entries := make([]*registry.Entry, 0)
	entries = append(entries, iam.RegisteryEntry())
	for _, entry := range entries {
		registry.RegisterCtyEncodeFunc(entry.GVK, entry.EncodeCtyCallback)
		registry.RegisterCtyDecodeFunc(entry.GVK, entry.DecodeCtyCallback)
		registry.RegisterResourceUnmarshalFunc(entry.GVK, entry.UnmarshalResourceCallback)
		registry.RegisterTerraformNameMapping(entry.TerraformResourceName, entry.GVK)
		registry.RegisterYAMLEncodeFunc(entry.GVK, entry.YamlEncodeCallback)
	}
}
