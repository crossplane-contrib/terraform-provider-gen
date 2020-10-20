package template

import (
	tmpl "text/template"

	"github.com/gobuffalo/packr"
)

var HackRelativePath = "../../hack"

// TemplateGetter provides a Get method which retrieves a Template
// parsed from the file indicated by the path argument. The path
// argument is relative to the hack directory at the root of this repo.
type TemplateGetter interface {
	Get(path string) (*tmpl.Template, error)
}

type tplGetter struct {
	files packr.Box
}

func (tg *tplGetter) Get(path string) (*tmpl.Template, error) {
	str, err := tg.files.FindString(path)
	if err != nil {
		return nil, err
	}
	return tmpl.New(path).Parse(str)
}

// NewTemplateGetter returns a template getter that can find templates
// from paths within the root "hack" directory of the repository
// it assumes its own relative path within the project structure
// see tests get_tests.go to help understand how it is used.
func NewTemplateGetter() TemplateGetter {
	return &tplGetter{
		files: packr.NewBox(HackRelativePath),
	}
}
