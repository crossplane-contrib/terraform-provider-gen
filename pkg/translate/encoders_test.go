package translate

import (
	"testing"

	"github.com/crossplane-contrib/terraform-provider-gen/pkg/generator"
	"github.com/zclconf/go-cty/cty"
)

func TestRenderPrimitiveType(t *testing.T) {
	f := generator.Field{
		Name: "SomeAttribute",
	}
	bt := &backTracker{
		TFName:  "some_attribute_tf_name",
		CtyType: cty.String,
	}
	funcPrefix := "encodeResource_Spec_ForProvider"
	receivedType := "ForProvider"

	actual := bt.RenderEncodeFn(funcPrefix, receivedType, f)
	expected := `func encodeResource_Spec_ForProvider_SomeAttribute(p *ForProvider, vals map[string]cty.Value) {
	vals["some_attribute_tf_name"] = cty.StringVal(p.SomeAttribute)
}`
	if actual != expected {
		t.Errorf("Expected:\n----\n%s\n----\nActual:\n----\n%s\n---", expected, actual)
	}
}

func TestRenderPrimitiveCollectionType(t *testing.T) {
	f := generator.Field{
		Name: "SomeAttribute",
	}
	ls := cty.List(cty.String)
	bt := &backTracker{
		TFName:         "some_attribute_tf_name",
		CtyType:        cty.String,
		CollectionType: &ls,
	}
	funcPrefix := "encodeResource_Spec_ForProvider"
	receivedType := "ForProvider"

	actual := bt.RenderEncodeFn(funcPrefix, receivedType, f)
	expected := `func encodeResource_Spec_ForProvider_SomeAttribute(p *ForProvider, vals map[string]cty.Value) {
	colVals := make([]cty.Value, 0)
	for _, value := range p.SomeAttribute {
		colVals = append(colVals, cty.StringVal(value))
	}
	vals["some_attribute_tf_name"] = cty.ListVal(colVals)
}`
	if actual != expected {
		t.Errorf("Expected:\n----\n%s\n----\nActual:\n----\n%s\n---", expected, actual)
	}
}

func TestRenderContainerType(t *testing.T) {
	f := generator.Field{
		Name: "NestedField",
		Fields: []generator.Field{
			{
				Name: "AttributeOne",
				Type: generator.FieldTypeAttribute,
				EncodeFnRenderer: &backTracker{
					TFName:  "attribute_one_tf_name",
					CtyType: cty.String,
				},
			},
			{
				Name: "DeeperField",
				Type: generator.FieldTypeStruct,
				EncodeFnRenderer: &backTracker{
					TFName: "deeper_field_tf_name",
					CtyType: cty.Object(map[string]cty.Type{
						"deeper_attribute_one_tf_name": cty.String,
					}),
				},
				Fields: []generator.Field{
					{
						Name: "DeeperAttributeOne",
						Type: generator.FieldTypeAttribute,
						EncodeFnRenderer: &backTracker{
							TFName:  "deeper_attribute_one_tf_name",
							CtyType: cty.String,
						},
					},
				},
			},
		},
		Type: generator.FieldTypeStruct,
		EncodeFnRenderer: &backTracker{
			TFName: "nested_field_tf_name",
			CtyType: cty.Object(map[string]cty.Type{
				"attribute_one_tf_name": cty.String,
				"deeper_attribute_one_tf_name": cty.Object(map[string]cty.Type{
					"deeper_attribute_one_tf_name": cty.String,
				}),
			}),
		},
	}
	funcPrefix := "encodeResource_Spec_ForProvider"
	receivedType := "NestedField"
	actual := f.EncodeFnRenderer.RenderEncodeFn(funcPrefix, receivedType, f)
	expected := `func encodeResource_Spec_ForProvider_NestedField(p *NestedField, vals map[string]cty.Value) {
	ctyVal = make(map[string]cty.Value)
	encodeResource_Spec_ForProvider_NestedField_AttributeOne(p, ctyVal)
	encodeResource_Spec_ForProvider_NestedField_DeeperField(p.DeeperField, ctyVal)
	vals["nested_field_tf_name"] = cty.ObjectVal(ctyVal)
}

func encodeResource_Spec_ForProvider_NestedField_AttributeOne(p *NestedField, vals map[string]cty.Value) {
	vals["attribute_one_tf_name"] = cty.StringVal(p.AttributeOne)
}

func encodeResource_Spec_ForProvider_NestedField_DeeperField(p *DeeperField, vals map[string]cty.Value) {
	ctyVal = make(map[string]cty.Value)
	encodeResource_Spec_ForProvider_NestedField_DeeperField_DeeperAttributeOne(p, ctyVal)
	vals["deeper_field_tf_name"] = cty.ObjectVal(ctyVal)
}

func encodeResource_Spec_ForProvider_NestedField_DeeperField_DeeperAttributeOne(p *DeeperField, vals map[string]cty.Value) {
	vals["deeper_attribute_one_tf_name"] = cty.StringVal(p.DeeperAttributeOne)
}`
	if actual != expected {
		t.Errorf("Expected:\n----\n%s\n----\nActual:\n----\n%s\n---", expected, actual)
	}
}

func TestRenderContainerCollectionType(t *testing.T) {
	lt := cty.List(cty.EmptyObject)
	f := generator.Field{
		Name: "NestedField",
		Fields: []generator.Field{
			{
				Name: "AttributeOne",
				Type: generator.FieldTypeAttribute,
				EncodeFnRenderer: &backTracker{
					TFName:  "attribute_one_tf_name",
					CtyType: cty.String,
				},
			},
			{
				Name: "DeeperField",
				Type: generator.FieldTypeStruct,
				EncodeFnRenderer: &backTracker{
					TFName:  "deeper_field_tf_name",
					CtyType: cty.EmptyObject,
				},
				Fields: []generator.Field{
					{
						Name: "DeeperAttributeOne",
						Type: generator.FieldTypeAttribute,
						EncodeFnRenderer: &backTracker{
							TFName:  "deeper_attribute_one_tf_name",
							CtyType: cty.String,
						},
					},
				},
			},
		},
		Type: generator.FieldTypeStruct,
		EncodeFnRenderer: &backTracker{
			TFName:         "nested_field_tf_name",
			CtyType:        cty.EmptyObject,
			CollectionType: &lt,
		},
	}
	funcPrefix := "encodeResource_Spec_ForProvider"
	receivedType := "NestedField"
	actual := f.EncodeFnRenderer.RenderEncodeFn(funcPrefix, receivedType, f)
	expected := `func encodeResource_Spec_ForProvider_NestedField(p *NestedField, vals map[string]cty.Value) {
	valsForCollection = make([]cty.Value, 0)
	for _, v := range p.NestedField {
		ctyVal = make(map[string]cty.Value)
		encodeResource_Spec_ForProvider_NestedField_AttributeOne(v, ctyVal)
		encodeResource_Spec_ForProvider_NestedField_DeeperField(v.DeeperField, ctyVal)
		valsForCollection = append(valsForCollection, cty.ObjectVal(ctyVal))
	}
	vals["nested_field_tf_name"] = cty.ListVal(valsForCollection)
}

func encodeResource_Spec_ForProvider_NestedField_AttributeOne(p *NestedField, vals map[string]cty.Value) {
	vals["attribute_one_tf_name"] = cty.StringVal(p.AttributeOne)
}

func encodeResource_Spec_ForProvider_NestedField_DeeperField(p *DeeperField, vals map[string]cty.Value) {
	ctyVal = make(map[string]cty.Value)
	encodeResource_Spec_ForProvider_NestedField_DeeperField_DeeperAttributeOne(p, ctyVal)
	vals["deeper_field_tf_name"] = cty.ObjectVal(ctyVal)
}

func encodeResource_Spec_ForProvider_NestedField_DeeperField_DeeperAttributeOne(p *DeeperField, vals map[string]cty.Value) {
	vals["deeper_attribute_one_tf_name"] = cty.StringVal(p.DeeperAttributeOne)
}`
	if actual != expected {
		t.Errorf("Expected:\n----\n%s\n----\nActual:\n----\n%s\n---", expected, actual)
	}
}
