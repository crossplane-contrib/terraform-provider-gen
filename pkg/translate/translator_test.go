package translate

import (
	"testing"

	"github.com/crossplane-contrib/terraform-provider-gen/pkg/generator"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/providers"
	"github.com/iancoleman/strcase"
	"github.com/zclconf/go-cty/cty"
)

func testFixtureOptionalStringField() (string, string, *configschema.Attribute) {
	return "optional_field", "OptionalField", &configschema.Attribute{
		Required: false,
		Optional: true,
		Computed: false,
		Type:     cty.String,
	}
}

func TestTypeToField(t *testing.T) {
	name, expectedName, attr := testFixtureOptionalStringField()
	f := TypeToField(name, attr.Type, "")
	if f.Name != expectedName {
		t.Errorf("Wrong value from TypeToField for Field.Name. expected=%s, actual=%s", expectedName, f.Name)
	}
	if (f.Type != generator.FieldTypeAttribute && f.AttributeField != generator.AttributeField{}) {
		t.Errorf("Expected TypeToField to return an Attribute field, instead saw=%s", f.Type.String())
	}
	if f.AttributeField.Type != generator.AttributeTypeString {
		t.Errorf("Expected attribute field to be a string type, instead saw =%s", f.AttributeField.Type.String())
	}
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

func TestSpecStatusAttributeFields(t *testing.T) {
	resourceName := "test"
	s := testFixtureFlatBlock()
	namer := generator.NewDefaultNamer(strcase.ToCamel(resourceName))
	fp, ap := SpecOrStatusAttributeFields(s.Block.Attributes, namer)
	total := len(s.Block.Attributes)
	expectedAP := 1
	expectedFP := total - expectedAP
	if len(fp) != expectedFP {
		t.Errorf("Expected %d/%d fields to be in ForProvider, saw=%d", expectedFP, total, len(fp))
	}
	if len(ap) != expectedAP {
		t.Errorf("Expected %d/%d fields to be in AtProvider, saw=%d", expectedAP, total, len(ap))
	}
}

func TestSchemaToManagedResourceRender(t *testing.T) {
	resourceName := "TestResource"
	// TODO: write some package naming stuff -- maybe start with a flat package name scheme
	packagePath := "github.com/crossplane/provider-terraform-aws/generated/test/v1alpha1"
	s := testFixtureFlatBlock()
	mr := SchemaToManagedResource(resourceName, packagePath, s)
	if mr.Name != mr.Namer().TypeName() {
		t.Errorf("expected ManagedResource.Name=%s, actual=%s", mr.Namer().TypeName(), mr.Name)
	}
}
