package provider

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

type StringTransformer func(string) string

func NewTerraformFieldRenamer() StringTransformer {
	return func(in string) string {
		return strcase.ToCamel(in)
	}
}

/*
type StringTransformer func(string) string
	return func(in string) string
ManagedResourceName() string
PackageName() string
APIVersion() string
strippedResourceName() string
TypeName() string
*/

type TerraformResourceNamer interface {
	PackageName() string
	ManagedResourceName() string
	APIVersion() string
	APIGroup() string
	KubernetesVersion() string
	TypeNameGroupKind() string
	TypeNameGroupVersionKind() string
	TerraformResourceName() string
}

type terraformResourceRenamer struct {
	terraformResourceName string
	apiVersion            string
	providerName          string
}

func (trr *terraformResourceRenamer) ManagedResourceName() string {
	return strcase.ToCamel(trr.strippedResourceName())
}
func (trr *terraformResourceRenamer) ManagedResourceListName() string {
	return fmt.Sprintf("%sList", trr.ManagedResourceName())
}

func (trr *terraformResourceRenamer) PackageName() string {
	return trr.strippedResourceName()
}

func (trr *terraformResourceRenamer) APIVersion() string {
	return trr.apiVersion
}

func (trr *terraformResourceRenamer) strippedResourceName() string {
	prefix := trr.providerName
	var offset int
	if prefix[len(prefix)-1:] == "_" {
		offset = len(prefix)
	} else {
		offset = len(prefix) + 1
	}
	return trr.terraformResourceName[offset:]
}

func NewTerraformResourceNamer(providerName, tfResourceName, apiVersion string) TerraformResourceNamer {
	return &terraformResourceRenamer{
		terraformResourceName: tfResourceName,
		apiVersion:            apiVersion,
		providerName:          providerName,
	}
}

func (trr *terraformResourceRenamer) APIGroup() string {
	return fmt.Sprintf("%s.terraform-provider-%s.crossplane.io",
		strcase.ToKebab(trr.PackageName()), trr.providerName)
}

func (trr *terraformResourceRenamer) TypeName() string {
	return trr.ManagedResourceName()
}

func (trr *terraformResourceRenamer) TerraformResourceName() string {
	return trr.terraformResourceName
}

// KubernetesVersion is an alias to .APIVersion
// TODO: this exists because some of the templates started using
// KubernetesVersion and I haven't made up my mind as to whether I want to change it
func (trr *terraformResourceRenamer) KubernetesVersion() string {
	return trr.APIVersion()
}

func (trr *terraformResourceRenamer) TypeNameGroupKind() string {
	return fmt.Sprintf("%sGroupKind", trr.ManagedResourceName())
}

func (trr *terraformResourceRenamer) TypeNameGroupVersionKind() string {
	return fmt.Sprintf("%sGroupVersionKind", trr.ManagedResourceName())
}

/*
< package {{ .KubernetesVersion}}
< 	name := managed.ControllerName({{ .TypeNameGroupKind }})
< 		resource.ManagedKind({{ .TypeNameGroupVersionKind }}),
< 		For(&{{ .ManagedResourceName }}{}).
*/
