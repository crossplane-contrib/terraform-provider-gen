package generated

import (
{{- range .PackageImports }}
    {{ .Name }} "{{ .Path }}"
{{- end }}

    "github.com/crossplane-contrib/terraform-runtime/pkg/plugin"
)

var generatedImplementations = []*plugin.Implementation{
{{- range .PackageImports }}
    {{ .Name }}.Implementation(),
{{- end}}
}

// this is deferred until init time to simplify the codegen workflow.
// index.go can be a simple templated, satisfying the needs of main.go so that
// the provider can be compiled (albeit in a non-functional state) enabling angryjet
// and controller-gen to run against the generated types.go before the a subsequent pass
// of terraform-provider-gen adds the compare/encode/decode methods.
func init() {
    for _, impl := range generatedImplementations {
        resourceImplementations = append(resourceImplementations, impl)
    }
}
