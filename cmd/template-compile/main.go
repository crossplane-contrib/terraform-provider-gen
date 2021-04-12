package main

import (
	"github.com/alecthomas/kong"
	"github.com/crossplane-contrib/terraform-provider-gen/internal/template"
	"github.com/spf13/afero"
)

var CompileTemplate struct {
	Root        string `help:"path to terraform-provider-gen repository root (default CWD)."`
	Output      string `help:"Compiled templates will be written to this directory."`
	PackageRoot string `help:"Go package that compiled templates will exist within."`
}

func main() {
	ctx := kong.Parse(&CompileTemplate)
	ctx.FatalIfErrorf(compileTemplates(CompileTemplate.Root, CompileTemplate.Output, CompileTemplate.PackageRoot))
}

func compileTemplates(root, output, packageRoot string) error {
	templateDir := "hack/template"
	tc := template.NewTemplateCompiler(afero.NewOsFs(), root, templateDir, output, packageRoot)
	return tc.CompileGeneratedTemplates()
}
