package v1alpha1

import (
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	xpresource "github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/hashicorp/terraform/providers"
	"github.com/zclconf/go-cty/cty"
	"sigs.k8s.io/yaml"
)

func UnmarshalServiceAccount(data []byte) (resource.Managed, error) {
	var sa resource.Managed = &ServiceAccount{}
	err := yaml.Unmarshal(data, sa)

	return sa, err
}

func AsCtyValue(resource xpresource.Managed, schema *providers.Schema) (cty.Value, error) {
	sa, ok := resource.(*ServiceAccount)
	if !ok {
		return cty.NilVal, fmt.Errorf("iam.AsCtyValue received a resource.Managed value which is not a ServiceAccount.")
	}
	ctyVal := make(map[string]cty.Value)
	//id := meta.GetExternalName(resource)
	//ctyVal["id"] = cty.StringVal(id)
	ctyVal["id"] = cty.StringVal(sa.Status.AtProvider.Name)
	ctyVal["account_id"] = cty.StringVal(*sa.Spec.ForProvider.DisplayName)
	ctyVal["display_name"] = cty.StringVal(*sa.Spec.ForProvider.DisplayName)
	ctyVal["description"] = cty.StringVal(*sa.Spec.ForProvider.Description)
	ctyVal["name"] = cty.StringVal(sa.Status.AtProvider.Name)
	ctyVal["project"] = cty.StringVal(sa.Status.AtProvider.Project)
	ctyVal["unique_id"] = cty.StringVal(sa.Status.AtProvider.UniqueID)
	ctyVal["email"] = cty.StringVal(sa.Status.AtProvider.Email)
	timeouts := make(map[string]cty.Value)
	timeouts["create"] = cty.StringVal("")
	ctyVal["timeouts"] = cty.ObjectVal(timeouts)
	//= cty.StringVal(*sa.TerraformConfig.Timeouts.Create)
	return cty.ObjectVal(ctyVal), nil
	//return unrollAndLookup(sa, schema.Block, "")
	/*
		if err != nil {
			return cty.NilVal, err
		}
		return cty.ObjectVal(ctyValue), nil
	*/
}

func FromCtyValue(previousManaged resource.Managed, ctyValue cty.Value, schema *providers.Schema) (resource.Managed, error) {
	prev, ok := previousManaged.(*ServiceAccount)
	if !ok {
		return nil, fmt.Errorf("iam.AsCtyValue received a resource.Managed value for previousManaged which is not a ServiceAccount")
	}
	valMap := ctyValue.AsValueMap()
	new := prev.DeepCopy()
	new.Status.AtProvider.Name = valMap["name"].AsString()
	new.Status.AtProvider.Project = valMap["project"].AsString()
	new.Status.AtProvider.UniqueID = valMap["unique_id"].AsString()
	new.Status.AtProvider.Email = valMap["email"].AsString()
	new.Spec.ForProvider.AccountID = valMap["account_id"].AsString()
	new.Spec.ForProvider.Description = stringReference(valMap["description"].AsString())
	new.Spec.ForProvider.DisplayName = stringReference(valMap["display_name"].AsString())
	return new, nil
}

func AsYAML(res resource.Managed) ([]byte, error) {
	sa, ok := res.(*ServiceAccount)
	if !ok {
		return nil, fmt.Errorf("iam.AsYAML received a resource.Managed value which is not a ServiceAccount.")
	}
	return yaml.Marshal(sa)
}

/*
func encodeValue(path string, sa *ServiceAccount) (cty.Value, error) {
	switch path {
	case "account_id":
		fmt.Printf("%s\n", *sa.Spec.ForProvider.DisplayName)
		return cty.StringVal(*sa.Spec.ForProvider.DisplayName), nil
	case "display_name":
		fmt.Printf("%s\n", *sa.Spec.ForProvider.DisplayName)
		return cty.StringVal(*sa.Spec.ForProvider.DisplayName), nil
	case "description":
		fmt.Printf("%s\n", *sa.Spec.ForProvider.Description)
		return cty.StringVal(*sa.Spec.ForProvider.Description), nil
	case "name":
		fmt.Printf("%s\n", sa.Status.AtProvider.Name)
		return cty.StringVal(sa.Status.AtProvider.Name), nil
	case "project":
		fmt.Printf("%s\n", sa.Status.AtProvider.Project)
		return cty.StringVal(sa.Status.AtProvider.Project), nil
	case "unique_id":
		fmt.Printf("%s\n", sa.Status.AtProvider.UniqueID)
		return cty.StringVal(sa.Status.AtProvider.UniqueID), nil
	case "email":
		fmt.Printf("%s\n", sa.Status.AtProvider.Email)
		return cty.StringVal(sa.Status.AtProvider.Email), nil
	case "timeouts.create":
		fmt.Printf("%s\n", *sa.TerraformConfig.Timeouts.Create)
		return cty.StringVal(*sa.TerraformConfig.Timeouts.Create), nil
	}
	return cty.NilVal, fmt.Errorf("Generated code does not understand path=%s", path)
}

func unrollAndLookupNested(sa *ServiceAccount, block *configschema.NestedBlock, path string) (cty.Value, error) {
	switch block.Nesting {
	case configschema.NestingGroup:
		return cty.NilVal, fmt.Errorf("not yet implemented")
	}
}

func unrollAndLookup(sa *ServiceAccount, block *configschema.Block, path string) (cty.Value, error) {
	ctyValue := make(map[string]cty.Value)
	for name, _ := range block.Attributes {
		// doing this to skip over `id`...
		_, exists := ctyValue[name]
		if exists {
			continue
		}
		if path != "" {
			name = path + "." + name
		}
		v, err := encodeValue(name, sa)
		if err != nil {
			return cty.NilVal, err
		}
		ctyValue[name] = v
	}
	for name, b := range block.BlockTypes {
		subpath := strings.Join([]string{path, name}, ".")
		fmt.Printf("subpath=%s\n", subpath)
		cv, err := unrollAndLookup(sa, &b.Block, subpath)
		if err != nil {
			return cty.NilVal, err
		}
		ctyValue[name] = cv
	}
}
*/

// StringValue converts the supplied string pointer to a string, returning the
// empty string if the pointer is nil.
func stringReference(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
