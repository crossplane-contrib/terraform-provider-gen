package generator

import (
	"fmt"

	"github.com/crossplane/terraform-provider-runtime/pkg/client"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/providers"
)

type GenerateSchemaOption func(*SchemaGenerator) error

type GenerateSchemaConfig struct {
	resourceName *string
}

type SchemaGenerator struct {
	provider *client.Provider
	schema   map[string]providers.Schema
	cfg      GenerateSchemaConfig
}

func WithResourceName(name *string) GenerateSchemaOption {
	return func(gen *SchemaGenerator) error {
		gen.cfg.resourceName = name
		return nil
	}
}

func (gen *SchemaGenerator) Generate() error {
	resp := gen.provider.GRPCProvider.GetSchema()
	if resp.Diagnostics.HasErrors() {
		return resp.Diagnostics.NonFatalErr()
	}
	schema := resp.ResourceTypes

	for resourceName, r := range schema {
		if gen.cfg.resourceName != nil {
			if *gen.cfg.resourceName != resourceName {
				continue
			}
		}
		fmt.Printf("name=%s", resourceName)
		UnrollBlocks(r.Block, "")
	}

	return nil
}

func UnrollBlocks(block *configschema.Block, indent string) {
	fmt.Println("Attributes")
	for key, attr := range block.Attributes {
		fmt.Printf("%s%s(type=%s, computed=%t, optional=%t, required=%t, sensitive=%t\n", indent, key, attr.Type.FriendlyName(), attr.Computed, attr.Optional, attr.Required, attr.Sensitive)
	}
	for key, b := range block.BlockTypes {
		var mode string
		switch b.Nesting {
		// NestingSingle indicates that only a single instance of a given
		// block type is permitted, with no labels, and its content should be
		// provided directly as an object value.
		case configschema.NestingSingle:
			mode = "NestingSingle"

		// NestingGroup is similar to NestingSingle in that it calls for only a
		// single instance of a given block type with no labels, but it additonally
		// guarantees that its result will never be null, even if the block is
		// absent, and instead the nested attributes and blocks will be treated
		// as absent in that case. (Any required attributes or blocks within the
		// nested block are not enforced unless the block is explicitly present
		// in the configuration, so they are all effectively optional when the
		// block is not present.)
		//
		// This is useful for the situation where a remote API has a feature that
		// is always enabled but has a group of settings related to that feature
		// that themselves have default values. By using NestingGroup instead of
		// NestingSingle in that case, generated plans will show the block as
		// present even when not present in configuration, thus allowing any
		// default values within to be displayed to the user.
		case configschema.NestingGroup:
			mode = "NestingGroup"

		// NestingList indicates that multiple blocks of the given type are
		// permitted, with no labels, and that their corresponding objects should
		// be provided in a list.
		case configschema.NestingList:
			mode = "NestingList"

		// NestingSet indicates that multiple blocks of the given type are
		// permitted, with no labels, and that their corresponding objects should
		// be provided in a set.
		case configschema.NestingSet:
			mode = "NestingSet"

		// NestingMap indicates that multiple blocks of the given type are
		// permitted, each with a single label, and that their corresponding
		// objects should be provided in a map whose keys are the labels.
		//
		// It's an error, therefore, to use the same label value on multiple
		// blocks.
		case configschema.NestingMap:
			mode = "NestingMap"
		default:
			mode = "invalid"
		}
		fmt.Printf("%sblock key=%s; mode=%s\n", indent, key, mode)
		UnrollBlocks(&b.Block, indent+"\t")
	}
}

func NewSchemaGenerator(provider *client.Provider, opts ...GenerateSchemaOption) (*SchemaGenerator, error) {
	gen := &SchemaGenerator{
		provider: provider,
	}
	for _, opt := range opts {
		err := opt(gen)
		if err != nil {
			return nil, err
		}
	}
	return gen, nil
}
