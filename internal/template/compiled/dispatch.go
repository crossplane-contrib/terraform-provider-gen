package dispatch

import (
	pkgGenerator "github.com/crossplane-contrib/terraform-provider-gen/internal/template/compiled/pkg/generator"
	pkgTemplate "github.com/crossplane-contrib/terraform-provider-gen/internal/template/compiled/pkg/template"
	providerCmdProvider "github.com/crossplane-contrib/terraform-provider-gen/internal/template/compiled/provider/cmd/provider"
	providerGenerated "github.com/crossplane-contrib/terraform-provider-gen/internal/template/compiled/provider/generated"
	providerGeneratedProviderV1Alpha1 "github.com/crossplane-contrib/terraform-provider-gen/internal/template/compiled/provider/generated/provider/v1alpha1"
)

var TemplateDispatchMap map[string]func() string = map[string]func() string{
	"pkg/generator/compare.go.tmpl":                     pkgGenerator.Compare,
	"pkg/generator/configure.go.tmpl":                   pkgGenerator.Configure,
	"pkg/generator/decode.go.tmpl":                      pkgGenerator.Decode,
	"pkg/generator/doc.go.tmpl":                         pkgGenerator.Doc,
	"pkg/generator/encode.go.tmpl":                      pkgGenerator.Encode,
	"pkg/generator/index.go.tmpl":                       pkgGenerator.Index,
	"pkg/generator/types.go.tmpl":                       pkgGenerator.Types,
	"pkg/template/test-template-getter.txt":             pkgTemplate.TestTemplateGetter,
	"provider/cmd/provider/main.go.tpl":                 providerCmdProvider.Main,
	"provider/generated/index.go.tpl":                   providerGenerated.Index,
	"provider/generated/index_provider.go.tpl":          providerGenerated.IndexProvider,
	"provider/generated/index_resources.go.tpl":         providerGenerated.IndexResources,
	"provider/generated/provider/v1alpha1/doc.go.tpl":   providerGeneratedProviderV1Alpha1.Doc,
	"provider/generated/provider/v1alpha1/index.go.tpl": providerGeneratedProviderV1Alpha1.Index,
	"provider/generated/provider/v1alpha1/types.go.tpl": providerGeneratedProviderV1Alpha1.Types,
}
