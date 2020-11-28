package translate

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/crossplane-contrib/terraform-provider-gen/pkg/generator"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/zclconf/go-cty/cty"
)

type backTracker struct {
	TFName         string    // the terraform field name will be different from the crd field name and json tag
	CtyType        cty.Type  // type of field, to be used to select a cty val conversion function (eg cty.StringVal)
	CollectionType *cty.Type // if this is a list/set/map then reflect the collection type here for another layer of translation
}

func (bt *backTracker) encodeFnRenderer(funcPrefix, receivedType string, f generator.Field) *encodeFnRenderer {
	return &encodeFnRenderer{
		FuncName:           fmt.Sprintf("%s_%s", funcPrefix, f.Name),
		ParentType:         receivedType,
		TerraformFieldName: bt.TFName,
		StructFieldName:    f.Name,
		Children:           f.Fields,
		CtyType:            bt.CtyType,
		CollectionType:     bt.CollectionType,
	}
}

func (bt *backTracker) RenderEncodeFn(funcPrefix, receivedType string, f generator.Field) string {
	efr := bt.encodeFnRenderer(funcPrefix, receivedType, f)
	switch true {
	case bt.CtyType.IsPrimitiveType():
		if bt.CollectionType != nil {
			return renderPrimitiveCollectionType(efr)
		}
		return renderPrimitiveType(efr)
	case bt.CtyType.IsMapType() || bt.CtyType.IsObjectType():
		if bt.CollectionType != nil {
			return renderContainerCollectionType(efr)
		}
		return renderContainerType(efr)
	//case bt.CtyType.IsCollectionType():
	// matches on map, set and list, but map should early return
	//return renderContainerCollectionType(efr)
	default:
		panic(fmt.Sprintf("Unknown cty type in RenderEncodeFn(), cannot render encoder for: %s", bt.CtyType.FriendlyName()))
	}
}

func renderPrimitiveType(efr *encodeFnRenderer) string {
	b := bytes.NewBuffer(make([]byte, 0))
	encoderTemplates["primitive"].Execute(b, efr)

	return b.String()
}

func renderPrimitiveCollectionType(efr *encodeFnRenderer) string {
	b := bytes.NewBuffer(make([]byte, 0))
	encoderTemplates["primitiveCollection"].Execute(b, efr)

	return b.String()
}

func renderContainerType(efr *encodeFnRenderer) string {
	b := bytes.NewBuffer(make([]byte, 0))
	encoderTemplates["container"].Execute(b, efr)

	rendered := []string{b.String()}
	rendered = append(rendered, renderChildEncoders(efr)...)

	return strings.Join(rendered, "\n\n")
}

func renderContainerCollectionType(efr *encodeFnRenderer) string {
	b := bytes.NewBuffer(make([]byte, 0))
	encoderTemplates["containerCollection"].Execute(b, efr)

	rendered := []string{b.String()}
	rendered = append(rendered, renderChildEncoders(efr)...)
	return strings.Join(rendered, "\n\n")
}

// this has been extracted into its own function so that it can be used by containers and container collections
func renderChildEncoders(efr *encodeFnRenderer) []string {
	rendered := make([]string, 0)
	for _, child := range efr.Children {
		receivedType := efr.StructFieldName
		if child.Type == generator.FieldTypeStruct {
			receivedType = child.Name
		}
		rendered = append(rendered, child.EncodeFnRenderer.RenderEncodeFn(efr.FuncName, receivedType, child))
	}
	return rendered
}

type encodeFnRenderer struct {
	FuncName           string // child func name is constructed by appending the child field name
	ParentType         string // name of the parent type to be used in constructing the function receiver
	TerraformFieldName string // the terraform field name will be different from the crd field name and json tag
	StructFieldName    string // name of the child field in the parent struct, to be used for assignment
	Children           []generator.Field
	CtyType            cty.Type
	CollectionType     *cty.Type
}

func (efr *encodeFnRenderer) ConversionFunc() string {
	switch efr.CtyType {
	case cty.String:
		return "cty.StringVal"
	case cty.Bool:
		return "cty.BoolVal"
	case cty.Number:
		return "cty.IntVal"
	}
	if efr.CtyType.IsObjectType() {
		return "cty.ObjectVal"
	}
	if efr.CtyType.IsMapType() {
		return "cty.MapVal"
	}
	panic(fmt.Sprintf("Unknown cty type in ConversionFunc(), cannot render convert function for: %s", efr.CtyType.FriendlyName()))
}

func (efr *encodeFnRenderer) CollectionConversionFunc() string {
	if efr.CollectionType.IsSetType() {
		return "cty.SetVal"
	}
	if efr.CollectionType.IsListType() {
		return "cty.ListVal"
	}
	panic(fmt.Sprintf("Unknown CollectionType in CollectionConversionFunc(), cannot render convert function for: %s", efr.CollectionType.FriendlyName()))
}

func indentLevelString(levels int) string {
	str := ""
	for i := 0; i < levels; i++ {
		str = str + "\t"
	}
	return str
}

func (efr *encodeFnRenderer) GenerateChildrenFuncCalls(indentLevels int, attr string) string {
	indent := indentLevelString(indentLevels)
	lines := make([]string, 0)
	for _, child := range efr.Children {
		if child.Type == generator.FieldTypeAttribute {
			l := fmt.Sprintf("%s%s_%s(%s, ctyVal)", indent, efr.FuncName, child.Name, attr)
			lines = append(lines, l)
		}
		if child.Type == generator.FieldTypeStruct {
			l := fmt.Sprintf("%s%s_%s(%s, ctyVal)", indent, efr.FuncName, child.Name, fmt.Sprintf("%s.%s", attr, child.Name))
			lines = append(lines, l)
		}
	}
	return strings.Join(lines, "\n")
}

var primitiveTypeTemplate = `func {{.FuncName}}(p *{{.ParentType}}, vals map[string]cty.Value) {
	vals["{{.TerraformFieldName}}"] = {{.ConversionFunc}}(p.{{.StructFieldName}})
}`

var primitiveCollectionTypeTemplate = `func {{.FuncName}}(p *{{.ParentType}}, vals map[string]cty.Value) {
	colVals := make([]cty.Value, 0)
	for _, value := range p.{{.StructFieldName}} {
		colVals = append(colVals, {{.ConversionFunc}}(value))
	}
	vals["{{.TerraformFieldName}}"] = {{.CollectionConversionFunc}}(colVals)
}`

var containerTypeTemplate = `func {{.FuncName}}(p *{{.ParentType}}, vals map[string]cty.Value) {
	ctyVal = make(map[string]cty.Value)
{{.GenerateChildrenFuncCalls 1 "p"}}
	vals["{{.TerraformFieldName}}"] = {{.ConversionFunc}}(ctyVal)
}`

var containerCollectionTypeTemplate = `func {{.FuncName}}(p *{{.ParentType}}, vals map[string]cty.Value) {
	valsForCollection = make([]cty.Value, 0)
	for _, v := range p.{{.StructFieldName}} {
		ctyVal = make(map[string]cty.Value)
{{.GenerateChildrenFuncCalls 2 "v"}}
		valsForCollection = append(valsForCollection, {{.ConversionFunc}}(ctyVal))
	}
	vals["{{.TerraformFieldName}}"] = {{.CollectionConversionFunc}}(valsForCollection)
}`

var encoderTemplates = map[string]*template.Template{
	"primitive":           template.Must(template.New("primitive").Parse(primitiveTypeTemplate)),
	"primitiveCollection": template.Must(template.New("primitiveCollection").Parse(primitiveCollectionTypeTemplate)),
	"container":           template.Must(template.New("container").Parse(containerTypeTemplate)),
	"containerCollection": template.Must(template.New("ContainerCollection").Parse(containerCollectionTypeTemplate)),
}

// ImpliedType returns the cty.Type that would result from decoding a
// configuration block using the receiving block schema.
//
// ImpliedType always returns a result, even if the given schema is
// inconsistent.
func impliedType(b *configschema.Block) cty.Type {
	if b == nil {
		return cty.EmptyObject
	}

	atys := make(map[string]cty.Type)

	for name, attrS := range b.Attributes {
		atys[name] = attrS.Type
	}

	for name, blockS := range b.BlockTypes {
		if _, exists := atys[name]; exists {
			panic("invalid schema, blocks and attributes cannot have the same name")
		}

		childType := blockS.Block.ImpliedType()

		switch blockS.Nesting {
		case configschema.NestingSingle, configschema.NestingGroup:
			atys[name] = childType
		case configschema.NestingList:
			// We prefer to use a list where possible, since it makes our
			// implied type more complete, but if there are any
			// dynamically-typed attributes inside we must use a tuple
			// instead, which means our type _constraint_ must be
			// cty.DynamicPseudoType to allow the tuple type to be decided
			// separately for each value.
			if childType.HasDynamicTypes() {
				atys[name] = cty.DynamicPseudoType
			} else {
				atys[name] = cty.List(childType)
			}
		case configschema.NestingSet:
			if childType.HasDynamicTypes() {
				panic("can't use cty.DynamicPseudoType inside a block type with NestingSet")
			}
			atys[name] = cty.Set(childType)
		case configschema.NestingMap:
			// We prefer to use a map where possible, since it makes our
			// implied type more complete, but if there are any
			// dynamically-typed attributes inside we must use an object
			// instead, which means our type _constraint_ must be
			// cty.DynamicPseudoType to allow the tuple type to be decided
			// separately for each value.
			if childType.HasDynamicTypes() {
				atys[name] = cty.DynamicPseudoType
			} else {
				atys[name] = cty.Map(childType)
			}
		default:
			panic("invalid nesting type")
		}
	}

	return cty.Object(atys)
}

var _ generator.EncodeFnRenderer = &backTracker{}
