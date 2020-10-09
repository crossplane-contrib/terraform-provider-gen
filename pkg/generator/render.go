package generator

import (
	j "github.com/dave/jennifer/jen"
)

// RenderManagedResource renders the top-level managed resource object
// including all subtypes.
// TODO: since we need to break out nested structs into their own named structs,
// this method should return a map[string]string where each key is the name of the struct
func RenderManagedResourceStructs(mr *ManagedResource) (string, error) {
	if err := mr.Validate(); err != nil {
		return "", err
	}
	mrStruct := j.Type().Id(mr.Name).Struct(
		j.Qual("metav1", "TypeMeta").Tag(map[string]string{"json": ",inline"}),
		j.Qual("metav1", "ObjectMeta").Tag(map[string]string{"json": "metadata,omitempty"}),
		j.Line(),
		j.Id("Spec").Qual("", mr.Spec.TypeName(mr.Name)).Tag(map[string]string{"json": "spec"}),
		j.Id("Status").Qual("", mr.Status.TypeName(mr.Name)).Tag(map[string]string{"json": "status,omitempty"}),
	)
	return RenderStatement(mrStruct), nil
}

func RenderStatement(s *j.Statement) string {
	return s.GoString()
}
