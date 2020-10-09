package generator

import (
	"fmt"

	"github.com/pkg/errors"
)

type FieldType int
type AttributeType int

const (
	FieldTypeStruct FieldType = iota
	FieldTypeAttribute
)

const (
	AttributeTypeUintptr AttributeType = iota
	AttributeTypeUint8
	AttributeTypeUint64
	AttributeTypeUint32
	AttributeTypeUint16
	AttributeTypeUint
	AttributeTypeString
	AttributeTypeRune
	AttributeTypeInt8
	AttributeTypeInt64
	AttributeTypeInt32
	AttributeTypeInt16
	AttributeTypeInt
	AttributeTypeFloat64
	AttributeTypeFloat32
	AttributeTypeComplex64
	AttributeTypeComplex128
	AttributeTypeByte
	AttributeTypeBool
)

type StructTag struct {
	JsonName      string
	JsonOmitempty bool
	Inline        bool
}

type Field struct {
	Name           string
	Type           FieldType
	Fields         map[string]Field
	StructField    StructField
	AttributeField AttributeField
	IsSlice        bool
	Tag            StructTag

	// struct comment "annotations"
	Computed  bool
	Optional  bool
	Required  bool
	Sensitive bool
}

type StructField struct {
	PackagePath string
}

type AttributeField struct {
	Type AttributeType
}

type ManagedResource struct {
	Name        string
	PackagePath string
	Spec        ResourceSpec
	Status      ResourceStatus
}

var InvalidMRNameEmpty error = errors.New(".Name is required")
var InvalidMRPackagePathEmpty error = errors.New(".PackagePath is required")

// Validate ensures that the ManagedResource can be rendered to code
func (mr *ManagedResource) Validate() error {
	fail := NewMultiError("ManagedResource.Validate() failed:")
	if mr.Name == "" {
		fail.Append(InvalidMRNameEmpty)
	}
	if mr.PackagePath == "" {
		fail.Append(InvalidMRPackagePathEmpty)
	}

	if len(fail.Errors()) > 0 {
		return fail
	}
	return nil
}

type ResourceSpec struct {
	ForProvider Field
}

func (rs *ResourceSpec) TypeName(resourceName string) string {
	return fmt.Sprintf("%sSpec", resourceName)
}

type ResourceStatus struct {
	AtProvider Field
}

func (rs *ResourceStatus) TypeName(resourceName string) string {
	return fmt.Sprintf("%sStatus", resourceName)
}
