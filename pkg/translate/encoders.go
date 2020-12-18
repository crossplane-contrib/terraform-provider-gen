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

const primitiveTypeTemplateName = "primitive"
const primitiveCollectionTypeTemplateName = "primitiveCollection"
const primitiveMapTypeTemplateName = "primitiveMap"
const containerTypeTemplateName = "container"
const containerCollectionTypeTemplateName = "containerCollection"
const containerCollectionSingletonTypeTemplateName = "containerCollectionSingleton"
const managedResourceTemplate = "managedResource"

func NewBlockEncodeFnGenerator(terraformName string, block *configschema.NestedBlock) generator.EncodeFnGenerator {
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

func NewAttributeEncodeFnGenerator(terraformName string, ctyType cty.Type) generator.EncodeFnGenerator {
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

type backTracker struct {
	tfName         string    // the terraform field name will be different from the crd field name and json tag
	ctyType        cty.Type  // type of field, to be used to select a cty val conversion function (eg cty.StringVal)
	collectionType *cty.Type // if this is a list/set/map then reflect the collection type here for another layer of translation
}

func (bt *backTracker) GenerateEncodeFn(funcPrefix, receivedType string, f generator.Field) string {
	efr := bt.encodeFnRenderer(funcPrefix, receivedType, f)
	switch true {
	case bt.ctyType.IsPrimitiveType():
		if bt.collectionType != nil {
			if bt.collectionType.IsMapType() {
				return renderPrimitiveType(efr, primitiveMapTypeTemplateName)
			}
			return renderPrimitiveType(efr, primitiveCollectionTypeTemplateName)
		}
		return renderPrimitiveType(efr, primitiveTypeTemplateName)
	case bt.ctyType.IsMapType() || bt.ctyType.IsObjectType():
		if bt.collectionType != nil {
			if !f.IsSlice {
				return renderContainerType(efr, containerCollectionSingletonTypeTemplateName)
			}
			return renderContainerType(efr, containerCollectionTypeTemplateName)
		}
		return renderContainerType(efr, containerTypeTemplateName)
	default:
		panic(fmt.Sprintf("Unknown cty type in RenderEncodeFn(), cannot render encoder for: %s", bt.ctyType.FriendlyName()))
	}
}

func (bt *backTracker) encodeFnRenderer(funcPrefix, receivedType string, f generator.Field) *encodeFnRenderer {
	return &encodeFnRenderer{
		FuncName:           fmt.Sprintf("%s_%s", funcPrefix, f.Name),
		ParentType:         receivedType,
		TerraformFieldName: bt.tfName,
		StructFieldName:    f.Name,
		Children:           f.Fields,
		CtyType:            bt.ctyType,
		CollectionType:     bt.collectionType,
	}
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

func renderPrimitiveType(efr *encodeFnRenderer, template string) string {
	b := bytes.NewBuffer(make([]byte, 0))
	encoderTemplates[template].Execute(b, efr)

	return b.String()
}

func renderContainerType(efr *encodeFnRenderer, template string) string {
	b := bytes.NewBuffer(make([]byte, 0))
	encoderTemplates[template].Execute(b, efr)

	rendered := []string{b.String()}
	for _, child := range efr.Children {
		receivedType := efr.StructFieldName
		if child.Type == generator.FieldTypeStruct {
			receivedType = child.Name
		}
		rendered = append(rendered, child.EncodeFnGenerator.GenerateEncodeFn(efr.FuncName, receivedType, child))
	}
	return strings.Join(rendered, "\n\n")
}

func (efr *encodeFnRenderer) ConversionFunc() string {
	switch efr.CtyType {
	case cty.String:
		return "cty.StringVal"
	case cty.Bool:
		return "cty.BoolVal"
	case cty.Number:
		return "cty.NumberIntVal"
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

func (efr *encodeFnRenderer) GenerateChildrenFuncCalls(indentLevels int, attr string) string {
	indent := indentLevelString(indentLevels)
	return generateChildrenFuncCalls(indent, efr.FuncName, attr, efr.Children)
}

func generateChildrenFuncCalls(indent, funcName string, attr string, children []generator.Field) string {
	lines := make([]string, 0)
	for _, child := range children {
		if child.Type == generator.FieldTypeAttribute {
			l := fmt.Sprintf("%s%s_%s(%s, ctyVal)", indent, funcName, child.Name, attr)
			lines = append(lines, l)
		}
		if child.Type == generator.FieldTypeStruct {
			l := fmt.Sprintf("%s%s_%s(%s, ctyVal)", indent, funcName, child.Name, fmt.Sprintf("%s.%s", attr, child.Name))
			lines = append(lines, l)
		}
	}
	return strings.Join(lines, "\n")
}

func indentLevelString(levels int) string {
	str := ""
	for i := 0; i < levels; i++ {
		str = str + "\t"
	}
	return str
}

var primitiveTypeTemplate = `func {{.FuncName}}(p {{.ParentType}}, vals map[string]cty.Value) {
	vals["{{.TerraformFieldName}}"] = {{.ConversionFunc}}(p.{{.StructFieldName}})
}`

var primitiveCollectionTypeTemplate = `func {{.FuncName}}(p {{.ParentType}}, vals map[string]cty.Value) {
	colVals := make([]cty.Value, 0)
	for _, value := range p.{{.StructFieldName}} {
		colVals = append(colVals, {{.ConversionFunc}}(value))
	}
	vals["{{.TerraformFieldName}}"] = {{.CollectionConversionFunc}}(colVals)
}`

var primitiveMapTypeTemplate = `func {{.FuncName}}(p {{.ParentType}}, vals map[string]cty.Value) {
	if len(p.{{.StructFieldName}}) == 0 {
		vals["{{.TerraformFieldName}}"] = cty.NullVal(cty.Map(cty.String))
		return
	}
	mVals := make(map[string]cty.Value)
	for key, value := range p.{{.StructFieldName}} {
		mVals[key] = {{.ConversionFunc}}(value)
	}
	vals["{{.TerraformFieldName}}"] = cty.MapVal(mVals)
}`

var containerTypeTemplate = `func {{.FuncName}}(p {{.ParentType}}, vals map[string]cty.Value) {
	ctyVal := make(map[string]cty.Value)
{{.GenerateChildrenFuncCalls 1 "p"}}
	vals["{{.TerraformFieldName}}"] = {{.ConversionFunc}}(ctyVal)
}`

var containerCollectionTypeTemplate = `func {{.FuncName}}(p []{{.ParentType}}, vals map[string]cty.Value) {
	valsForCollection := make([]cty.Value, 0)
	for _, v := range p {
		ctyVal := make(map[string]cty.Value)
{{.GenerateChildrenFuncCalls 2 "v"}}
		valsForCollection = append(valsForCollection, {{.ConversionFunc}}(ctyVal))
	}
	vals["{{.TerraformFieldName}}"] = {{.CollectionConversionFunc}}(valsForCollection)
}`

var containerCollectionSingletonTypeTemplate = `func {{.FuncName}}(p {{.ParentType}}, vals map[string]cty.Value) {
	valsForCollection := make([]cty.Value, 1)
	ctyVal := make(map[string]cty.Value)
{{.GenerateChildrenFuncCalls 1 "p"}}
	valsForCollection[0] = {{.ConversionFunc}}(ctyVal)
	vals["{{.TerraformFieldName}}"] = {{.CollectionConversionFunc}}(valsForCollection)
}`

var managedResourceEntrypointTemplate = `type ctyEncoder struct{}

func (e *ctyEncoder) EncodeCty(mr resource.Managed, schema *providers.Schema) (cty.Value, error) {
	r, ok := mr.(*{{ .TypeName}})
	if !ok {
		return cty.NilVal, fmt.Errorf("EncodeType received a resource.Managed value which is not a {{ .TypeName}}.")
	}
	return {{.EncodeFnName}}(*r), nil
}

func {{.EncodeFnName}}(r {{.TypeName}}) cty.Value {
	ctyVal := make(map[string]cty.Value)
{{.ForProviderCalls}}
{{.AtProviderCalls}}
	// always set id = external-name if it exists
	// TODO: we should trim Id off schemas in an "optimize" pass
	// before code generation
	en := meta.GetExternalName(&r)
	ctyVal["id"] = cty.StringVal(en)
	return cty.ObjectVal(ctyVal)
}`

var encoderTemplates = map[string]*template.Template{
	primitiveTypeTemplateName:                    template.Must(template.New(primitiveTypeTemplateName).Parse(primitiveTypeTemplate)),
	primitiveCollectionTypeTemplateName:          template.Must(template.New(primitiveCollectionTypeTemplateName).Parse(primitiveCollectionTypeTemplate)),
	primitiveMapTypeTemplateName:                 template.Must(template.New(primitiveMapTypeTemplateName).Parse(primitiveMapTypeTemplate)),
	containerTypeTemplateName:                    template.Must(template.New(containerTypeTemplateName).Parse(containerTypeTemplate)),
	containerCollectionTypeTemplateName:          template.Must(template.New(containerCollectionTypeTemplateName).Parse(containerCollectionTypeTemplate)),
	containerCollectionSingletonTypeTemplateName: template.Must(template.New(containerCollectionSingletonTypeTemplateName).Parse(containerCollectionSingletonTypeTemplate)),
	managedResourceTemplate:                      template.Must(template.New(managedResourceTemplate).Parse(managedResourceEntrypointTemplate)),
}

var _ generator.EncodeFnGenerator = &backTracker{}

func GenerateEncoders(mr *generator.ManagedResource, tg tpl.TemplateGetter) (string, error) {
	funcName := fmt.Sprintf("Encode%s", mr.Namer().TypeName())
	forProvider := mr.Parameters
	atProvider := mr.Observation
	typeName := mr.Namer().TypeName()

	ttpl, err := tg.Get("hack/template/pkg/generator/encode.go.tmpl")
	if err != nil {
		return "", err
	}

	forProviderCalls := generateChildrenFuncCalls("\t", funcName, "r.Spec.ForProvider", forProvider.Fields)
	atProviderCalls := generateChildrenFuncCalls("\t", funcName, "r.Status.AtProvider", atProvider.Fields)

	b := bytes.NewBuffer(make([]byte, 0))
	encoderTemplates[managedResourceTemplate].Execute(b, struct {
		EncodeFnName     string
		TypeName         string
		ForProviderCalls string
		AtProviderCalls  string
	}{
		EncodeFnName:     funcName,
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
			rendered = append(rendered, child.EncodeFnGenerator.GenerateEncodeFn(funcName, receivedType, child))
		}
	}
	buf := new(bytes.Buffer)
	tplParams := struct {
		Encoders string
	}{strings.Join(rendered, "\n\n")}
	err = ttpl.Execute(buf, tplParams)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
