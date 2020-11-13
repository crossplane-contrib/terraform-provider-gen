package generator

import (
	"fmt"
	"testing"

	j "github.com/dave/jennifer/jen"
)

const FakeResourceName string = "Test"
const FakePackagePath string = "github.com/crossplane-contrib/fake"

func DefaultTestResource() *ManagedResource {
	return NewManagedResource(FakeResourceName, FakePackagePath).WithNamer(NewDefaultNamer(FakeResourceName))
}

func NestedFieldFixture(outerTypeName, nestedTypeName, deeplyNestedTypeName string) Field {
	f := Field{
		// "Name" is appended to help visually differentiate field and type names
		Name: deeplyNestedTypeName + "Name",
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: FakePackagePath,
			TypeName:    deeplyNestedTypeName,
		},
		Fields: []Field{
			{
				Name:           "aString",
				Type:           FieldTypeAttribute,
				AttributeField: AttributeField{Type: AttributeTypeString},
				Tag: &StructTag{
					Json: &StructTagJson{
						Name: "a_string",
					},
				},
			},
		},
		Tag: &StructTag{
			Json: &StructTagJson{
				Name: "deeper_sub_field",
			},
		},
	}
	nf := Field{
		// "Name" is appended to help visually differentiate field and type names
		Name: nestedTypeName + "Name",
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: FakePackagePath,
			TypeName:    nestedTypeName,
		},
		Fields: []Field{
			f,
		},
		Tag: &StructTag{
			Json: &StructTagJson{
				Name: "sub_field",
			},
		},
	}
	test := Field{
		// "Name" is appended to help visually differentiate field and type names
		Name: outerTypeName + "Name",
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: FakePackagePath,
			TypeName:    outerTypeName,
		},
		Fields: []Field{
			nf,
		},
	}
	return test
}

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
	err := mr.Validate()
	expectValidationError(InvalidMRNameEmpty, err, t)
	expectValidationError(InvalidMRPackagePathEmpty, err, t)

	mr.Name = resourceName
	mr.PackagePath = "github.com/crossplane-contrib/fake"
	actual := ResourceTypeFragment(mr).Render()

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

	if actual != expected {
		t.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}
}

func TestDefaultIsValid(t *testing.T) {
	mr := DefaultTestResource()
	err := mr.Validate()
	if err != nil {
		t.Errorf("Unexpected error from ManagedResource.Validate(): %v", err)
	}
}

func TestResourceList(t *testing.T) {
	mr := DefaultTestResource()
	expected := "// +kubebuilder:object:root=true\n" +
		"\n" +
		"// Test contains a list of TestList\n" +
		"type TestList struct {\n" +
		"	metav1.TypeMeta `json:\",inline\"`\n" +
		"	metav1.ListMeta `json:\"metadata,omitempty\"`\n" +
		"	Items           []Test `json:\"items\"`\n" +
		"}"
	actual := TypeListFragment(mr).Render()
	if actual != expected {
		t.Errorf("Unexpected output from TypeListFragment.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}
}

func TestSpec(t *testing.T) {
	mr := DefaultTestResource()
	expected := "// A TestSpec defines the desired state of a Test\n" +
		"type TestSpec struct {\n" +
		"	runtimev1alpha1.ResourceSpec `json:\",inline\"`\n" +
		"	ForProvider                  TestParameters `json:\",inline\"`\n" +
		"}"
	actual := SpecFragment(mr).Render()
	if actual != expected {
		t.Errorf("Unexpected value for SpecFragment.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}
}

func TestStatus(t *testing.T) {
	mr := DefaultTestResource()
	expected := "// A TestStatus defines the observed state of a Test\n" +
		"type TestStatus struct {\n" +
		"	runtimev1alpha1.ResourceStatus `json:\",inline\"`\n" +
		"	AtProvider                     TestObservation `json:\",inline\"`\n" +
		"}"
	actual := StatusFragment(mr).Render()
	if actual != expected {
		t.Errorf("Unexpected value for StatusFragment.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}
}

func TestAttributeStatementStruct(t *testing.T) {
	parent := Field{
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: FakePackagePath,
		},
	}
	// TODO: should we even allow this in our managed resource types?
	// this case is an unqualified, embedded struct
	test := Field{
		Name: "Fieldname",
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: FakePackagePath,
			TypeName:    "TypeName",
		},
	}
	stmt := j.Type().Id("fakeStruct").Struct(AttributeStatement(test, parent))
	actual := renderStatement(stmt)
	expected := "type fakeStruct struct {\n" +
		"	Fieldname TypeName\n" +
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
		"	Fieldname different.TypeName\n" +
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
		"	different.TypeName\n" +
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
			TypeName:    pkgName,
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
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: FakePackagePath,
			TypeName:    FakeResourceName,
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
	if len(frags) != 1 {
		t.Errorf("Expected %d results from FieldFragments, saw %d", 1, len(frags))
	}
	actual := frags[0].Render() // assumes single result from FieldFragments
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
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: FakePackagePath,
			TypeName:    FakeResourceName,
		},
		Fields: []Field{
			ezAttrField("a", AttributeTypeUintptr, withFieldTag(tag)),
		},
	}
	frags := FieldFragments(test)
	if len(frags) != 1 {
		t.Errorf("Expected %d results from FieldFragments, saw %d", 1, len(frags))
	}
	actual := frags[0].Render()
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
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: FakePackagePath,
			TypeName:    FakeResourceName,
		},
		Fields: []Field{
			field,
		},
	}
	frags := FieldFragments(test)
	actual := frags[0].Render()
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
			PackagePath: FakePackagePath,
			TypeName:    "AnotherName",
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

// TestFieldFragmentsNested assumes elements are ordered from outside-in
// this ordering assumption is important because it translates to how
// we want things to look in the final output
func TestFieldFragmentsNested(t *testing.T) {
	deeplyNestedTypeName := "DeeplyNestedField"
	nestedTypeName := "NestedField"
	test := NestedFieldFixture(FakeResourceName, nestedTypeName, deeplyNestedTypeName)
	frags := FieldFragments(test)
	if len(frags) != 3 {
		t.Errorf("Expected %d results from FieldFragments, saw %d", 3, len(frags))
	}
	actual := frags[0].Render()
	expected := "type Test struct {\n" +
		"	NestedFieldName NestedField `json:\"sub_field\"`\n" +
		"}"
	if actual != expected {
		t.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}

	expected = "type NestedField struct {\n" +
		"	DeeplyNestedFieldName DeeplyNestedField `json:\"deeper_sub_field\"`\n" +
		"}"
	actual = frags[1].Render()
	if actual != expected {
		t.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}

	expected = "type DeeplyNestedField struct {\n" +
		"	aString string `json:\"a_string\"`\n" +
		"}"
	actual = frags[2].Render()
	if actual != expected {
		t.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}
}

func TestMapFragment(t *testing.T) {
	f := Field{
		// "Name" is appended to help visually differentiate field and type names
		Name: FakeResourceName + "Name",
		Type: FieldTypeStruct,
		StructField: StructField{
			PackagePath: FakePackagePath,
			TypeName:    FakeResourceName,
		},
		Fields: []Field{
			{
				Name: "StrMap",
				Type: FieldTypeAttribute,
				AttributeField: AttributeField{
					Type:         AttributeTypeMapStringKey,
					MapValueType: AttributeTypeString,
				},
			},
			{
				Name: "BoolMap",
				Type: FieldTypeAttribute,
				AttributeField: AttributeField{
					Type:         AttributeTypeMapStringKey,
					MapValueType: AttributeTypeBool,
				},
			},
		},
	}
	frags := FieldFragments(f)
	actual := frags[0].Render()
	expected := "type Test struct {\n" +
		"	StrMap  map[string]string\n" +
		"	BoolMap map[string]bool\n" +
		"}"

	if expected != actual {
		t.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, actual)
	}
}
