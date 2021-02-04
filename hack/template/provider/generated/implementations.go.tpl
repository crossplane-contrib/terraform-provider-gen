package generated

import (
{{- range .PackageImports }}
    {{ .Name }} "{{ .Path }}"
{{- end }}

    "github.com/crossplane-contrib/terraform-runtime/pkg/plugin"
)

var ResourceImplementations = []*plugin.Implementation{
{{- range .PackageImports }}
    {{ .Name }}.Implementation(),
{{- end}}
}
