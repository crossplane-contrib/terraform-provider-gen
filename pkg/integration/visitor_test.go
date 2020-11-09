// +build integration

package integration

import (
	"testing"

	"github.com/hashicorp/terraform/configs/configschema"
)

func TestNestingModeAggregator(t *testing.T) {
	p, err := getProvider()
	if err != nil {
		t.Fatal(err)
	}
	rt := p.GetSchema().ResourceTypes
	name := "aws_kinesis_analytics_application"
	block := rt[name].Block
	unmm := make(UniqueNestingModeMap)
	VisitAllBlocks(unmm.Visitor, name, *block)
	if len(unmm) != 4 {
		t.Errorf("Expected a single entry in UniqueNestingModeMap, found =%d", len(unmm))
	}
	for unq, name := range unmm {
		t.Logf("%v : %v", unq, name)
	}
}

func TestProbe(t *testing.T) {
	name := "aws_appmesh_route"
	p, err := getProvider()
	if err != nil {
		t.Fatal(err)
	}
	rt := p.GetSchema().ResourceTypes
	block := rt[name].Block
	if len(block.Attributes) == 0 {
		t.Log("huh")
	}
}

type probe struct {
	resourceType string
	foundIt      bool
}

func (p *probe) Visitor(names []string, blocks []*configschema.NestedBlock) {
	b := blocks[0]
	if b.Nesting == configschema.NestingSingle {
		if names[0] == p.resourceType {
			p.foundIt = true
		}
	}
}
