package api

import (
	"fmt"

	"github.com/crossplane/hiveworld/pkg/client"
	"github.com/hashicorp/terraform/providers"
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
)

type ProviderSchema struct {
	Name       string
	Attributes map[string]schemaAttribute
}

type schemaAttribute struct {
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Optional    bool   `json:"optional"`
	Computed    bool   `json:"computed"`
	Description string `json:"description"`
}

func GetProviderSchema(p *client.Provider) (ProviderSchema, error) {
	resp := p.GRPCProvider.GetSchema()
	ps := ProviderSchema{
		Name:       p.Config.ProviderName,
		Attributes: make(map[string]schemaAttribute),
	}
	if resp.Diagnostics.HasErrors() {
		return ps, resp.Diagnostics.NonFatalErr()
	}
	cfgSchema := resp.Provider.Block
	for key, attr := range cfgSchema.Attributes {
		ps.Attributes[key] = schemaAttribute{
			Type:        attr.Type.FriendlyName(),
			Required:    attr.Required,
			Optional:    attr.Optional,
			Computed:    attr.Computed,
			Description: attr.Description,
		}
		//fmt.Printf("%s : type=%s, required=%b, optional=%b, computed=%b, description=%s\n", key, attr.Type.FriendlyName(), attr.Required, attr.Optional, attr.Computed, attr.Description)
		// Attribute represents a configuration attribute, within a block.
	}
	return ps, nil
}

func GetSchema(p *client.Provider) (map[string]providers.Schema, error) {
	resp := p.GRPCProvider.GetSchema()
	if resp.Diagnostics.HasErrors() {
		return nil, resp.Diagnostics.NonFatalErr()
	}

	return resp.ResourceTypes, nil
}

func SchemaForGVK(gvk k8schema.GroupVersionKind, p *client.Provider) (*providers.Schema, error) {
	schema, err := GetSchema(p)
	if err != nil {
		return nil, err
	}
	rs, ok := schema[gvk.Kind]
	if !ok {
		return nil, fmt.Errorf("Could not find schema for GVK=%s", gvk.String())
	}
	return &rs, nil
}
