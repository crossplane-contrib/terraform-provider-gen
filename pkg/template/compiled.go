package template

import (
	"fmt"
	compiled "github.com/crossplane-contrib/terraform-provider-gen/internal/template/compiled"
	tmpl "text/template"
)

type compiledTemplateGetter struct {
}

func (ctg *compiledTemplateGetter) Get(path string) (*tmpl.Template, error) {
	cb, ok := compiled.TemplateDispatchMap[path]
	if !ok {
		return nil, fmt.Errorf("Compiled template not found for path %s", path)
	}
	return tmpl.New(path).Parse(cb())
}

func NewCompiledTemplateGetter() TemplateGetter {
	return &compiledTemplateGetter{}
}
