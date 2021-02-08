package translate

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/crossplane-contrib/terraform-provider-gen/pkg/generator"
	tpl "github.com/crossplane-contrib/terraform-provider-gen/pkg/template"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/zclconf/go-cty/cty"
)

const primitiveTypeDecodeTemplateName = "primitive"
const primitiveCollectionTypeDecodeTemplateName = "primitiveCollection"
const primitiveMapTypeDecodeTemplateName = "primitiveMap"
const containerTypeDecodeTemplateName = "container"
const containerCollectionTypeDecodeTemplateName = "containerCollection"
const containerCollectionSingletonTypeDecodeTemplateName = "containerCollectionSingleton"
const managedResourceDecodeTemplate = "managedResource"

func NewBlockDecodeFnGenerator(terraformName string, block *configschema.NestedBlock) generator.DecodeFnGenerator {
	//ctyType cty.Type, collectionType *cty.Type)
	var colType *cty.Type
	switch block.Nesting {
	case configschema.NestingSingle, configschema.NestingGroup:
		// this is not a collection type, signal that it is a null type
		colType = nil
	case configschema.NestingList:
		colType = &ctyListCollectionType
	case configschema.NestingSet:
		colType = &ctySetCollectionType
	case configschema.NestingMap:
		colType = &ctyMapCollectionType
	default:
		panic("Unrecognized nesting type")
	}
	return &backTracker{
		tfName:         terraformName,
		ctyType:        cty.EmptyObject,
		collectionType: colType,
	}
}

func NewAttributeDecodeFnGenerator(terraformName string, ctyType cty.Type) generator.DecodeFnGenerator {
	if ctyType.IsCollectionType() {
		ct := ctyType.ElementType()
		return &backTracker{
			tfName:         terraformName,
			ctyType:        ct,
			collectionType: &ctyType,
		}
	}
	return &backTracker{
		tfName:  terraformName,
		ctyType: ctyType,
	}
}

// TODO: convert to decode style
func (bt *backTracker) GenerateDecodeFn(funcPrefix, receivedType string, f generator.Field) string {
	efr := bt.decodeFnRenderer(funcPrefix, receivedType, f)
	switch true {
	case bt.ctyType.IsPrimitiveType():
		if bt.collectionType != nil {
			if bt.collectionType.IsMapType() {
				return renderPrimitiveTypeDecoder(efr, primitiveMapTypeTemplateName)
			}
			return renderPrimitiveTypeDecoder(efr, primitiveCollectionTypeTemplateName)
		}
		return renderPrimitiveTypeDecoder(efr, primitiveTypeTemplateName)
	case bt.ctyType.IsMapType() || bt.ctyType.IsObjectType():
		if bt.collectionType != nil {
			if !f.IsSlice {
				return renderContainerTypeDecoder(efr, containerCollectionSingletonTypeTemplateName)
			}
			return renderContainerTypeDecoder(efr, containerCollectionTypeTemplateName)
		}
		return renderContainerTypeDecoder(efr, containerTypeTemplateName)
	default:
		panic(fmt.Sprintf("Unknown cty type in RenderDecodeFn(), cannot render decoder for: %s", bt.ctyType.FriendlyName()))
	}
}

// TODO: convert to decode style
func (bt *backTracker) decodeFnRenderer(funcPrefix, receivedType string, f generator.Field) *decodeFnRenderer {
	return &decodeFnRenderer{
		FuncName:           fmt.Sprintf("%s_%s", funcPrefix, f.Name),
		ParentType:         receivedType,
		TerraformFieldName: bt.tfName,
		StructFieldName:    f.Name,
		Children:           f.Fields,
		CtyType:            bt.ctyType,
		CollectionType:     bt.collectionType,
		Field:              f,
	}
}

type decodeFnRenderer struct {
	FuncName           string // child func name is constructed by appending the child field name
	ParentType         string // name of the parent type to be used in constructing the function receiver
	TerraformFieldName string // the terraform field name will be different from the crd field name and json tag
	StructFieldName    string // name of the child field in the parent struct, to be used for assignment
	Children           []generator.Field
	CtyType            cty.Type
	CollectionType     *cty.Type
	Field              generator.Field
}

// TODO: convert to decode style
func renderPrimitiveTypeDecoder(efr *decodeFnRenderer, template string) string {
	b := bytes.NewBuffer(make([]byte, 0))
	decoderTemplates[template].Execute(b, efr)

	return b.String()
}

// TODO: convert to decode style
func renderContainerTypeDecoder(efr *decodeFnRenderer, template string) string {
	b := bytes.NewBuffer(make([]byte, 0))
	decoderTemplates[template].Execute(b, efr)

	rendered := []string{b.String()}
	for _, child := range efr.Children {
		receivedType := efr.StructFieldName
		if child.Type == generator.FieldTypeStruct {
			receivedType = child.Name
		}
		rendered = append(rendered, child.DecodeFnGenerator.GenerateDecodeFn(efr.FuncName, receivedType, child))
	}
	return strings.Join(rendered, "\n\n")
}

func (efr *decodeFnRenderer) ConversionFunc() string {
	switch efr.CtyType {
	case cty.String:
		return "ctwhy.ValueAsString"
	case cty.Bool:
		return "ctwhy.ValueAsBool"
	case cty.Number:
		return "ctwhy.ValueAsInt64"
	}
	if efr.CtyType.IsObjectType() {
		return "ctwhy.ValueAsObject"
	}
	if efr.CtyType.IsMapType() {
		return "ctwhy.ValueAsMap"
	}
	panic(fmt.Sprintf("Unknown cty type in ConversionFunc(), cannot render convert function for: %s", efr.CtyType.FriendlyName()))
}

func (efr *decodeFnRenderer) CollectionConversionFunc() string {
	if efr.CollectionType.IsSetType() {
		return "ctwhy.ValueAsSet"
	}
	if efr.CollectionType.IsListType() {
		return "ctwhy.ValueAsList"
	}
	panic(fmt.Sprintf("Unknown CollectionType in CollectionConversionFunc(), cannot render convert function for: %s", efr.CollectionType.FriendlyName()))
}

// TODO: convert to decode style
func (efr decodeFnRenderer) GenerateChildrenDecodeFuncCalls(indentLevels int, attr string) string {
	indent := indentLevelString(indentLevels)
	return generateChildrenDecodeFuncCalls(indent, efr.FuncName, attr, efr.Children)
}

func (efr *decodeFnRenderer) PrimitiveFieldType() string {
	return generator.AttributeTypeDeclaration(efr.Field)
}

// TODO: convert to decode style
func generateChildrenDecodeFuncCalls(indent, funcName string, attr string, children []generator.Field) string {
	lines := make([]string, 0)
	for _, child := range children {
		if child.Type == generator.FieldTypeAttribute {
			l := fmt.Sprintf("%s%s_%s(%s, valMap)", indent, funcName, child.Name, attr)
			lines = append(lines, l)
		}
		if child.Type == generator.FieldTypeStruct {
			l := fmt.Sprintf("%s%s_%s(%s, valMap)", indent, funcName, child.Name, fmt.Sprintf("%s.%s", attr, child.Name))
			lines = append(lines, l)
		}
	}
	return strings.Join(lines, "\n")
}

var primitiveTypeDecodeTemplate = `//primitiveTypeDecodeTemplate
func {{.FuncName}}(p *{{.ParentType}}, vals map[string]cty.Value) {
	p.{{.StructFieldName}} = {{.ConversionFunc}}(vals["{{.TerraformFieldName}}"])
}`

var primitiveCollectionTypeDecodeTemplate = `//primitiveCollectionTypeDecodeTemplate
func {{.FuncName}}(p *{{.ParentType}}, vals map[string]cty.Value) {
	goVals := make([]{{ .PrimitiveFieldType }}, 0)
	for _, value := range {{.CollectionConversionFunc}}(vals["{{.TerraformFieldName}}"]) {
		goVals = append(goVals, {{.ConversionFunc}}(value))
	}
	p.{{.StructFieldName}} = goVals
}`

var primitiveMapTypeDecodeTemplate = `//primitiveMapTypeDecodeTemplate
func {{.FuncName}}(p *{{.ParentType}}, vals map[string]cty.Value) {
	// TODO: generalize generation of the element type, string elements are hard-coded atm
	if vals["{{.TerraformFieldName}}"].IsNull() {
		p.{{.StructFieldName}} = nil
        return
    }
	vMap := make(map[string]string)
	v := vals["{{.TerraformFieldName}}"].AsValueMap()
	for key, value := range v {
		vMap[key] = {{.ConversionFunc}}(value)
	}
	p.{{.StructFieldName}} = vMap
}`

var containerTypeDecodeTemplate = `//containerTypeDecodeTemplate
func {{.FuncName}}(p *{{.ParentType}}, vals map[string]cty.Value) {
	valMap := vals["{{.TerraformFieldName}}"].AsValueMap()
{{.GenerateChildrenDecodeFuncCalls 1 "p"}}
}`

// TODO: the collection types are busted
// nesting works a little differently from encode, so these need to a bigger rewrite
// TODO: convert to decode
var containerCollectionTypeDecodeTemplate = `//containerCollectionTypeDecodeTemplate
func {{.FuncName}}(p *[]{{.ParentType}}, vals map[string]cty.Value) {
	valsForCollection := make([]cty.Value, 0)
	for _, v := range p {
		ctyVal := make(map[string]cty.Value)
{{.GenerateChildrenDecodeFuncCalls 2 "v"}}
		valsForCollection = append(valsForCollection, {{.ConversionFunc}}(ctyVal))
	}
	vals["{{.TerraformFieldName}}"] = {{.CollectionConversionFunc}}(valsForCollection)
}`

// TODO: the collection types are busted
// nesting works a little differently from encode, so these need to a bigger rewrite
// TODO: convert to decode
var containerCollectionSingletonTypeDecodeTemplate = `//containerCollectionSingletonTypeDecodeTemplate
func {{.FuncName}}(p *{{.ParentType}}, vals map[string]cty.Value) {
	v := {{.CollectionConversionFunc}}(vals["{{.TerraformFieldName}}"])
	valMap := v[0]
{{.GenerateChildrenDecodeFuncCalls 1 "p"}}
}`

var decodeManagedResourceEntrypointTemplate = `type ctyDecoder struct{}

func (e *ctyDecoder) DecodeCty(mr resource.Managed, ctyValue cty.Value, schema *providers.Schema) (resource.Managed, error) {
	r, ok := mr.(*{{ .TypeName}})
	if !ok {
		return nil, fmt.Errorf("DecodeCty received a resource.Managed value that does not assert to the expected type")
	}
	return {{.DecodeFnName}}(r, ctyValue)
}

func {{.DecodeFnName}}(prev *{{.TypeName}}, ctyValue cty.Value) (resource.Managed, error) {
	valMap := ctyValue.AsValueMap()
	new := prev.DeepCopy()
{{.ForProviderCalls}}
{{.AtProviderCalls}}
	eid := valMap["id"].AsString()
	if len(eid) > 0 {
		meta.SetExternalName(new, eid)
	}
	return new, nil
}`

var decoderTemplates = map[string]*template.Template{
	primitiveTypeTemplateName:                    template.Must(template.New(primitiveTypeDecodeTemplateName).Parse(primitiveTypeDecodeTemplate)),
	primitiveCollectionTypeTemplateName:          template.Must(template.New(primitiveCollectionTypeDecodeTemplateName).Parse(primitiveCollectionTypeDecodeTemplate)),
	primitiveMapTypeTemplateName:                 template.Must(template.New(primitiveMapTypeDecodeTemplateName).Parse(primitiveMapTypeDecodeTemplate)),
	containerTypeTemplateName:                    template.Must(template.New(containerTypeDecodeTemplateName).Parse(containerTypeDecodeTemplate)),
	containerCollectionTypeTemplateName:          template.Must(template.New(containerCollectionTypeDecodeTemplateName).Parse(containerCollectionTypeDecodeTemplate)),
	containerCollectionSingletonTypeTemplateName: template.Must(template.New(containerCollectionSingletonTypeDecodeTemplateName).Parse(containerCollectionSingletonTypeDecodeTemplate)),
	managedResourceTemplate:                      template.Must(template.New(managedResourceDecodeTemplate).Parse(decodeManagedResourceEntrypointTemplate)),
}

var _ generator.DecodeFnGenerator = &backTracker{}

func GenerateDecoders(mr *generator.ManagedResource, tg tpl.TemplateGetter) (string, error) {
	funcName := fmt.Sprintf("Decode%s", mr.Namer().TypeName())
	forProvider := mr.Parameters
	atProvider := mr.Observation
	typeName := mr.Namer().TypeName()

	ttpl, err := tg.Get("hack/template/pkg/generator/decode.go.tmpl")
	if err != nil {
		return "", err
	}

	// TODO: convert forProviderCalls/atProviderCalls to pass values correctly
	forProviderCalls := generateChildrenDecodeFuncCalls("\t", funcName, "&new.Spec.ForProvider", forProvider.Fields)
	atProviderCalls := generateChildrenDecodeFuncCalls("\t", funcName, "&new.Status.AtProvider", atProvider.Fields)

	b := bytes.NewBuffer(make([]byte, 0))
	decoderTemplates[managedResourceTemplate].Execute(b, struct {
		DecodeFnName     string
		TypeName         string
		ForProviderCalls string
		AtProviderCalls  string
	}{
		DecodeFnName:     funcName,
		TypeName:         typeName,
		ForProviderCalls: forProviderCalls,
		AtProviderCalls:  atProviderCalls,
	})
	rendered := []string{b.String()}
	for _, field := range []generator.Field{forProvider, atProvider} {
		for _, child := range field.Fields {
			receivedType := field.Name
			if child.Type == generator.FieldTypeStruct {
				receivedType = child.Name
			}
			rendered = append(rendered, child.DecodeFnGenerator.GenerateDecodeFn(funcName, receivedType, child))
		}
	}

	buf := new(bytes.Buffer)
	tplParams := struct {
		Decoders string
	}{strings.Join(rendered, "\n\n")}
	err = ttpl.Execute(buf, tplParams)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
