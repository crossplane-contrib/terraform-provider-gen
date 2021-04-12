package template

import (
	"bytes"
	"io"
	"os"
	"path"
	tmpl "text/template"
)

type fsTplGetter struct {
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

func (tg *fsTplGetter) Get(p string) (*tmpl.Template, error) {
	str, err := fileAsString(path.Join(tg.basepath, p))
	if err != nil {
		return nil, err
	}
	return tmpl.New(p).Parse(str)
}

// NewFSTemplateGetter returns a template getter that can find templates
// from paths within the root "hack" directory of the repository
// it assumes its own relative path within the project structure
// see tests get_tests.go to help understand how it is used.
func NewFSTemplateGetter(basepath string) TemplateGetter {
	return &fsTplGetter{
		basepath: basepath,
	}
}
