package translate

import "github.com/iancoleman/strcase"

type NameBeautifier struct {
	transformers []StringTransformer
}

type StringTransformer func(string) string

type terraformProviderTruncater struct {
	leftOffset int
}

func (t *terraformProviderTruncater) Transform(in string) string {
	return in[t.leftOffset:]
}

func NewTerraformResourceRenamer(prefix string) StringTransformer {
	var offset int
	if prefix[len(prefix)-1:] == "_" {
		offset = len(prefix)
	} else {
		offset = len(prefix) + 1
	}
	return func(in string) string {
		return strcase.ToCamel(in[offset:])
	}
}

func NewTerraformFieldRenamer() StringTransformer {
	return func(in string) string {
		return strcase.ToCamel(in)
	}
}
