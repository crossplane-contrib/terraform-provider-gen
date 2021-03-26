package generated

func IndexProvider() string {
	return "package generated\n\nimport (\n    \"{{ .RootPackage }}/generated/provider/{{ .ProviderConfigVersion }}\"\n)\n\nfunc init() {\n    providerInit = {{ .ProviderConfigVersion }}.GetProviderInit()\n}\n"
}
