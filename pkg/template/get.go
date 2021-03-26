package template

import (
	tmpl "text/template"
)

// TemplateGetter encapsulates the task of getting template data from
// the repository.
type TemplateGetter interface {
	// Get retrieves a Template parsed from the file indicated by the path argument.
	// The path argument is relative to the hack directory at the root of this repo.
	Get(path string) (*tmpl.Template, error)
}
