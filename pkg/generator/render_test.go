package generator

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/crossplane/terraform-provider-gen/pkg/template"
)

func getFixture(path string) (string, error) {
	fh, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("Error while trying to read fixture path %s: %s", path, err)
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, fh)
	if err != nil {
		return "", fmt.Errorf("Error while trying to read fixture path %s: %s", path, err)
	}
	return buf.String(), nil
}

func TestRenderTypesFile(t *testing.T) {
	tg := template.NewTemplateGetter("../../")
	actual, fxpath, err := getTestRenderTypesFileResult(tg)
	expected, err := getFixture(fxpath)
	if err != nil {
		t.Error(err)
	}

	if err != nil {
		t.Error(err)
	}
	if actual != expected {
		t.Errorf("Unexpected output from RenderTypesFile.\nExpected:\n ---- \n%s\n ---- \nActual:\n ---- \n%s\n ---- \n", expected, actual)
	}
}
