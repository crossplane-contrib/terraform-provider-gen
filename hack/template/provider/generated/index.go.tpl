package generated

import (
	"{{ .RootPackage }}/generated/provider/{{ .ProviderConfigVersion }}"
	"github.com/crossplane-contrib/terraform-runtime/pkg/plugin"
)

const ProviderReferenceName string = "{{ .Name }}"

func Index(idxr *plugin.Indexer) {
	for _, impl := range ResourceImplementations {
		idxr.Overlay(impl)
	}
}

func ProviderInit() *plugin.ProviderInit {
	return v1alpha1.GetProviderInit()
}
