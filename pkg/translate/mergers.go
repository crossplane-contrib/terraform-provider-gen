package translate

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/crossplane-contrib/terraform-provider-gen/pkg/generator"
	tpl "github.com/crossplane-contrib/terraform-provider-gen/pkg/template"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/zclconf/go-cty/cty"
)

const mergePrimitiveTemplateName = "primitive"
const mergePrimitiveContainerTemplateName = "primitiveContainer"
const mergeStructTemplateName = "struct"
const mergeStructSliceTemplateName = "structContainer"

func NewBlockMergeFnGenerator(terraformName string, block *configschema.NestedBlock) generator.MergeFnGenerator {
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

func NewAttributeMergeFnGenerator(terraformName string, ctyType cty.Type) generator.MergeFnGenerator {
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

func (bt *backTracker) GenerateMergeFn(funcPrefix, receivedType string, f generator.Field, spec bool) string {
	efr := bt.mergeFnRenderer(funcPrefix, receivedType, f)
	switch true {
	case f.Type == generator.FieldTypeAttribute:
		if f.AttributeField.Type == generator.AttributeTypeMapStringKey {
			return renderPrimitiveTypeMerger(efr, mergePrimitiveContainerTemplateName, spec)
		}
		if f.IsSlice {
			return renderPrimitiveTypeMerger(efr, mergePrimitiveContainerTemplateName, spec)
		}
		return renderPrimitiveTypeMerger(efr, mergePrimitiveTemplateName, spec)
	case f.Type == generator.FieldTypeStruct:
		if f.IsSlice {
			return renderContainerTypeMerger(efr, mergeStructSliceTemplateName, spec)
		}
		return renderContainerTypeMerger(efr, mergeStructTemplateName, spec)
	default:
		panic(fmt.Sprintf("Unknown cty type in RenderMergeFn(), cannot render merger for: %s", bt.ctyType.FriendlyName()))
	}
}

func (bt *backTracker) mergeFnRenderer(funcPrefix, receivedType string, f generator.Field) *mergeFnRenderer {
	return &mergeFnRenderer{
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

type mergeFnRenderer struct {
	FuncName           string // child func name is constructed by appending the child field name
	ParentType         string // name of the parent type to be used in constructing the function receiver
	TerraformFieldName string // the terraform field name will be different from the crd field name and json tag
	StructFieldName    string // name of the child field in the parent struct, to be used for assignment
	Children           []generator.Field
	CtyType            cty.Type
	CollectionType     *cty.Type
	Field              generator.Field
}

func renderPrimitiveTypeMerger(efr *mergeFnRenderer, template string, isSpec bool) string {
	b := bytes.NewBuffer(make([]byte, 0))
	if isSpec {
		specTemplates[template].Execute(b, efr)
	} else {
		statusTemplates[template].Execute(b, efr)
	}
	return b.String()
}

func renderContainerTypeMerger(efr *mergeFnRenderer, template string, isSpec bool) string {
	b := bytes.NewBuffer(make([]byte, 0))
	if isSpec {
		specTemplates[template].Execute(b, efr)
	} else {
		statusTemplates[template].Execute(b, efr)
	}

	rendered := []string{b.String()}
	sort.Stable(generator.NamedFields(efr.Children))
	for _, child := range efr.Children {
		receivedType := efr.ParentType
		if child.Type == generator.FieldTypeStruct {
			receivedType = child.StructField.TypeName
		}
		rendered = append(rendered, child.MergeFnGenerator.GenerateMergeFn(efr.FuncName, receivedType, child, isSpec))
	}
	return strings.Join(rendered, "\n\n")
}

func (efr *mergeFnRenderer) PrimitiveFieldType() string {
	return generator.AttributeTypeDeclaration(efr.Field)
}

var updatedCallHandlingTmpl *template.Template = template.Must(template.New("updatedCallHandlingTmpl ").Parse(`{{ .ChildCall }}
{{ .Indent }}if updated {
{{ .Indent }}	anyChildUpdated = true
{{ .Indent }}}
`))

func renderUpdatedHandling(call, indent string) string {
	b := bytes.NewBuffer(make([]byte, 0))
	updatedCallHandlingTmpl.Execute(b, struct {
		Indent    string
		ChildCall string
	}{
		Indent:    indent,
		ChildCall: call,
	})
	return b.String()
}

func (efr mergeFnRenderer) GenerateChildrenMergeFuncCalls(indentLevels int, isSpec bool) string {
	indent := indentLevelString(indentLevels)
	return generateChildrenMergeFuncCalls(indent, efr.FuncName, efr.Children, isSpec, "k", "p", false)
}

// TODO: this needs a better design, the attrRefs hack is shameful
// the issue is that we always receive pointers to nested functions, so in the case of a simple attribute where we pass
// down the pointer, we can pass it along directly to the receiver which already wants a pointer. when the child is a struct
// we want to pass down a reference to the struct instead. but there is a special case in the top-level merge entrypoint
// function where we want to pass down a reference to a nested struct member, kind of like when the child is a struct,
// which is why we have the attrRefs hack -- we need to pass these attributes down as refs because they aren't already
// this could be fixed by rewriting the entrypoint
func generateChildrenMergeFuncCalls(indent, funcName string, children []generator.Field, isSpec bool, leftAttr string, rightAttr string, attrRefs bool) string {
	lines := make([]string, 0)
	sort.Stable(generator.NamedFields(children))
	for _, child := range children {
		if child.Type == generator.FieldTypeAttribute {
			l := fmt.Sprintf("%supdated = %s_%s(%s, %s, md)", indent, funcName, child.Name, leftAttr, rightAttr)
			if attrRefs {
				l = fmt.Sprintf("%supdated = %s_%s(&%s, &%s, md)", indent, funcName, child.Name, leftAttr, rightAttr)
			}
			lines = append(lines, renderUpdatedHandling(l, indent))
		}
		if child.Type == generator.FieldTypeStruct {
			l := fmt.Sprintf("%supdated = %s_%s(&%s.%s, &%s.%s, md)", indent, funcName, child.Name, leftAttr, child.Name, rightAttr, child.Name)
			lines = append(lines, renderUpdatedHandling(l, indent))
		}
	}
	return strings.Join(lines, "\n")
}

func (efr *mergeFnRenderer) PrimitiveContainerComparison() string {
	f := efr.Field
	switch f.AttributeField.Type {
	case generator.AttributeTypeMapStringKey:
		switch f.AttributeField.MapValueType {
		case generator.AttributeTypeBool:
			return "plugin.CompareMapBool"
		case generator.AttributeTypeString:
			return "plugin.CompareMapString"
		case generator.AttributeTypeInt:
			return "plugin.CompareMapInt64"
		}
	case generator.AttributeTypeInt64:
		if !f.IsSlice {
			panic(fmt.Sprintf("Attribute treated as container but is not a slice or map %v", f.Name))
		}
		return "plugin.CompareInt64Slices"
	case generator.AttributeTypeString:
		if !f.IsSlice {
			panic(fmt.Sprintf("Attribute treated as container but is not a slice or map %v", f.Name))
		}
		return "plugin.CompareStringSlices"
	}
	panic(fmt.Sprintf("Attribute treated as container but does not match supported cases: name=%s, type=%s, isSlice=%t", f.Name, f.AttributeField.Type.String(), f.IsSlice))
}

var lateInitializePrimitiveTemplate = `//lateInitializePrimitiveTemplate
func {{.FuncName}}(k *{{.ParentType}}, p *{{.ParentType}}, md *plugin.MergeDescription) bool {
	if *p.{{ .StructFieldName }} != *k.{{ .StructFieldName }} {
		if k.{{ .StructFieldName }} == nil {
			*k.{{ .StructFieldName }} = *p.{{ .StructFieldName }}
			md.LateInitializedSpec = true
		} else {
			*p.{{ .StructFieldName }} = *k.{{ .StructFieldName }}
			md.NeedsProviderUpdate = true
		}
		return true
	}
	return false
}`

var mergePrimitiveTemplateStatus = `//mergePrimitiveTemplateStatus
func {{.FuncName}}(k *{{.ParentType}}, p *{{.ParentType}}, md *plugin.MergeDescription) bool {
	if k.{{ .StructFieldName }} != p.{{ .StructFieldName }} {
		k.{{ .StructFieldName }} = p.{{ .StructFieldName }}
		md.StatusUpdated = true
		return true
	}
	return false
}`

var mergePrimitiveTemplateSpec = `//mergePrimitiveTemplateSpec
func {{.FuncName}}(k *{{.ParentType}}, p *{{.ParentType}}, md *plugin.MergeDescription) bool {
	if k.{{ .StructFieldName }} != p.{{ .StructFieldName }} {
		p.{{ .StructFieldName }} = k.{{ .StructFieldName }}
		md.NeedsProviderUpdate = true
		return true
	}
	return false
}`

var mergePrimitiveContainerTemplateStatus = `//mergePrimitiveContainerTemplateStatus
func {{.FuncName}}(k *{{.ParentType}}, p *{{.ParentType}}, md *plugin.MergeDescription) bool {
	if !{{.PrimitiveContainerComparison }}(k.{{ .StructFieldName }}, p.{{ .StructFieldName }}) {
		k.{{ .StructFieldName }} = p.{{ .StructFieldName }}
		md.StatusUpdated = true
		return true
	}
	return false
}`

var mergePrimitiveContainerTemplateSpec = `//mergePrimitiveContainerTemplateSpec
func {{.FuncName}}(k *{{.ParentType}}, p *{{.ParentType}}, md *plugin.MergeDescription) bool {
	if !{{.PrimitiveContainerComparison }}(k.{{ .StructFieldName }}, p.{{ .StructFieldName }}) {
		p.{{ .StructFieldName }} = k.{{ .StructFieldName }}
		md.NeedsProviderUpdate = true
		return true
	}
	return false
}`

var mergeStructTemplateStatus = `//mergeStructTemplateStatus
func {{.FuncName}}(k *{{.ParentType}}, p *{{.ParentType}}, md *plugin.MergeDescription) bool {
	updated := false
	anyChildUpdated := false
{{.GenerateChildrenMergeFuncCalls 1 false}}
	if anyChildUpdated {
		md.StatusUpdated = true
	}
	return anyChildUpdated
}`

var mergeStructTemplateSpec = `//mergeStructTemplateSpec
func {{.FuncName}}(k *{{.ParentType}}, p *{{.ParentType}}, md *plugin.MergeDescription) bool {
	updated := false
	anyChildUpdated := false
{{.GenerateChildrenMergeFuncCalls 1 true}}
	if anyChildUpdated {
		md.NeedsProviderUpdate = true
	}
	return anyChildUpdated
}`

var mergeStructSliceTemplateStatus = `//mergeStructSliceTemplateStatus
func {{.FuncName}}(ksp *[]{{.ParentType}}, psp *[]{{.ParentType}}, md *plugin.MergeDescription) bool {
	if len(*ksp) != len(*psp) {
		*ksp = *psp
		md.NeedsProviderUpdate = true
		return true
	}
	ks := *ksp
	ps := *psp
	anyChildUpdated := false
	for i := range ps {
		updated := false
		k := &ks[i]
		p := &ps[i]
{{.GenerateChildrenMergeFuncCalls 2 false}}
	}
	if anyChildUpdated {
		md.StatusUpdated = true
	}
	return anyChildUpdated
}`

var mergeStructSliceTemplateSpec = `//mergeStructSliceTemplateSpec
func {{.FuncName}}(ksp *[]{{.ParentType}}, psp *[]{{.ParentType}}, md *plugin.MergeDescription) bool {
	if len(*ksp) != len(*psp) {
		*psp = *ksp
		md.NeedsProviderUpdate = true
		return true
	}
	ks := *ksp
	ps := *psp
	anyChildUpdated := false
	for i := range ps {
		updated := false
		k := &ks[i]
		p := &ps[i]
{{.GenerateChildrenMergeFuncCalls 2 true }}
	}
	if anyChildUpdated {
		md.NeedsProviderUpdate = true
	}
	return anyChildUpdated
}`

var mergeManagedResourceEntrypointTemplate = `//mergeManagedResourceEntrypointTemplate
type resourceMerger struct{}

func (r *resourceMerger) MergeResources(kube resource.Managed, prov resource.Managed) plugin.MergeDescription {
	k := kube.(*{{ .TypeName }})
	p := prov.(*{{ .TypeName }})
	md := &plugin.MergeDescription{}
	updated := false
	anyChildUpdated := false

{{.ForProviderCalls}}
{{.AtProviderCalls}}
	for key, v := range p.Annotations {
		if k.Annotations[key] != v {
			k.Annotations[key] = v
			md.AnnotationsUpdated = true
		}
	}
	md.AnyFieldUpdated = anyChildUpdated
	return *md
}`

var specTemplates = map[string]*template.Template{
	mergePrimitiveTemplateName:          template.Must(template.New(mergePrimitiveTemplateName).Parse(mergePrimitiveTemplateSpec)),
	mergePrimitiveContainerTemplateName: template.Must(template.New(mergePrimitiveContainerTemplateName).Parse(mergePrimitiveContainerTemplateSpec)),
	mergeStructTemplateName:             template.Must(template.New(mergeStructTemplateName).Parse(mergeStructTemplateSpec)),
	mergeStructSliceTemplateName:        template.Must(template.New(mergeStructSliceTemplateName).Parse(mergeStructSliceTemplateSpec)),
}

var statusTemplates = map[string]*template.Template{
	mergePrimitiveTemplateName:          template.Must(template.New(mergePrimitiveTemplateName).Parse(mergePrimitiveTemplateStatus)),
	mergePrimitiveContainerTemplateName: template.Must(template.New(mergePrimitiveContainerTemplateName).Parse(mergePrimitiveContainerTemplateStatus)),
	mergeStructTemplateName:             template.Must(template.New(mergeStructTemplateName).Parse(mergeStructTemplateStatus)),
	mergeStructSliceTemplateName:        template.Must(template.New(mergeStructSliceTemplateName).Parse(mergeStructSliceTemplateStatus)),
}

var _ generator.MergeFnGenerator = &backTracker{}

func GenerateMergers(mr *generator.ManagedResource, tg tpl.TemplateGetter) (string, error) {
	funcName := fmt.Sprintf("Merge%s", mr.Namer().TypeName())
	forProvider := mr.Parameters
	atProvider := mr.Observation
	typeName := mr.Namer().TypeName()

	ttpl, err := tg.Get("pkg/generator/compare.go.tmpl")
	if err != nil {
		return "", err
	}

	forProviderCalls := generateChildrenMergeFuncCalls("\t", funcName, forProvider.Fields, true, "k.Spec.ForProvider", "p.Spec.ForProvider", true)
	atProviderCalls := generateChildrenMergeFuncCalls("\t", funcName, atProvider.Fields, false, "k.Status.AtProvider", "p.Status.AtProvider", true)

	b := bytes.NewBuffer(make([]byte, 0))
	tmpl := template.Must(template.New("mrtpl").Parse(mergeManagedResourceEntrypointTemplate))
	tmpl.Execute(b, struct {
		TypeName         string
		ForProviderCalls string
		AtProviderCalls  string
	}{
		TypeName:         typeName,
		ForProviderCalls: forProviderCalls,
		AtProviderCalls:  atProviderCalls,
	})
	rendered := []string{b.String()}
	for _, field := range []generator.Field{forProvider} {
		for _, child := range field.Fields {
			receivedType := field.Name
			if child.Type == generator.FieldTypeStruct {
				receivedType = child.StructField.TypeName
			}
			rendered = append(rendered, child.MergeFnGenerator.GenerateMergeFn(funcName, receivedType, child, true))
		}
	}
	for _, field := range []generator.Field{atProvider} {
		for _, child := range field.Fields {
			receivedType := field.Name
			if child.Type == generator.FieldTypeStruct {
				receivedType = child.StructField.TypeName
			}
			rendered = append(rendered, child.MergeFnGenerator.GenerateMergeFn(funcName, receivedType, child, false))
		}
	}
	buf := new(bytes.Buffer)
	tplParams := struct {
		Mergers string
	}{strings.Join(rendered, "\n\n")}
	err = ttpl.Execute(buf, tplParams)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
