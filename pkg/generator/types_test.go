package generator

import (
	"testing"
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

func TestBaseStruct(t *testing.T) {
	mr := &ManagedResource{}
	mrs, err := RenderManagedResourceStructs(mr)
	expectValidationError(InvalidMRNameEmpty, err, t)
	expectValidationError(InvalidMRPackagePathEmpty, err, t)

	mr.Name = "Test"
	mr.PackagePath = "github.com/crossplane-contrib/fake"
	mrs, err = RenderManagedResourceStructs(mr)
	if err != nil {
		t.Errorf("Unexpected error from RenderManagedResourceStructs: %v", err)
	}

	expected := "type Test struct {\n" +
		"	metav1.TypeMeta   `json:\",inline\"`\n" +
		"	metav1.ObjectMeta `json:\"metadata,omitempty\"`\n" +
		"\n" +
		"	Spec   TestSpec   `json:\"spec\"`\n" +
		"	Status TestStatus `json:\"status,omitempty\"`\n" +
		"}"
	if mrs != expected {
		t.Errorf("Unexpected output from jen render.\nExpected:\n ---- \n%s\n ---- \nActual:\n%s", expected, mrs)
	}
}
