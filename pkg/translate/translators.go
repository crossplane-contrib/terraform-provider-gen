package translate

import (
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/generator"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/providers"
	"github.com/iancoleman/strcase"
)

type SpecOrStatusField int

const (
	ForProviderField SpecOrStatusField = iota
	AtProviderField
)

// AttributeToField converts a terraform *configschema.Attribute
// to a crossplane generator.Field
func AttributeToField(name string, tfAttr *configschema.Attribute) generator.Field {
	return generator.Field{
		Name:           strcase.ToCamel(name),
		Type:           generator.FieldTypeAttribute,
		AttributeField: generator.AttributeField{Type: generator.AttributeTypeString},
		Tag: &generator.StructTag{
			&generator.StructTagJson{
				Name: name,
			},
		},
	}
}

func SpecOrStatus(attr *configschema.Attribute) SpecOrStatusField {
	if attr.Computed {
		return AtProviderField
	}
	return ForProviderField
}

// SpecStatusAttributeFields iterates through the terraform configschema.Attribute map
// found under Block.Attributes, translating each attribute to a generator.Field and
// grouping them as spec or status based on their optional/required/computed properties.
func SpecOrStatusAttributeFields(attributes map[string]*configschema.Attribute) ([]generator.Field, []generator.Field) {
	forProvider := make([]generator.Field, 0)
	atProvider := make([]generator.Field, 0)
	for name, attr := range attributes {
		f := AttributeToField(name, attr)
		switch SpecOrStatus(attr) {
		case ForProviderField:
			forProvider = append(forProvider, f)
		case AtProviderField:
			atProvider = append(atProvider, f)
		}
	}
	return forProvider, atProvider
}

func SchemaToManagedResource(name, packagePath string, s providers.Schema) *generator.ManagedResource {
	namer := generator.NewDefaultNamer(name)
	mr := generator.NewManagedResource(namer.TypeName(), packagePath).WithNamer(namer)
	spec, status := SpecOrStatusAttributeFields(s.Block.Attributes)
	mr.Parameters = generator.Field{
		Type: generator.FieldTypeStruct,
		StructField: generator.StructField{
			PackagePath: packagePath,
			TypeName:    namer.ForProviderTypeName(),
		},
		Fields: spec,
	}
	mr.Observation = generator.Field{
		Type: generator.FieldTypeStruct,
		StructField: generator.StructField{
			PackagePath: packagePath,
			TypeName:    namer.AtProviderTypeName(),
		},
		Fields: status,
	}
	return mr
}

func BlockToField(name, typeName string, tfBlock *configschema.Block, enclosingField *generator.Field) *generator.Field {
	f := &generator.Field{
		Name: name,
		Type: generator.FieldTypeStruct,
		StructField: generator.StructField{
			PackagePath: enclosingField.StructField.PackagePath,
			TypeName:    typeName,
		},
	}
	return f
}
