package generated

import (
    "{{ .RootPackage }}/generated/provider/{{ .ProviderConfigVersion }}"
)

func init() {
    providerInit = {{ .ProviderConfigVersion }}.GetProviderInit()
}
