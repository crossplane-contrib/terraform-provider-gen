package schema

import (
	"encoding/json"
	"fmt"

	"github.com/crossplane-contrib/terraform-runtime/pkg/client"
)

// Dump prints out the schema returned by the provider named by the provider arg
// if the json flag is true, formats the output to json
func Dump(provider *client.Provider, jsonOut bool) error {
	schema, err := getProviderSchema(provider)
	if err != nil {
		return err
	}
	if jsonOut {
		jsonb, err := json.Marshal(schema)
		if err != nil {
			return err
		}
		fmt.Println(string(jsonb))
		return nil
	}
	return nil
}

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

func getProviderSchema(p *client.Provider) (ProviderSchema, error) {
	resp := p.GRPCProvider.GetSchema()
	ps := ProviderSchema{
		Name:       p.Name,
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
