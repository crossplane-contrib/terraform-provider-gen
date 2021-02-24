package integration

import (
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/generator"
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/optimize"
	"sort"
	"testing"
)

func TestDeduplicator(t *testing.T) {
	mr := DefaultTestResource()
	// TODO: wonky thing that we have to do to satisfy matching package names to exclude
	// the qualifier. Might want to add generator.FakePackagePath as an arg to the fixture instead
	// of assuming it everywhere
	mr.Parameters.StructField.PackagePath = FakePackagePath
	mr.Parameters.Fields = []generator.Field{NestedFieldsWithDuplicates()}
	_, err := optimize.Deduplicate(mr)
	if err != nil {
		t.Errorf("error from optimize.Deduplicate: %s", err)
	}
	fm := make(map[string][]*generator.Field)
	optimize.UnrollFields(&mr.Parameters, fm)
	AssertExistsWithLength(t, fm, "duplicator", 1)
	AssertExistsWithLength(t, fm, "duplicator0", 1)
	AssertExistsWithLength(t, fm, "duplicator1", 1)
}

func TestDeduplicatorIdempotent(t *testing.T) {
	mr := DefaultTestResource()
	// TODO: wonky thing that we have to do to satisfy matching package names to exclude
	// the qualifier. Might want to add generator.FakePackagePath as an arg to the fixture instead
	// of assuming it everywhere
	mr.Parameters.StructField.PackagePath = FakePackagePath
	mr.Parameters.Fields = []generator.Field{NestedFieldsWithDuplicates()}
	_, err := optimize.Deduplicate(mr)
	if err != nil {
		t.Errorf("error from optimize.Deduplicate: %s", err)
	}
	fm := make(map[string][]*generator.Field)
	optimize.UnrollFields(&mr.Parameters, fm)
	AssertExistsWithLength(t, fm, "duplicator", 1)
	AssertExistsWithLength(t, fm, "duplicator0", 1)
	AssertExistsWithLength(t, fm, "duplicator1", 1)
	_, err = optimize.Deduplicate(mr)
	if err != nil {
		t.Errorf("error from optimize.Deduplicate: %s", err)
	}
	fm = make(map[string][]*generator.Field)
	optimize.UnrollFields(&mr.Parameters, fm)
	AssertExistsWithLength(t, fm, "duplicator", 1)
	AssertExistsWithLength(t, fm, "duplicator0", 1)
	AssertExistsWithLength(t, fm, "duplicator1", 1)
}

func AssertExistsWithLength(t *testing.T, fm map[string][]*generator.Field, name string, l int) {
	dup, ok := fm[name]
	if !ok {
		t.Errorf("Expected there to still be a field with TypeName='%s' after de-duplication", name)
	}
	if len(dup) != l {
		t.Errorf("Expected %d field with TypeName='%s' after de-duplication, observed=%d", l, name, len(dup))
	}
}

func TestUnrollFields(t *testing.T) {
	mr := DefaultTestResource()
	// TODO: wonky thing that we have to do to satisfy matching package names to exclude
	// the qualifier. Might want to add generator.FakePackagePath as an arg to the fixture instead
	// of assuming it everywhere
	mr.Parameters.StructField.PackagePath = FakePackagePath
	mr.Parameters.Fields = []generator.Field{NestedFieldsWithDuplicates()}
	fm := make(map[string][]*generator.Field)
	optimize.UnrollFields(&mr.Parameters, fm)
	keys := make([]string, 0)
	for k, _ := range fm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	expected := []string{"duplicator", "middleOne", "middleTwo", "outer"}
	if !compareStringSlices(keys, expected) {
		t.Errorf("Unexpected set of keys after unrolling. expected=%v, actual=%v", expected, keys)
	}
	if len(fm["duplicator"]) != 3 {
		t.Errorf("Expected 3 duplicate fields with struct name = 'duplicator', observed=%d", len(fm["duplicator"]))
	}
}

func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, _ := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
