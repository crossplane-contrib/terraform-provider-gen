// +build integration

package integration

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform/providers"
	"github.com/zclconf/go-cty/cty"
)

func TestSchemaIterate(t *testing.T) {
	c, err := getProvider(&IntegrationTestConfig{})
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
}

func TestNetworkACL(t *testing.T) {
	c, err := getProvider(&IntegrationTestConfig{})
	if err != nil {
		t.Error(err)
	}
	rt := c.GetSchema().ResourceTypes
	attrs := rt["aws_network_acl"].Block.Attributes
	var ingressType cty.Type = attrs["ingress"].Type
	if !ingressType.IsSetType() {
		t.Errorf("Expected aws_network_acl.ingress to be a set of objects")
	}
	var ingElemsType cty.Type = ingressType.ElementType()
	if !ingElemsType.IsObjectType() {
		t.Errorf("Expected aws_network_acl.ingress to be a set of objects")
	}
	attrTypes := ingElemsType.AttributeTypes()
	for key, val := range attrTypes {
		t.Logf("%s=%v", key, val.FriendlyName())
	}
}
