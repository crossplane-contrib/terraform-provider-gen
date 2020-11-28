package generator

import (
	"strings"

	"github.com/pkg/errors"
)

type FieldType int
type AttributeType int

const (
	FieldTypeStruct FieldType = iota
	FieldTypeAttribute
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=AttributeType,FieldType -output=types_stringers.go
const (
	AttributeTypeUnsupported AttributeType = iota
	AttributeTypeUintptr
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

	// AttributeTypeMapStringKey means that the field is a map[string]<type>
	// where <type> is defined by the value of Field.AttributeField.MapValueType
	AttributeTypeMapStringKey
)

var InvalidMRNameEmpty error = errors.New(".Name is required")
var InvalidMRPackagePathEmpty error = errors.New(".PackagePath is required")

type StructTagJson struct {
	Name      string
	Omitempty bool
	Inline    bool
}

type StructTag struct {
	Json *StructTagJson
}

type Field struct {
	Name             string
	Type             FieldType
	Fields           []Field
	StructField      StructField
	AttributeField   AttributeField
	IsSlice          bool
	Tag              *StructTag
	EncodeFnRenderer EncodeFnRenderer

	// struct comment "annotations"
	Computed  bool
	Optional  bool
	Required  bool
	Sensitive bool
}

type StructField struct {
	PackagePath string
	TypeName    string
}

type AttributeField struct {
	Type         AttributeType
	MapValueType AttributeType
}

type ManagedResource struct {
	Name         string
	PackagePath  string
	Parameters   Field
	Observation  Field
	namer        ResourceNamer
	CategoryTags []string
}

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

// CategoryTagsCSV returns a comma separated list respresenting CategoryTags
// this is used in the kubebuilder resource categories comment annotation
// eg: +kubebuilder:resource:categories={crossplane,managed,aws}
func (mr *ManagedResource) CategoryCSV() string {
	return strings.Join(mr.CategoryTags, ",")
}

func (mr *ManagedResource) Namer() ResourceNamer {
	return mr.namer
}

func (mr *ManagedResource) WithNamer(n ResourceNamer) *ManagedResource {
	mr.namer = n
	return mr
}

func NewManagedResource(name, packagePath string) *ManagedResource {
	return &ManagedResource{
		Name:        name,
		PackagePath: packagePath,
	}
}
