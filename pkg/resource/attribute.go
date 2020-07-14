package resource

import "github.com/zclconf/go-cty/cty"

type ResourceAttribute struct {
	Name      string
	Type      cty.Type
	Value     interface{}
	CtyEncode func() (cty.Value, error)
}
