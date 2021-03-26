package generated

func IndexResources() string {
	return "package generated\n\nimport (\n{{- range .PackageImports }}\n    {{ .Name }} \"{{ .Path }}\"\n{{- end }}\n\n    \"github.com/crossplane-contrib/terraform-runtime/pkg/plugin\"\n)\n\nvar generatedImplementations = []*plugin.Implementation{\n{{- range .PackageImports }}\n    {{ .Name }}.Implementation(),\n{{- end}}\n}\n\n// this is deferred until init time to simplify the codegen workflow.\n// index.go can be a simple templated, satisfying the needs of main.go so that\n// the provider can be compiled (albeit in a non-functional state) enabling angryjet\n// and controller-gen to run against the generated types.go before the a subsequent pass\n// of terraform-provider-gen adds the compare/encode/decode methods.\nfunc init() {\n    for _, impl := range generatedImplementations {\n        resourceImplementations = append(resourceImplementations, impl)\n    }\n}\n"
}
