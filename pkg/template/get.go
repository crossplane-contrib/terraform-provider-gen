package template

//go:generate sh -c "CGO_ENABLED=0 go run ../../hack/.packr/main.go $PWD"

import (
	"bytes"
	"io"
	"os"
	"path"
	tmpl "text/template"
)

// TemplateGetter encapsulates the task of getting template data from
// the repository.
type TemplateGetter interface {
	// Get retrieves a Template parsed from the file indicated by the path argument.
	// The path argument is relative to the hack directory at the root of this repo.
	Get(path string) (*tmpl.Template, error)
}

type tplGetter struct {
	//files packr.Box
	basepath string
}

func fileAsString(fname string) (string, error) {
	fh, err := os.Open(fname)
	defer fh.Close()
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, fh)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (tg *tplGetter) Get(p string) (*tmpl.Template, error) {
	str, err := fileAsString(path.Join(tg.basepath, p))
	if err != nil {
		return nil, err
	}
	return tmpl.New(p).Parse(str)
}

// NewTemplateGetter returns a template getter that can find templates
// from paths within the root "hack" directory of the repository
// it assumes its own relative path within the project structure
// see tests get_tests.go to help understand how it is used.
func NewTemplateGetter(basepath string) TemplateGetter {
	return &tplGetter{
		basepath: basepath,
	}
}
