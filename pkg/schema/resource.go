package schema

import (
	"github.com/hashicorp/terraform/providers"
	"github.com/zclconf/go-cty/cty"
)

type SchemaAttribute struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Computed  bool   `json:"computed"`
	Optional  bool   `json:"optional"`
	Required  bool   `json:"required"`
	Sensitive bool   `json:"sensitive"`
	encodeFn  func(*SchemaAttribute) (cty.Value, error)
}

type SchemaGrouping struct {
	Parameters  map[string]*SchemaAttribute `json:"parameters"`
	Observation map[string]*SchemaAttribute `json:"observation"`
}

type Resource struct {
	Name           string            `json:"name"`
	ProviderSchema *providers.Schema `json:"schema"`
	Grouping       *SchemaGrouping   `json:"grouping"`
	encodeCtyFn    func(*Resource) (cty.Value, error)
}

func (r *Resource) AsCtyValue() (cty.Value, error) {
	return r.encodeCtyFn(r)
}

func NewResourceFromProviderSchema(name string, schema *providers.Schema) (*Resource, error) {
	g := &SchemaGrouping{
		Parameters:  make(map[string]*SchemaAttribute),
		Observation: make(map[string]*SchemaAttribute),
	}
	for name, attr := range schema.Block.Attributes {
		a := &SchemaAttribute{
			Name:      name,
			Type:      attr.Type.FriendlyName(),
			Computed:  attr.Computed,
			Optional:  attr.Optional,
			Required:  attr.Required,
			Sensitive: attr.Sensitive,
		}
		if a.Computed {
			g.Observation[name] = a
			continue
		}
		g.Parameters[name] = a
	}
	return &Resource{
		Name:           name,
		ProviderSchema: schema,
		Grouping:       g,
	}, nil
}
