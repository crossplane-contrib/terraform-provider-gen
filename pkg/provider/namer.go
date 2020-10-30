package provider

import (
	"github.com/iancoleman/strcase"
)

type StringTransformer func(string) string

func NewTerraformFieldRenamer() StringTransformer {
	return func(in string) string {
		return strcase.ToCamel(in)
	}
}

type TerraformResourceNamer interface {
	PackageName() string
	ManagedResourceName() string
}

type terraformResourceRenamer struct {
	terraformResourceName string
	strippedResourceName  string
}

func (trr *terraformResourceRenamer) ManagedResourceName() string {
	return strcase.ToCamel(trr.strippedResourceName)
}

func (trr *terraformResourceRenamer) PackageName() string {
	return trr.strippedResourceName
}

func NewTerraformResourceNamer(prefix, tfResourceName string) TerraformResourceNamer {
	var offset int
	if prefix[len(prefix)-1:] == "_" {
		offset = len(prefix)
	} else {
		offset = len(prefix) + 1
	}
	return &terraformResourceRenamer{
		terraformResourceName: tfResourceName,
		strippedResourceName:  tfResourceName[offset:],
	}
}
