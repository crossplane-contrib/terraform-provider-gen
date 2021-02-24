package optimize

import (
	"fmt"
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/generator"
)

func Deduplicate(mr *generator.ManagedResource) (*generator.ManagedResource, error) {
	fm := make(map[string][]*generator.Field)
	UnrollFields(&mr.Parameters, fm)
	UnrollFields(&mr.Observation, fm)
	for k, f := range fm {
		if k == "" {
			continue
		}
		if len(f) > 1 {
			for i, _ := range f[1:] {
				o := i + 1
				collision := f[o]
				collision.StructField.TypeName = fmt.Sprintf("%s%d", collision.StructField.TypeName, i)
			}
		}
	}
	return mr, nil
}

// Build a map of Fields, keyed by Field.StructField.TypeName
// this is used to by the Deduplicate optimizer find sets of duplicately named fields.
// It is exported as a public function so that it can be used in the integration package to help verify
// correct behavior of Deduplicate.
func UnrollFields(fld *generator.Field, fm map[string][]*generator.Field) {
	if fld.StructField.TypeName != "" {
		fs, ok := fm[fld.StructField.TypeName]
		if !ok {
			fs = make([]*generator.Field, 0)
		}
		fs = append(fs, fld)
		fm[fld.StructField.TypeName] = fs
	}
	for i, _ := range fld.Fields {
		f := &fld.Fields[i]
		if f.Type == generator.FieldTypeStruct {
			UnrollFields(f, fm)
		}
	}
}

var _ Optimizer = Deduplicate
