package generator

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/crossplane-contrib/terraform-runtime/pkg/client"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/providers"
	"github.com/zclconf/go-cty/cty"
)

const (
	//awsRelativePluginPath string = "testdata/.terraform/plugins/registry.terraform.io/hashicorp/aws/3.9.0/darwin_amd64/terraform-provider-aws_v3.9.0_x5"
	awsRelativePluginPath string = "testdata/.terraform/plugins/registry.terraform.io/hashicorp/aws/3.9.0/darwin_amd64/"
	awsProviderName       string = "aws"
)

func pluginPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(wd, awsRelativePluginPath), nil
}

func TestSchemaIterate(t *testing.T) {
	pluginPath, err := pluginPath()
	if err != nil {
		t.Fatal(err)
	}
	c, err := client.NewGRPCProvider(awsProviderName, pluginPath)
	if err != nil {
		t.Fatal(err)
	}
	vpcSchema := make(map[string]providers.Schema)
	vpcKeys := make([]string, 0)
	for name, s := range c.GetSchema().ResourceTypes {
		if strings.Contains(name, "aws_vpc") {
			vpcSchema[name] = s
			vpcKeys = append(vpcKeys, name)
		}
	}

	if _, ok := vpcSchema["aws_vpc"]; !ok {
		t.Errorf("Could not find vpc network in list of resource schemas! keys containing vpc=%v", vpcKeys)
	}

	vpc := vpcSchema["aws_vpc"]
	t.Logf("%v", vpc)
}

func TestFlatBlock(t *testing.T) {
	s := providers.Schema{
		Block: &configschema.Block{
			Attributes: make(map[string]*configschema.Attribute),
			BlockTypes: make(map[string]*configschema.NestedBlock),
		},
	}
	// I think "id" should probably not be part of the schema, it is like our external-name
	// TODO: check how this was implemented in the prototype
	//s.Block.Attributes["id"] =
	s.Block.Attributes["different_resource_ref_id"] = &configschema.Attribute{
		Required: false,
		Optional: true,
		Computed: false,
		Type:     cty.String,
	}
	s.Block.Attributes["perform_optional_action"] = &configschema.Attribute{
		Required: false,
		Optional: true,
		Computed: false,
		Type:     cty.Bool,
	}
	s.Block.Attributes["labels"] = &configschema.Attribute{
		Required: false,
		Optional: true,
		Computed: false,
		Type:     cty.Map(cty.String),
	}
	s.Block.Attributes["number_list"] = &configschema.Attribute{
		Required: false,
		Optional: true,
		Computed: false,
		Type:     cty.List(cty.Number),
	}
	s.Block.Attributes["computed_owner_id"] = &configschema.Attribute{
		Required: false,
		Optional: false,
		Computed: true,
		Type:     cty.String,
	}
	s.Block.Attributes["required_name"] = &configschema.Attribute{
		Required: true,
		Optional: false,
		Computed: false,
		Type:     cty.String,
	}
	t.Logf("%v", s)
}
