package generator

import (
	"fmt"

	"github.com/hashicorp/terraform/configs/configschema"
)

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

type visitor func(string, []*configschema.NestedBlock)

func VisitAllBlocks(v visitor, name string, block configschema.Block) {
	fmt.Printf("resource=%s\n", name)
	nb := configschema.NestedBlock{
		Block:   block,
		Nesting: configschema.NestingSingle,
	}
	VisitBlock(v, name, []*configschema.NestedBlock{&nb})
}

func VisitBlock(v visitor, name string, blocks []*configschema.NestedBlock) {
	block := blocks[len(blocks)-1]
	for n, b := range block.BlockTypes {
		v(n, append(blocks, b))
		VisitBlock(v, n, append(blocks, b))
	}
}

func NestingModePrinter(name string, blocks []*configschema.NestedBlock) {
	b := blocks[len(blocks)-1]
	fmt.Printf("%s%s: %s (%d, %d)", tabs(len(blocks)), name, nestingModeString(b.Nesting), b.MinItems, b.MaxItems)
	if b.Deprecated {
		fmt.Print(" (DEPRECATED)")
	}
	fmt.Print("\n")
}

func tabs(n int) string {
	t := ""
	for i := 0; i < n; i++ {
		t += "\t"
	}
	return t
}

func nestingModeString(nm configschema.NestingMode) string {
	switch nm {
	// NestingSingle indicates that only a single instance of a given
	// block type is permitted, with no labels, and its content should be
	// provided directly as an object value.
	case configschema.NestingSingle:
		return "NestingSingle"

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
		return "NestingGroup"

	// NestingList indicates that multiple blocks of the given type are
	// permitted, with no labels, and that their corresponding objects should
	// be provided in a list.
	case configschema.NestingList:
		return "NestingList"

	// NestingSet indicates that multiple blocks of the given type are
	// permitted, with no labels, and that their corresponding objects should
	// be provided in a set.
	case configschema.NestingSet:
		return "NestingSet"

	// NestingMap indicates that multiple blocks of the given type are
	// permitted, each with a single label, and that their corresponding
	// objects should be provided in a map whose keys are the labels.
	//
	// It's an error, therefore, to use the same label value on multiple
	// blocks.
	case configschema.NestingMap:
		return "NestingMap"
	default:
		return "invalid"
	}
}
