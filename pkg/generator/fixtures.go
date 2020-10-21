package generator

import (
	"io"
	"os"
	"path"

	"github.com/crossplane/terraform-provider-gen/pkg/template"
)

const fakeResourceName string = "Test"
const fakePackagePath string = "github.com/crossplane-contrib/fake"

func defaultTestResource() *ManagedResource {
	return NewManagedResource(fakeResourceName, fakePackagePath).WithNamer(NewDefaultNamer(fakeResourceName))
}

func getTestRenderTypesFileResult(tg template.TemplateGetter) (string, string, error) {
	fixturePath := "testdata/test-render-types-file.go"
	mr := defaultTestResource()
	result, err := RenderTypesFile(mr, tg)
	return result, fixturePath, err
}

func UpdateAllFixtures(basepath string) error {
	tg := template.NewTemplateGetter(basepath)
	fixtureGenerators := []func(template.TemplateGetter) (string, string, error){
		getTestRenderTypesFileResult,
	}
	for _, f := range fixtureGenerators {
		contents, fxpath, err := f(tg)
		p := path.Join(basepath, "pkg/generator", fxpath)
		f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		_, err = io.WriteString(f, contents)
		f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
