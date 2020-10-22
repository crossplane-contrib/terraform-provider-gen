package generator

import (
	"bytes"
	"fmt"

	"github.com/crossplane/terraform-provider-gen/pkg/template"
	j "github.com/dave/jennifer/jen"
)

const CommentBlankLine = "" // assumes comments are joined with newlines
const KubebuilderObjectRoot = "+kubebuilder:object:root=true"
const KubebuilderMarkStatusSubresource = "+kubebuilder:subresource:status"

// RenderKubebuilderResourceAnnotation renderes the kubebuilder resource tag
// which indicates whether the resources is namespace- or cluster-scoped
// and sets the categories field in the CRD which allow kubectl to select
// sets of resources based on matching category tags.
func RenderKubebuilderResourceAnnotation(mr *ManagedResource) string {
	catCSV := mr.CategoryCSV()
	if catCSV == "" {
		return "+kubebuilder:resource:scope=Cluster"
	}
	return fmt.Sprintf(" +kubebuilder:resource:scope=Cluster,categories={%s}", catCSV)
}

func ResourceTypeFragment(mr *ManagedResource) *Fragment {
	namer := mr.Namer()
	mrStruct := j.Type().Id(namer.TypeName()).Struct(
		j.Qual("metav1", "TypeMeta").Tag(map[string]string{"json": ",inline"}),
		j.Qual("metav1", "ObjectMeta").Tag(map[string]string{"json": "metadata,omitempty"}),
		j.Line(),
		j.Id("Spec").Qual("", namer.SpecTypeName()).Tag(map[string]string{"json": "spec"}),
		j.Id("Status").Qual("", namer.StatusTypeName()).Tag(map[string]string{"json": "status,omitempty"}),
	)

	comments := []string{
		KubebuilderObjectRoot,
		CommentBlankLine,
		// TODO: check if there is a reasonable description field we can use for this comment
		fmt.Sprintf("%s is a managed resource representing a resource mirrored in the cloud", namer.TypeName()),
		// TODO: handle printcolumn lines
		// we always mark ou resources
		KubebuilderMarkStatusSubresource,
		RenderKubebuilderResourceAnnotation(mr),
	}

	return &Fragment{
		comments:  comments,
		statement: mrStruct,
	}
}

func TypeListFragment(mr *ManagedResource) *Fragment {
	namer := mr.Namer()
	stmt := j.Type().Id(namer.TypeListName()).Struct(
		j.Qual("metav1", "TypeMeta").Tag(map[string]string{"json": ",inline"}),
		j.Qual("metav1", "ListMeta").Tag(map[string]string{"json": "metadata,omitempty"}),
		j.Id("Items").Index().Qual("", namer.TypeName()).Tag(map[string]string{"json": "items"}),
	)
	comments := []string{
		KubebuilderObjectRoot,
		CommentBlankLine,
		fmt.Sprintf("%s contains a list of %s", namer.TypeName(), namer.TypeListName()),
	}
	return &Fragment{
		statement: stmt,
		comments:  comments,
	}
}

func SpecFragment(mr *ManagedResource) *Fragment {
	namer := mr.Namer()
	stmt := j.Type().Id(namer.SpecTypeName()).Struct(
		j.Qual("runtimev1alpha1", "ResourceSpec").Tag(map[string]string{"json": ",inline"}),
		j.Id("ForProvider").Qual("", namer.ForProviderTypeName()).Tag(map[string]string{"json": ",inline"}),
	)
	comment := fmt.Sprintf("A %s defines the desired state of a %s", namer.SpecTypeName(), namer.TypeName())
	return &Fragment{
		name:      namer.SpecTypeName(),
		statement: stmt,
		comments:  []string{comment},
	}
}

func ForProviderFragments(mr *ManagedResource) []*Fragment {
	namer := mr.Namer()
	if mr.Parameters.StructField.TypeName != namer.ForProviderTypeName() {
		mr.Parameters.StructField.TypeName = namer.ForProviderTypeName()
	}
	frags := FieldFragments(mr.Parameters)
	// frags[0] is the outermost element, aka ForProvider
	frags[0].comments = []string{
		fmt.Sprintf("A %s defines the desired state of a %s", namer.ForProviderTypeName(), namer.TypeName()),
	}
	return frags
}

func StatusFragments(mr *ManagedResource) []*Fragment {
	namer := mr.Namer()
	stmt := j.Type().Id(namer.StatusTypeName()).Struct(
		j.Qual("runtimev1alpha1", "ResourceStatus").Tag(map[string]string{"json": ",inline"}),
		j.Id("AtProvider").Qual("", namer.AtProviderTypeName()).Tag(map[string]string{"json": ",inline"}),
	)
	comment := fmt.Sprintf("A %s defines the observed state of a %s", namer.StatusTypeName(), namer.TypeName())
	return []*Fragment{{
		name:      namer.StatusTypeName(),
		statement: stmt,
		comments:  []string{comment},
	}}
}

func FieldFragments(f Field) []*Fragment {
	attributes := make([]j.Code, 0)
	nested := make([]*Fragment, 0)
	for _, a := range f.Fields {
		attributes = append(attributes, AttributeStatement(a, f))
		if a.Type == FieldTypeStruct {
			for _, frag := range FieldFragments(a) {
				nested = append(nested, frag)
			}
		}
	}
	// append the nested fields onto the tail so that the results are
	// in recursive-descent order.
	return append([]*Fragment{{
		name:      f.Name,
		statement: j.Type().Id(f.StructField.TypeName).Struct(attributes...),
	}}, nested...)
}

func AttributeStatement(f, parent Field) *j.Statement {
	id := j.Id(f.Name)
	if f.IsSlice {
		id = id.Index()
	}
	switch f.Type {
	case FieldTypeAttribute:
		id = TypeStatement(f, id)
	case FieldTypeStruct:
		// TODO: since you can have an embedded struct, we need to allow
		// the name to be excluded, and since we can have relative packages
		// package path can be empty, but we should always have a TypeName
		path := f.StructField.PackagePath
		// this check kind of assumes that we don't refer to types that claim to
		// be in a different package and also nest other types. probably a safe
		// assumption since nesting should only be for types nested within
		// the same terraform package
		if f.StructField.PackagePath == parent.StructField.PackagePath {
			path = ""
		}
		id = id.Qual(path, f.StructField.TypeName)
	}
	if f.Tag != nil {
		if f.Tag.Json != nil {
			jsonTag := ""
			if f.Tag.Json.Name != "" {
				jsonTag = f.Tag.Json.Name
			}
			if f.Tag.Json.Inline {
				jsonTag = jsonTag + ",inline"
			}
			if f.Tag.Json.Omitempty {
				jsonTag = jsonTag + ",omitempty"
			}
			id.Tag(map[string]string{"json": jsonTag})
		}
	}
	return id
}

func TypeStatement(f Field, s *j.Statement) *j.Statement {
	switch f.AttributeField.Type {
	case AttributeTypeUintptr:
		return s.Uintptr()
	case AttributeTypeUint8:
		return s.Uint8()
	case AttributeTypeUint64:
		return s.Uint64()
	case AttributeTypeUint32:
		return s.Uint32()
	case AttributeTypeUint16:
		return s.Uint16()
	case AttributeTypeUint:
		return s.Uint()
	case AttributeTypeString:
		return s.String()
	case AttributeTypeRune:
		return s.Rune()
	case AttributeTypeInt8:
		return s.Int8()
	case AttributeTypeInt64:
		return s.Int64()
	case AttributeTypeInt32:
		return s.Int32()
	case AttributeTypeInt16:
		return s.Int16()
	case AttributeTypeInt:
		return s.Int()
	case AttributeTypeFloat64:
		return s.Float64()
	case AttributeTypeFloat32:
		return s.Float32()
	case AttributeTypeComplex64:
		return s.Complex64()
	case AttributeTypeComplex128:
		return s.Complex128()
	case AttributeTypeByte:
		return s.Byte()
	case AttributeTypeBool:
		return s.Bool()
	}

	panic(fmt.Sprintf("Unable to determine type for %s", f.Name))
}

type managedResourceTypeDefRenderer struct {
	mr *ManagedResource
	tg template.TemplateGetter
}

func NewManagedResourceTypeDefRenderer(mr *ManagedResource, tg template.TemplateGetter) *managedResourceTypeDefRenderer {
	return &managedResourceTypeDefRenderer{
		mr: mr,
		tg: tg,
	}
}

func (tdr *managedResourceTypeDefRenderer) Render() (string, error) {
	mr := tdr.mr
	if err := mr.Validate(); err != nil {
		return "", err
	}
	typeDefs := make([]*Fragment, 0)

	typeDefs = append(typeDefs, ResourceTypeFragment(mr))
	typeDefs = append(typeDefs, TypeListFragment(mr))
	typeDefs = append(typeDefs, SpecFragment(mr))
	for _, frag := range ForProviderFragments(mr) {
		typeDefs = append(typeDefs, frag)
	}

	for _, frag := range StatusFragments(mr) {
		typeDefs = append(typeDefs, frag)
	}

	tpl, err := tdr.tg.Get("hack/template/pkg/generator/types.go.tmpl")
	if err != nil {
		return "", err
	}

	typeDefsString := ""
	for _, f := range typeDefs {
		typeDefsString = fmt.Sprintf("%s\n\n%s", typeDefsString, f.Render())
	}

	buf := new(bytes.Buffer)
	tplParams := struct {
		TypeDefs string
	}{typeDefsString}
	err = tpl.Execute(buf, tplParams)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
