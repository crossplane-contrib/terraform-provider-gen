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
	case "set of object": // TODO: figure out how to support "set of object"
		f.AttributeField.Type = generator.AttributeTypeUnsupported
	case "list of object": // TODO: figure out how to support "list of object"
		f.AttributeField.Type = generator.AttributeTypeUnsupported
	case "list of string": // TODO: figure out how to support "list of string"
		f.AttributeField.Type = generator.AttributeTypeString
		f.IsSlice = true
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

func NestedBlockFields(blocks map[string]*configschema.NestedBlock, packagePath string) []generator.Field {
	fields := make([]generator.Field, 0)
	for name, block := range blocks {
		f := generator.Field{
			Name:   strcase.ToCamel(name),
			Fields: make([]generator.Field, 0),
			Type:   generator.FieldTypeStruct,
			StructField: generator.StructField{
				PackagePath: packagePath,
				// TODO: the output would look nicer if we pluralized names when IsBlockList is true
				TypeName: strcase.ToCamel(name),
			},
			Tag: &generator.StructTag{
				&generator.StructTagJson{
					Name: name,
				},
			},
			Required: IsBlockRequired(block),
			IsSlice:  IsBlockSlice(block),
		}
		for n, attr := range block.Attributes {
			f.Fields = append(f.Fields, AttributeToField(n, attr))
		}
		f.Fields = append(f.Fields, NestedBlockFields(block.BlockTypes, packagePath)...)
		fields = append(fields, f)
	}
	return fields
}

func SchemaToManagedResource(name, packagePath string, s providers.Schema) *generator.ManagedResource {
	namer := generator.NewDefaultNamer(strcase.ToCamel(name))
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
	nb := NestedBlockFields(s.Block.BlockTypes, packagePath)
	if len(nb) > 0 {
		mr.Parameters.Fields = append(mr.Parameters.Fields, nb...)
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

func IsBlockRequired(nb *configschema.NestedBlock) bool {
	if nb.MinItems > 0 {
		return true
	}
	return false
}

func IsBlockSlice(nb *configschema.NestedBlock) bool {
	if nb.MaxItems != 1 {
		return true
	}
	return false
}
