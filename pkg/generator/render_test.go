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

func TestManagedResourceTypeDefRenderer(t *testing.T) {
	if err := AssertConsistentFixture(TestManagedResourceTypeDefRendererPath); err != nil {
		t.Error(err)
	}
}

func TestRenderNestedStatus(t *testing.T) {
	if err := AssertConsistentFixture(TestRenderNestedStatusPath); err != nil {
		t.Error(err)
	}
}

func TestRenderNestedSpec(t *testing.T) {
	if err := AssertConsistentFixture(TestRenderNestedSpecPath); err != nil {
		t.Error(err)
	}
}

func AssertConsistentFixture(fixturePath string) error {
	tg := template.NewTemplateGetter("../../")
	fr := FixtureGenerators[fixturePath]
	actual, err := fr(tg)
	if err != nil {
		return err
	}

	expected, err := getFixture(fixturePath)
	if err != nil {
		return err
	}
	if actual != expected {
		return fmt.Errorf("Unexpected output from managedResourceTypeDefRenderer.Render().\nExpected:\n ---- \n%s\n ---- \nActual:\n ---- \n%s\n ---- \n", expected, actual)
	}
	return nil
}
