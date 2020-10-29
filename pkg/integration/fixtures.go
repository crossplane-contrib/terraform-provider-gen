package integration

import (
	"io"
	"os"
	"path"

	"github.com/crossplane-contrib/terraform-provider-gen/pkg/generator"
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/template"
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/translate"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/providers"
	"github.com/zclconf/go-cty/cty"
)

const FakeResourceName string = "Test"
const FakePackagePath string = "github.com/crossplane-contrib/fake"

func DefaultTestResource() *generator.ManagedResource {
	return generator.NewManagedResource(FakeResourceName, FakePackagePath).WithNamer(generator.NewDefaultNamer(FakeResourceName))
}

func NestedFieldFixture(outerTypeName, nestedTypeName, deeplyNestedTypeName string) generator.Field {
	f := generator.Field{
		// "Name" is appended to help visually differentiate field and type names
		Name: deeplyNestedTypeName + "Name",
		Type: generator.FieldTypeStruct,
		StructField: generator.StructField{
			PackagePath: FakePackagePath,
			TypeName:    deeplyNestedTypeName,
		},
		Fields: []generator.Field{
			{
				Name:           "aString",
				Type:           generator.FieldTypeAttribute,
				AttributeField: generator.AttributeField{Type: generator.AttributeTypeString},
				Tag: &generator.StructTag{
					Json: &generator.StructTagJson{
						Name: "a_string",
					},
				},
			},
		},
		Tag: &generator.StructTag{
			Json: &generator.StructTagJson{
				Name: "deeper_sub_field",
			},
		},
	}
	nf := generator.Field{
		// "Name" is appended to help visually differentiate field and type names
		Name: nestedTypeName + "Name",
		Type: generator.FieldTypeStruct,
		StructField: generator.StructField{
			PackagePath: FakePackagePath,
			TypeName:    nestedTypeName,
		},
		Fields: []generator.Field{
			f,
		},
		Tag: &generator.StructTag{
			Json: &generator.StructTagJson{
				Name: "sub_field",
			},
		},
	}
	test := generator.Field{
		// "Name" is appended to help visually differentiate field and type names
		Name: outerTypeName + "Name",
		Type: generator.FieldTypeStruct,
		StructField: generator.StructField{
			PackagePath: FakePackagePath,
			TypeName:    outerTypeName,
		},
		Fields: []generator.Field{
			nf,
		},
	}
	return test
}

func testFixtureFlatBlock() providers.Schema {
	s := providers.Schema{
		Block: &configschema.Block{
			Attributes: make(map[string]*configschema.Attribute),
			BlockTypes: make(map[string]*configschema.NestedBlock),
		},
	}
	// I think "id" should probably not be part of the schema, it is like our external-name
	// TODO: check how this was implemented in the prototype
	//s.Block.Attributes["id"] =
	s.Block.Attributes["different_resource_ref_id"] = &configschema.Attribute{
		Required: false,
		Optional: true,
		Computed: false,
		Type:     cty.String,
	}
	s.Block.Attributes["perform_optional_action"] = &configschema.Attribute{
		Required: false,
		Optional: true,
		Computed: false,
		Type:     cty.Bool,
	}
	s.Block.Attributes["labels"] = &configschema.Attribute{
		Required: false,
		Optional: true,
		Computed: false,
		Type:     cty.Map(cty.String),
	}
	s.Block.Attributes["number_list"] = &configschema.Attribute{
		Required: false,
		Optional: true,
		Computed: false,
		Type:     cty.List(cty.Number),
	}
	s.Block.Attributes["computed_owner_id"] = &configschema.Attribute{
		Required: false,
		Optional: false,
		Computed: true,
		Type:     cty.String,
	}
	s.Block.Attributes["required_name"] = &configschema.Attribute{
		Required: true,
		Optional: false,
		Computed: false,
		Type:     cty.String,
	}
	return s
}

type fixtureGenerator func(template.TemplateGetter) (string, error)

var (
	TestManagedResourceTypeDefRendererPath = "testdata/test-render-types-file.go"
	TestRenderNestedStatusPath             = "testdata/test-render-nested-status.go"
	TestRenderNestedSpecPath               = "testdata/test-render-nested-spec.go"
	TestSchemaToManagedResourceRender      = "testdata/test-schema-to-managed-resource-render.go"
)

var FixtureGenerators map[string]fixtureGenerator = map[string]fixtureGenerator{
	TestManagedResourceTypeDefRendererPath: func(tg template.TemplateGetter) (string, error) {
		mr := DefaultTestResource()
		renderer := generator.NewManagedResourceTypeDefRenderer(mr, tg)
		result, err := renderer.Render()
		return result, err
	},
	TestRenderNestedStatusPath: func(tg template.TemplateGetter) (string, error) {
		mr := DefaultTestResource()
		// TODO: wonky thing that we have to do to satisfy matching package names to exclude
		// the qualifier. Might want to add generator.FakePackagePath as an arg to the fixture instead
		// of assuming it everywhere
		mr.Observation.StructField.PackagePath = FakePackagePath
		mr.Observation.Fields = []generator.Field{NestedFieldFixture("SubObservation", "nestedField", "deeplyNestedField")}
		renderer := generator.NewManagedResourceTypeDefRenderer(mr, tg)
		return renderer.Render()
	},
	TestRenderNestedSpecPath: func(tg template.TemplateGetter) (string, error) {
		mr := DefaultTestResource()
		// TODO: wonky thing that we have to do to satisfy matching package names to exclude
		// the qualifier. Might want to add generator.FakePackagePath as an arg to the fixture instead
		// of assuming it everywhere
		mr.Parameters.StructField.PackagePath = FakePackagePath
		mr.Parameters.Fields = []generator.Field{NestedFieldFixture("SubParameters", "nestedField", "deeplyNestedField")}
		renderer := generator.NewManagedResourceTypeDefRenderer(mr, tg)
		return renderer.Render()
	},
	TestSchemaToManagedResourceRender: func(tg template.TemplateGetter) (string, error) {
		terraform_resource_name := "aws_test_resource"
		// TODO: write some package naming stuff -- maybe start with a flat package name scheme
		packagePath := "github.com/crossplane/provider-terraform-aws/generated/test/v1alpha1"
		renamer := translate.NewTerraformResourceRenamer("aws")
		s := testFixtureFlatBlock()
		mr := translate.SchemaToManagedResource(renamer(terraform_resource_name), packagePath, s)
		renderer := generator.NewManagedResourceTypeDefRenderer(mr, tg)
		return renderer.Render()
	},
}

func UpdateAllFixtures(basepath string) error {
	tg := template.NewTemplateGetter(basepath)
	for fxpath, f := range FixtureGenerators {
		contents, err := f(tg)
		p := path.Join(basepath, "pkg/integration", fxpath)
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
