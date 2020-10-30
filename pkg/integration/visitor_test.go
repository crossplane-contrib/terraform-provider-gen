// +build integration

package integration

import (
	"testing"
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
	VisitAllBlocks(unmm.visitor, name, *block)
}
