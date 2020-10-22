package generator

import (
	"fmt"
	"testing"

	j "github.com/dave/jennifer/jen"
)

func expectValidationError(expectedErr, actualErr error, t *testing.T) {
	if actualErr != nil {
		me, ok := actualErr.(MultiError)
		if !ok {
			t.Errorf("Could not type assert Validate error to MultiError")
		}
		for _, e := range me.Errors() {
			if e == expectedErr {
				return
			}
		}
	}
	t.Errorf("Did not find expected validation error=%v", expectedErr)
}

func TestBaseWithValidation(t *testing.T) {
	resourceName := "Test"
	mr := &ManagedResource{}
	mr.WithNamer(NewDefaultNamer(resourceName))
	mrs, err := RenderManagedResourceTypes(mr)
	expectValidationError(InvalidMRNameEmpty, err, t)
	expectValidationError(InvalidMRPackagePathEmpty, err, t)

	mr.Name = resourceName
	mr.PackagePath = "github.com/crossplane-contrib/fake"
	mrs, err = RenderManagedResourceTypes(mr)
	if err != nil {
		t.Errorf("Unexpected error from RenderManagedResourceTypes: %v", err)
	}

	expected := "// +kubebuilder:object:root=true\n" +
		"\n" +
		"// Test is a managed resource representing a resource mirrored in the cloud\n" +
		"// +kubebuilder:subresource:status\n" +
		"// +kubebuilder:resource:scope=Cluster\n" +
		"type Test struct {\n" +
		"	metav1.TypeMeta   `json:\",inline\"`\n" +
		"	metav1.ObjectMeta `json:\"metadata,omitempty\"`\n" +
		"\n" +
		"	Spec   TestSpec   `json:\"spec\"`\n" +
		"	Status TestStatus `json:\"status,omitempty\"`\n" +
		"}"

	actual := mrs[mr.Namer().TypeName()]
	if actual != expected {
		t.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}
}

func TestDefaultIsValid(t *testing.T) {
	mr := defaultTestResource()
	_, err := RenderManagedResourceTypes(mr)
	if err != nil {
		t.Errorf("Unexpected error from RenderManagedResourceTypes: %v", err)
	}
}

func assertRenderExpected(t *testing.T, mr *ManagedResource, typeName string, expected string) error {
	// checking this error is handled by TestDefaultIsValid
	mrs, _ := RenderManagedResourceTypes(mr)
	actual, ok := mrs[typeName]
	if !ok {
		return fmt.Errorf("Could not find TypeListName()=%s in output from RenderManagedResourceTypes", typeName)
	}
	if actual != expected {
		return fmt.Errorf("Unexpected output from jen render for %s.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", typeName, expected, actual)
	}
	return nil
}

func TestResourceList(t *testing.T) {
	mr := defaultTestResource()
	expected := "// +kubebuilder:object:root=true\n" +
		"\n" +
		"// Test contains a list of TestList\n" +
		"type TestList struct {\n" +
		"	metav1.TypeMeta `json:\",inline\"`\n" +
		"	metav1.ListMeta `json:\"metadata,omitempty\"`\n" +
		"	Items           []Test `json:\"items\"`\n" +
		"}"
	err := assertRenderExpected(t, mr, mr.Namer().TypeListName(), expected)
	if err != nil {
		t.Error(err)
	}
}

func TestSpec(t *testing.T) {
	mr := defaultTestResource()
	expected := "// A TestSpec defines the desired state of a Test\n" +
		"type TestSpec struct {\n" +
		"	runtimev1alpha1.ResourceSpec `json:\",inline\"`\n" +
		"	ForProvider                  TestParameters `json:\",inline\"`\n" +
		"}"
	err := assertRenderExpected(t, mr, mr.Namer().SpecTypeName(), expected)
	if err != nil {
		t.Error(err)
	}
}

func TestStatus(t *testing.T) {
	mr := defaultTestResource()
	expected := "// A TestStatus defines the observed state of a Test\n" +
		"type TestStatus struct {\n" +
		"	runtimev1alpha1.ResourceStatus `json:\",inline\"`\n" +
		"	AtProvider                     TestObservation `json:\",inline\"`\n" +
		"}"
	err := assertRenderExpected(t, mr, mr.Namer().StatusTypeName(), expected)
	if err != nil {
		t.Error(err)
	}
}

func TestAttributeStatementStruct(t *testing.T) {
	parent := Field{
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: fakePackagePath,
		},
	}
	// TODO: should we even allow this in our managed resource types?
	// this case is an unqualified, embedded struct
	test := Field{
		Name: "Fieldname",
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: fakePackagePath,
			PackageName: "PackageName",
		},
	}
	stmt := j.Type().Id("fakeStruct").Struct(AttributeStatement(test, parent))
	actual := renderStatement(stmt)
	expected := "type fakeStruct struct {\n" +
		"	Fieldname PackageName\n" +
		"}"
	if actual != expected {
		t.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}

	// modify the package name different from the parent to get a package qualifier
	test.StructField.PackagePath = test.StructField.PackagePath + "/different"
	stmt = j.Type().Id("fakeStruct").Struct(AttributeStatement(test, parent))
	actual = renderStatement(stmt)
	// this is because jen chops the qualifier down to the imported -package
	expected = "type fakeStruct struct {\n" +
		"	Fieldname different.PackageName\n" +
		"}"
	if actual != expected {
		t.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}

	// TODO: should we support this? probably not, & disallowing it could reduce possible errors
	// remove the field name to get an embedded type
	test.Name = ""
	stmt = j.Type().Id("fakeStruct").Struct(AttributeStatement(test, parent))
	actual = renderStatement(stmt)
	expected = "type fakeStruct struct {\n" +
		"	different.PackageName\n" +
		"}"
	if actual != expected {
		t.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}
}

type fieldOpt func(f *Field)

func withFieldTag(t StructTag) fieldOpt {
	return func(f *Field) {
		f.Tag = &t
	}
}

func ezAttrField(name string, attrType AttributeType, opts ...fieldOpt) Field {
	f := Field{
		Name: name,
		Type: FieldTypeAttribute,
		AttributeField: AttributeField{
			Type: attrType,
		},
	}
	for _, opt := range opts {
		opt(&f)
	}
	return f
}

func ezStructField(name, pkgPath, pkgName string, opts ...fieldOpt) Field {
	f := Field{
		Name: name,
		Type: FieldTypeStruct,
		StructField: StructField{
			PackageName: pkgName,
			PackagePath: pkgPath,
		},
	}
	for _, opt := range opts {
		opt(&f)
	}
	return f
}

func TestAttributeStatementPrimitives(t *testing.T) {
	test := Field{
		Name: fakeResourceName,
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: fakePackagePath,
			PackageName: "PackageName",
		},
		Fields: []Field{
			ezAttrField("a", AttributeTypeUintptr),
			ezAttrField("b", AttributeTypeUint8),
			ezAttrField("c", AttributeTypeUint64),
			ezAttrField("d", AttributeTypeUint32),
			ezAttrField("e", AttributeTypeUint16),
			ezAttrField("f", AttributeTypeUint),
			ezAttrField("g", AttributeTypeString),
			ezAttrField("h", AttributeTypeRune),
			ezAttrField("i", AttributeTypeInt8),
			ezAttrField("j", AttributeTypeInt64),
			ezAttrField("k", AttributeTypeInt32),
			ezAttrField("l", AttributeTypeInt16),
			ezAttrField("m", AttributeTypeInt),
			ezAttrField("n", AttributeTypeFloat64),
			ezAttrField("o", AttributeTypeFloat32),
			ezAttrField("p", AttributeTypeComplex64),
			ezAttrField("q", AttributeTypeComplex128),
			ezAttrField("r", AttributeTypeByte),
			ezAttrField("s", AttributeTypeBool),
		},
	}

	frags := FieldFragments(test)
	actual := frags[fakeResourceName].Render()
	expected := "type Test struct {\n" +
		"	a uintptr\n" +
		"	b uint8\n" +
		"	c uint64\n" +
		"	d uint32\n" +
		"	e uint16\n" +
		"	f uint\n" +
		"	g string\n" +
		"	h rune\n" +
		"	i int8\n" +
		"	j int64\n" +
		"	k int32\n" +
		"	l int16\n" +
		"	m int\n" +
		"	n float64\n" +
		"	o float32\n" +
		"	p complex64\n" +
		"	q complex128\n" +
		"	r byte\n" +
		"	s bool\n" +
		"}"
	if actual != expected {
		t.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}
}

func testTag(t *testing.T, tagJson *StructTagJson, expected string) error {
	tag := StructTag{
		Json: tagJson,
	}
	test := Field{
		Name: fakeResourceName,
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: fakePackagePath,
			PackageName: "PackageName",
		},
		Fields: []Field{
			ezAttrField("a", AttributeTypeUintptr, withFieldTag(tag)),
		},
	}
	frags := FieldFragments(test)
	actual := frags[fakeResourceName].Render()
	if actual != expected {
		return fmt.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}
	return nil
}

func TestAttrFieldTags(t *testing.T) {
	var tj *StructTagJson
	var expected string

	tj = &StructTagJson{
		Name:      "",
		Omitempty: true,
	}
	expected = "type Test struct {\n" +
		"	a uintptr `json:\",omitempty\"`\n" +
		"}"
	if err := testTag(t, tj, expected); err != nil {
		t.Error(err)
	}

	tj = &StructTagJson{
		Name:   "",
		Inline: true,
	}
	expected = "type Test struct {\n" +
		"	a uintptr `json:\",inline\"`\n" +
		"}"
	if err := testTag(t, tj, expected); err != nil {
		t.Error(err)
	}

	tj = &StructTagJson{
		Name: "a_field",
	}
	expected = "type Test struct {\n" +
		"	a uintptr `json:\"a_field\"`\n" +
		"}"
	if err := testTag(t, tj, expected); err != nil {
		t.Error(err)
	}

	tj = &StructTagJson{
		Name:   "a_field",
		Inline: true,
	}
	expected = "type Test struct {\n" +
		"	a uintptr `json:\"a_field,inline\"`\n" +
		"}"
	if err := testTag(t, tj, expected); err != nil {
		t.Error(err)
	}

	tj = &StructTagJson{
		Name:      "a_field",
		Omitempty: true,
	}
	expected = "type Test struct {\n" +
		"	a uintptr `json:\"a_field,omitempty\"`\n" +
		"}"
	if err := testTag(t, tj, expected); err != nil {
		t.Error(err)
	}
}

func testStructFieldTag(t *testing.T, field Field, expected string) error {
	test := Field{
		Name: fakeResourceName,
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: fakePackagePath,
			PackageName: "PackageName",
		},
		Fields: []Field{
			field,
		},
	}
	frags := FieldFragments(test)
	actual := frags[fakeResourceName].Render()
	if actual != expected {
		return fmt.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}
	return nil
}

func TestStructFieldTag(t *testing.T) {
	f := Field{
		Name: "SubField",
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: fakePackagePath,
			PackageName: "AnotherName",
		},
		Tag: &StructTag{
			Json: &StructTagJson{
				Name: "sub_field",
			},
		},
	}

	// we won't repeat testing all the permutations of the json tag fields
	// as in TestAttrFieldTags on the assumption that the code path of tag
	// construction is otherwise the same, just want to exercise the struct
	// path of AttributeStatement
	expected := "type Test struct {\n" +
		"	SubField AnotherName `json:\"sub_field\"`\n" +
		"}"
	if err := testStructFieldTag(t, f, expected); err != nil {
		t.Error(err)
	}
}

func TestFieldFragmentsNested(t *testing.T) {
	deeplyNestedTypeName := "DeeplyNestedField"
	nestedTypeName := "NestedField"
	test := nestedFieldFixture(nestedTypeName, deeplyNestedTypeName)
	frags := FieldFragments(test)
	actual := frags[fakeResourceName].Render()
	expected := "type Test struct {\n" +
		"	NestedField `json:\"sub_field\"`\n" +
		"}"
	if actual != expected {
		t.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}

	expected = "type NestedField struct {\n" +
		"	DeeplyNestedField `json:\"deeper_sub_field\"`\n" +
		"}"
	actual = frags[nestedTypeName].Render()
	if actual != expected {
		t.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}
}
