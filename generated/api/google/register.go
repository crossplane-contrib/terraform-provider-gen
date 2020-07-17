package google

import (
	"github.com/crossplane/provider-terraform-plugin/pkg/registry"
)

func Register(r *registry.Registry) {
	for _, entry := range RegistryEntries {
		/*
			r.RegisterCtyEncodeFunc(entry.GVK, entry.EncodeCtyCallback)
			r.RegisterCtyDecodeFunc(entry.GVK, entry.DecodeCtyCallback)
			r.RegisterResourceUnmarshalFunc(entry.GVK, entry.UnmarshalResourceCallback)
			r.RegisterTerraformNameMapping(entry.TerraformResourceName, entry.GVK)
			r.RegisterYAMLEncodeFunc(entry.GVK, entry.YamlEncodeCallback)
			r.RegisterReconcilerConfigureFunc(entry.GVK, entry.ReconcilerConfigurer)
		*/
		r.Register(entry)
	}
}
