package translate

import (
	"fmt"

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

func FieldTypeForAttribute(tfAttr *configschema.Attribute) generator.AttributeType {
	switch tfAttr.Type.FriendlyName() {
	case "string": // TODO: figure out how to support "string"
		return generator.AttributeTypeString
	case "number": // TODO: figure out how to support "number"
		return generator.AttributeTypeInt
	case "bool": // TODO: figure out how to support "bool"
		return generator.AttributeTypeBool
	case "map of string": // TODO: figure out how to support "map of string"
		return generator.AttributeTypeUnsupported
	case "set of string": // TODO: figure out how to support "set of string"
		return generator.AttributeTypeUnsupported
	case "set of object": // TODO: figure out how to support "set of object"
		return generator.AttributeTypeUnsupported
	case "list of object": // TODO: figure out how to support "list of object"
		return generator.AttributeTypeUnsupported
	case "list of string": // TODO: figure out how to support "list of string"
		return generator.AttributeTypeUnsupported
	case "map of bool": // TODO: figure out how to support "map of bool"
		return generator.AttributeTypeUnsupported
	case "set of map of string": // TODO: oh... oh no
		return generator.AttributeTypeUnsupported
	}
	return generator.AttributeTypeUnsupported
}

// AttributeToField converts a terraform *configschema.Attribute
// to a crossplane generator.Field
func AttributeToField(name string, tfAttr *configschema.Attribute) generator.Field {
	f := generator.Field{
		Name:           strcase.ToCamel(name),
		Type:           generator.FieldTypeAttribute,
		AttributeField: generator.AttributeField{},
		Tag: &generator.StructTag{
			&generator.StructTagJson{
				Name: name,
			},
		},
	}
	switch tfAttr.Type.FriendlyName() {
	case "string":
		f.AttributeField.Type = generator.AttributeTypeString
	case "number":
		f.AttributeField.Type = generator.AttributeTypeInt
	case "bool":
		f.AttributeField.Type = generator.AttributeTypeBool
	case "map of string": // TODO: figure out how to support "map of string"
		f.AttributeField.Type = generator.AttributeTypeUnsupported
	case "set of string": // TODO: figure out how to support "set of string"
		f.AttributeField.Type = generator.AttributeTypeString
		f.IsSlice = true
		fmt.Printf("saw a []string! name=%s/%s\n", name, f.Name)
	case "set of object": // TODO: figure out how to support "set of object"
		f.AttributeField.Type = generator.AttributeTypeUnsupported
	case "list of object": // TODO: figure out how to support "list of object"
		f.AttributeField.Type = generator.AttributeTypeUnsupported
	case "list of string": // TODO: figure out how to support "list of string"
		f.AttributeField.Type = generator.AttributeTypeString
		f.IsSlice = true
		fmt.Printf("saw a []string (set)! name=%s/%s\n", name, f.Name)
	case "map of bool": // TODO: figure out how to support "map of bool"
		f.AttributeField.Type = generator.AttributeTypeUnsupported
	case "set of map of string": // TODO: oh... oh no
		f.AttributeField.Type = generator.AttributeTypeUnsupported
	default:
		f.AttributeField.Type = generator.AttributeTypeUnsupported
	}
	return f
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
