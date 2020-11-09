package integration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/configs/configschema"

	"github.com/crossplane-contrib/terraform-provider-gen/pkg/translate"
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

type visitor func([]string, []*configschema.NestedBlock)

func VisitAllBlocks(v visitor, name string, block configschema.Block) {
	nb := configschema.NestedBlock{
		Block:   block,
		Nesting: configschema.NestingSingle,
	}
	VisitBlock(v, []string{name}, []*configschema.NestedBlock{&nb})
}

func VisitBlock(v visitor, names []string, blocks []*configschema.NestedBlock) {
	block := blocks[len(blocks)-1]
	v(names, blocks)
	for n, b := range block.BlockTypes {
		VisitBlock(v, append(names, n), append(blocks, b))
	}
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

type UniqueNestingMode struct {
	Mode       string
	MinItems   int
	MaxItems   int
	IsRequired bool
}

type UniqueNestingModeMap map[UniqueNestingMode][]string

func (unmm UniqueNestingModeMap) Visitor(names []string, blocks []*configschema.NestedBlock) {
	b := blocks[len(blocks)-1]
	np := namePath(names)
	unm := UniqueNestingMode{
		Mode:       nestingModeString(b.Nesting),
		MinItems:   b.MinItems,
		MaxItems:   b.MaxItems,
		IsRequired: translate.IsBlockRequired(b),
	}
	l, ok := unmm[unm]
	if !ok {
		l = make([]string, 0)
	}
	unmm[unm] = append(l, np)
}

func namePath(names []string) string {
	return strings.Join(names, ".")
}

func NestingModePrinter(names []string, blocks []*configschema.NestedBlock) {
	b := blocks[len(blocks)-1]
	//fmt.Printf("%s%s: %s (%d, %d)", tabs(len(blocks)), names[len(names)-1], nestingModeString(b.Nesting), b.MinItems, b.MaxItems)
	fmt.Printf("%s: %s (%d, %d)", namePath(names), nestingModeString(b.Nesting), b.MinItems, b.MaxItems)
	if b.Deprecated {
		fmt.Print(" (DEPRECATED)")
	}
	fmt.Print("\n")
}

func MultiVisitor(vs ...visitor) visitor {
	return func(names []string, blocks []*configschema.NestedBlock) {
		for _, v := range vs {
			v(names, blocks)
		}
	}
}

type FlatResourceFinder []string

func (frf *FlatResourceFinder) Visitor(names []string, blocks []*configschema.NestedBlock) {
	if len(blocks) > 1 {
		return
	}
	b := blocks[0]
	if len(b.BlockTypes) == 0 {
		*frf = append(*frf, names[0])
	}
}

var _ visitor = NestingModePrinter
