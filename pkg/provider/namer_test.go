package provider

import "testing"

func TestTerraformTypeRenamer(t *testing.T) {
	tfName := "aws_resource"
	expected := "Resource"
	prefix := "aws_"
	r := NewTerraformResourceNamer(prefix, tfName)
	actual := r.ManagedResourceName()
	if actual != expected {
		t.Errorf("Unexpected renaming of '%s' to '%s' using NewTerraformResourceRenamer('%s'). expected=%s", tfName, actual, prefix, expected)
	}
	prefix = "aws"
	r = NewTerraformResourceNamer(prefix, tfName)
	actual = r.ManagedResourceName()
	if actual != expected {
		t.Errorf("Unexpected renaming of '%s' to '%s' using NewTerraformResourceRenamer('%s'). expected=%s", tfName, actual, prefix, expected)
	}

	tfName = "aws_longer_resource_name"
	expected = "LongerResourceName"
	prefix = "aws_"
	r = NewTerraformResourceNamer(prefix, tfName)
	actual = r.ManagedResourceName()
	if actual != expected {
		t.Errorf("Unexpected renaming of '%s' to '%s' using NewTerraformResourceRenamer('%s'). expected=%s", tfName, actual, prefix, expected)
	}
}

func TestTerraformFieldRenamer(t *testing.T) {
	field_name := "meandering_long_field_name"
	expected := "MeanderingLongFieldName"
	r := NewTerraformFieldRenamer()
	actual := r(field_name)
	if actual != expected {
		t.Errorf("Unexpected renaming of '%s' to '%s' using NewTerraformFieldRenamer. expected=%s", field_name, actual, expected)
	}
}
