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

func nestedFieldFixture(nestedTypeName, deeplyNestedTypeName string) Field {
	f := Field{
		Name: deeplyNestedTypeName + "Name",
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: fakePackagePath,
			TypeName:    deeplyNestedTypeName,
		},
		Fields: []Field{
			{
				Name:           "aString",
				Type:           FieldTypeAttribute,
				AttributeField: AttributeField{Type: AttributeTypeString},
				Tag: &StructTag{
					Json: &StructTagJson{
						Name: "a_string",
					},
				},
			},
		},
		Tag: &StructTag{
			Json: &StructTagJson{
				Name: "deeper_sub_field",
			},
		},
	}
	nf := Field{
		Name: nestedTypeName + "Name",
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: fakePackagePath,
			TypeName:    nestedTypeName,
		},
		Fields: []Field{
			f,
		},
		Tag: &StructTag{
			Json: &StructTagJson{
				Name: "sub_field",
			},
		},
	}
	test := Field{
		Name: fakeResourceName + "Name",
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: fakePackagePath,
			TypeName:    fakeResourceName,
		},
		Fields: []Field{
			nf,
		},
	}
	return test
}

type fixtureGenerator func(template.TemplateGetter) (string, error)

var (
	TestManagedResourceTypeDefRendererPath = "testdata/test-render-types-file.go"
	TestRenderNestedStatusPath             = "testdata/test-render-nested-status.go"
	TestRenderNestedSpecPath               = "testdata/test-render-nested-spec.go"
)

var FixtureGenerators map[string]fixtureGenerator = map[string]fixtureGenerator{
	TestManagedResourceTypeDefRendererPath: func(tg template.TemplateGetter) (string, error) {
		mr := defaultTestResource()
		renderer := NewManagedResourceTypeDefRenderer(mr, tg)
		result, err := renderer.Render()
		return result, err
	},
	TestRenderNestedStatusPath: func(tg template.TemplateGetter) (string, error) {
		mr := defaultTestResource()
		mr.Observation.Fields = []Field{nestedFieldFixture("nestedField", "deeplyNestedField")}
		renderer := NewManagedResourceTypeDefRenderer(mr, tg)
		return renderer.Render()
	},
	TestRenderNestedSpecPath: func(tg template.TemplateGetter) (string, error) {
		mr := defaultTestResource()
		// TODO: wonky thing that we have to do to satisfy matching package names to exclude
		// the qualifier. Might want to add fakePackagePath as an arg to the fixture instead
		// of assuming it everywhere
		mr.Parameters.StructField.PackagePath = fakePackagePath
		mr.Parameters.Fields = []Field{nestedFieldFixture("nestedField", "deeplyNestedField")}
		renderer := NewManagedResourceTypeDefRenderer(mr, tg)
		return renderer.Render()
	},
}

func UpdateAllFixtures(basepath string) error {
	tg := template.NewTemplateGetter(basepath)
	for fxpath, f := range FixtureGenerators {
		contents, err := f(tg)
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
