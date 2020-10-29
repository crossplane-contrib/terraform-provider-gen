// Code generated by "stringer -type=AttributeType,FieldType -output=types_stringers.go"; DO NOT EDIT.

package generator

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AttributeTypeUintptr-0]
	_ = x[AttributeTypeUint8-1]
	_ = x[AttributeTypeUint64-2]
	_ = x[AttributeTypeUint32-3]
	_ = x[AttributeTypeUint16-4]
	_ = x[AttributeTypeUint-5]
	_ = x[AttributeTypeString-6]
	_ = x[AttributeTypeRune-7]
	_ = x[AttributeTypeInt8-8]
	_ = x[AttributeTypeInt64-9]
	_ = x[AttributeTypeInt32-10]
	_ = x[AttributeTypeInt16-11]
	_ = x[AttributeTypeInt-12]
	_ = x[AttributeTypeFloat64-13]
	_ = x[AttributeTypeFloat32-14]
	_ = x[AttributeTypeComplex64-15]
	_ = x[AttributeTypeComplex128-16]
	_ = x[AttributeTypeByte-17]
	_ = x[AttributeTypeBool-18]
}

const _AttributeType_name = "AttributeTypeUintptrAttributeTypeUint8AttributeTypeUint64AttributeTypeUint32AttributeTypeUint16AttributeTypeUintAttributeTypeStringAttributeTypeRuneAttributeTypeInt8AttributeTypeInt64AttributeTypeInt32AttributeTypeInt16AttributeTypeIntAttributeTypeFloat64AttributeTypeFloat32AttributeTypeComplex64AttributeTypeComplex128AttributeTypeByteAttributeTypeBool"

var _AttributeType_index = [...]uint16{0, 20, 38, 57, 76, 95, 112, 131, 148, 165, 183, 201, 219, 235, 255, 275, 297, 320, 337, 354}

func (i AttributeType) String() string {
	if i < 0 || i >= AttributeType(len(_AttributeType_index)-1) {
		return "AttributeType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _AttributeType_name[_AttributeType_index[i]:_AttributeType_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[FieldTypeStruct-0]
	_ = x[FieldTypeAttribute-1]
}

const _FieldType_name = "FieldTypeStructFieldTypeAttribute"

var _FieldType_index = [...]uint8{0, 15, 33}

func (i FieldType) String() string {
	if i < 0 || i >= FieldType(len(_FieldType_index)-1) {
		return "FieldType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FieldType_name[_FieldType_index[i]:_FieldType_index[i+1]]
}