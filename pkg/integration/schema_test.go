// +build integration

package integration

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform/providers"
)

func TestSchemaIterate(t *testing.T) {
	c, err := getProvider()
	if err != nil {
		t.Error(err)
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
