package generated

import (
	"{{ .RootPackage }}/generated/provider/{{ .ProviderConfigVersion }}"
	"github.com/crossplane-contrib/terraform-runtime/pkg/plugin"
)

const ProviderReferenceName string = "{{ .Name }}"

func Index(idxr *plugin.Indexer) error {
	for _, impl := range ResourceImplementations {
		err := idxr.Overlay(impl)
		if err != nil {
			return err
		}
	}
	return nil
}

func ProviderInit() *plugin.ProviderInit {
	return v1alpha1.GetProviderInit()
}
