package template

import (
	"bytes"
	"testing"
)

func TestTemplateGetter(t *testing.T) {
	g := NewTemplateGetter("../..")
	f, err := g.Get("hack/template/pkg/template/test-template-getter.txt")
	if err != nil {
		t.Error(err)
	}
	buf := new(bytes.Buffer)
	err = f.Execute(buf, struct{}{})
	if err != nil {
		t.Error(err)
	}
	actual := buf.String()
	expected := "see tests for TemplateGetter in pkg/template to understand why this file is here.\n"
	if actual != expected {
		t.Errorf("Expected to find fixture content:\n%s\nInstead found content\n%s\n", expected, actual)
	}
}
