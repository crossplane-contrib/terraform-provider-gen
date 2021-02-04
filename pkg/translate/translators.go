package translate

import (
	"fmt"

	"github.com/crossplane-contrib/terraform-provider-gen/pkg/generator"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/providers"
	"github.com/iancoleman/strcase"
	"github.com/zclconf/go-cty/cty"
)

type SpecOrStatusField int

const (
	ForProviderField SpecOrStatusField = iota
	AtProviderField
)

type FieldBuilder struct {
	f *generator.Field
}

func isReferenceType(t cty.Type) bool {
	if t.IsMapType() {
		return true
	}
	if t.IsCollectionType() {
		return true
	}
	return false
}

func NewFieldBuilder(name string, ctyType cty.Type) *FieldBuilder {
	encFnGen := NewAttributeEncodeFnGenerator(name, ctyType)
	decFnGen := NewAttributeDecodeFnGenerator(name, ctyType)
	mergeFnGen := NewAttributeMergeFnGenerator(name, ctyType)
	st := &generator.StructTag{
		Json: &generator.StructTagJson{
			Name: name,
		},
	}
	if isReferenceType(ctyType) {
		st.Json.Omitempty = true
	}
	return &FieldBuilder{
		f: &generator.Field{
			Name: strcase.ToCamel(name),
			Tag: st,
			EncodeFnGenerator: encFnGen,
			DecodeFnGenerator: decFnGen,
			MergeFnGenerator:  mergeFnGen,
		},
	}
}

func (fb *FieldBuilder) AttributeField(af generator.AttributeField) *FieldBuilder {
	fb.f.Type = generator.FieldTypeAttribute
	fb.f.AttributeField = af
	return fb
}

func (fb *FieldBuilder) StructField(typeName string, fields []generator.Field) *FieldBuilder {
	fb.f.Type = generator.FieldTypeStruct
	// do we need to plumb through package path? that whole concept is getting tiresome
	// and isn't used for anything so far
	fb.f.StructField = generator.StructField{
		TypeName: typeName,
	}
	fb.f.Fields = fields
	return fb
}

func (fb *FieldBuilder) IsSlice(is bool) *FieldBuilder {
	fb.f.IsSlice = is
	return fb
}

func (fb *FieldBuilder) Unsupported() generator.Field {
	return fb.AttributeField(
		generator.AttributeField{Type: generator.AttributeTypeUnsupported}).Build()
}

func (fb *FieldBuilder) ObjectField(typeName string, attrType cty.Type, schemaPath string) *FieldBuilder {
	fields := make([]generator.Field, 0)
	for k, t := range attrType.ElementType().AttributeTypes() {
		fields = append(fields, TypeToField(k, t, schemaPath))
	}
	return fb.StructField(typeName, fields)
}

func (fb *FieldBuilder) Build() generator.Field {
	return *fb.f
}

// TypeToField converts a terraform *configschema.Attribute
// to a crossplane generator.Field
func TypeToField(name string, attrType cty.Type, parentPath string) generator.Field {
	sp := appendToSchemaPath(parentPath, name)
	fb := NewFieldBuilder(name, attrType)
	switch attrType.FriendlyName() {
	case "bool":
		return fb.AttributeField(
			generator.AttributeField{Type: generator.AttributeTypeBool}).Build()
	case "number":
		return fb.AttributeField(
			generator.AttributeField{Type: generator.AttributeTypeInt64}).Build()
	case "string":
		return fb.AttributeField(
			generator.AttributeField{Type: generator.AttributeTypeString}).Build()
	case "map of bool":
		return fb.AttributeField(
			generator.AttributeField{
				Type:         generator.AttributeTypeMapStringKey,
				MapValueType: generator.AttributeTypeBool,
			}).Build()
	case "map of number":
		return fb.AttributeField(
			generator.AttributeField{
				Type:         generator.AttributeTypeMapStringKey,
				MapValueType: generator.AttributeTypeInt64,
			}).Build()
	case "map of string":
		return fb.AttributeField(
			generator.AttributeField{
				Type:         generator.AttributeTypeMapStringKey,
				MapValueType: generator.AttributeTypeString,
			}).Build()
	case "list of number":
		return fb.IsSlice(true).AttributeField(
			generator.AttributeField{Type: generator.AttributeTypeInt64}).Build()
	case "list of string":
		return fb.IsSlice(true).AttributeField(
			generator.AttributeField{Type: generator.AttributeTypeString}).Build()
	case "set of number":
		return fb.IsSlice(true).AttributeField(
			generator.AttributeField{Type: generator.AttributeTypeInt64}).Build()
	case "set of string":
		return fb.IsSlice(true).AttributeField(
			generator.AttributeField{Type: generator.AttributeTypeString}).Build()
	case "set of map of string":
		return fb.IsSlice(true).AttributeField(
			generator.AttributeField{
				Type:         generator.AttributeTypeMapStringKey,
				MapValueType: generator.AttributeTypeString,
			}).Build()

	// TODO: the set/list of objects types can probably be []map[string]string
	// but we need to spot check and confirm this.
	case "list of object": // TODO: probably can be []map[string]string
		//f.AttributeField.Type = generator.AttributeTypeUnsupported
		// TODO: see note on "set of object" re object schemas, this may also apply
		// to constructing lists of object
		if !attrType.IsListType() {
			return fb.Unsupported()
		}
		if !attrType.ElementType().IsObjectType() {
			return fb.Unsupported()
		}
		return fb.IsSlice(true).ObjectField(strcase.ToCamel(name), attrType, sp).Build()
	case "set of object":
		// TODO: sets of objects have a fixed schema that we need to track in order to
		// marshal them later. I think it may be an error to try to declare a set with
		// an object definition consisting just of set fields (a subset of the fields)
		// i've also seen (in provider configs) errors when optional set types are not
		// declared. look at the provider config construction in provider_terraform_aws
		// for an example.
		if !attrType.IsSetType() {
			return fb.Unsupported()
		}
		if !attrType.ElementType().IsObjectType() {
			return fb.Unsupported()
		}
		return fb.IsSlice(true).ObjectField(strcase.ToCamel(name), attrType, sp).Build()
	default:
		// TODO: need better error handling here to help generate error messages
		// which would describe why the field is unsupported
		// maybe this panic, either here or further up the stack
		return fb.Unsupported()
	}
}

func SpecOrStatus(attr *configschema.Attribute) SpecOrStatusField {
	// if attr.Computed is true, it can either be an attribute (status) or an argument (spec)
	// but arguments will always either be required or optional
	if attr.Required || attr.Optional {
		return ForProviderField
	}
	return AtProviderField
}

func appendToSchemaPath(sp, name string) string {
	return fmt.Sprintf("%s_%s", sp, name)
}

// SpecStatusAttributeFields iterates through the terraform configschema.Attribute map
// found under Block.Attributes, translating each attribute to a generator.Field and
// grouping them as spec or status based on their optional/required/computed properties.
func SpecOrStatusAttributeFields(attributes map[string]*configschema.Attribute, namer generator.ResourceNamer) ([]generator.Field, []generator.Field) {
	forProvider := make([]generator.Field, 0)
	atProvider := make([]generator.Field, 0)
	forProviderPath := fmt.Sprintf("%s_%s_%s", namer.TypeName(), namer.SpecTypeName(), namer.ForProviderTypeName())
	atProviderPath := fmt.Sprintf("%s_%s_%s", namer.TypeName(), namer.StatusTypeName(), namer.AtProviderTypeName())
	for name, attr := range attributes {
		// filter the top-level terraform id field out of the schema, these are
		// manually handled in the generated encode/decode methods
		if name == "id" {
			continue
		}
		switch SpecOrStatus(attr) {
		case ForProviderField:
			f := TypeToField(name, attr.Type, forProviderPath)
			forProvider = append(forProvider, f)
		case AtProviderField:
			f := TypeToField(name, attr.Type, atProviderPath)
			atProvider = append(atProvider, f)
		}
	}
	return forProvider, atProvider
}

var (
	ctyListCollectionType = cty.List(cty.EmptyObject)
	ctySetCollectionType  = cty.Set(cty.EmptyObject)
	ctyMapCollectionType  = cty.Map(cty.EmptyObject)
)

func NestedBlockFields(blocks map[string]*configschema.NestedBlock, packagePath, schemaPath string) []generator.Field {
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
				Json: &generator.StructTagJson{
					Name: name,
				},
			},
			Required: IsBlockRequired(block),
			IsSlice:  IsBlockSlice(block),
		}
		f.EncodeFnGenerator = NewBlockEncodeFnGenerator(name, block)
		f.DecodeFnGenerator = NewBlockDecodeFnGenerator(name, block)
		f.MergeFnGenerator = NewBlockMergeFnGenerator(name, block)

		sp := appendToSchemaPath(schemaPath, f.Name)
		for n, attr := range block.Attributes {
			f.Fields = append(f.Fields, TypeToField(n, attr.Type, sp))
		}
		f.Fields = append(f.Fields, NestedBlockFields(block.BlockTypes, packagePath, sp)...)
		fields = append(fields, f)
	}
	return fields
}

func SchemaToManagedResource(name, packagePath string, s providers.Schema) *generator.ManagedResource {
	namer := generator.NewDefaultNamer(strcase.ToCamel(name))
	mr := generator.NewManagedResource(namer.TypeName(), packagePath).WithNamer(namer)
	spec, status := SpecOrStatusAttributeFields(s.Block.Attributes, namer)
	mr.Parameters = generator.Field{
		Tag: &generator.StructTag{
			Json: &generator.StructTagJson{
				Name: "forProvider",
			},
		},
		Type: generator.FieldTypeStruct,
		StructField: generator.StructField{
			PackagePath: packagePath,
			TypeName:    namer.ForProviderTypeName(),
		},
		Fields: spec,
		Name:   namer.ForProviderTypeName(),
	}
	mr.Observation = generator.Field{
		Tag: &generator.StructTag{
			Json: &generator.StructTagJson{
				Name: "atProvider",
			},
		},
		Type: generator.FieldTypeStruct,
		StructField: generator.StructField{
			PackagePath: packagePath,
			TypeName:    namer.AtProviderTypeName(),
		},
		Fields: status,
		Name:   namer.AtProviderTypeName(),
	}
	nb := NestedBlockFields(s.Block.BlockTypes, packagePath, namer.TypeName())
	if len(nb) > 0 {
		// currently the assumption is that the nested types are spec fields
		// TODO: write an analyzer to ensure that deeply nested types are not common in status
		// we could do tree search into the structure of a NestedBlock
		mr.Parameters.Fields = append(mr.Parameters.Fields, nb...)
	}
	return mr
}

func IsBlockRequired(nb *configschema.NestedBlock) bool {
	if nb.MinItems > 0 {
		return true
	}
	return false
}

func IsBlockSlice(nb *configschema.NestedBlock) bool {
	// this is used to indicate a single optional block, like aws_db_proxy.timeouts
	// idk why it does not have a MaxItems = 1
	// take a look w/:
	// ./terraform-provider-gen analyze --plugin-path=$PLUGIN_PATH --providerName=$PROVIDER_NAME nesting | grep 'NestingList (0, 0'
	if nb.MaxItems == 0 && nb.MinItems == 0 {
		return false
	}
	if nb.MaxItems != 1 {
		return true
	}
	return false
}
