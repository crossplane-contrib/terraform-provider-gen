package generator

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/crossplane/terraform-provider-runtime/pkg/client"
	"github.com/hashicorp/terraform/providers"
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
}
