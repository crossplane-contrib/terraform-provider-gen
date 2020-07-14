package client

import (
	"io/ioutil"

	"github.com/hashicorp/terraform/configs/configschema"
	tfplugin "github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/providers"
	"github.com/zclconf/go-cty/cty"
	"sigs.k8s.io/yaml"
)

// FakeTerraformVersion is a nice lie we tell to providers to keep them happy
// TODO: is there a more sane way to negotiate version compat w/ providers?
const FakeTerraformVersion string = "v0.12.26"

// VeryBadHardcodedPluginDirectory simplifies plugin directory discovery by assuming
// that we'll eventually just have a single canonical location for plugins, maybe w/ a flag to override
// This will get replaced by some reasonable default that will always be used in docker containers
// with the explicit value only being specified for dev builds outside of containers.
const VeryBadHardcodedPluginDirectory string = "/Users/kasey/src/crossplane/hiveworld/.terraform/plugins/darwin_amd64/"

// Provider wraps grpcProvider with some additional metadata like the provider name
type Provider struct {
	GRPCProvider *tfplugin.GRPCProvider
	Config       ProviderConfig
}

// ProviderConfig models the on-disk yaml config for providers
type ProviderConfig struct {
	ProviderName string            `json:"provider_name"`
	Version      string            `json:"version"`
	PluginDir    string            `json:"plugin_dir"`
	Config       map[string]string `json:"config"`
}

// Configure calls the provider's grpc configuration interface,
// also translating the ProviderConfig structure to the
// Provider's encoded HCL representation.
func (p *Provider) Configure() error {
	schema, err := GetProviderSchema(p)
	if err != nil {
		return err
	}
	ctyCfg := PopulateConfig(schema, p.Config)

	cfgReq := providers.ConfigureRequest{
		TerraformVersion: FakeTerraformVersion,
		Config:           ctyCfg,
	}
	cfgResp := p.GRPCProvider.Configure(cfgReq)
	if cfgResp.Diagnostics.HasErrors() {
		return cfgResp.Diagnostics.Err()
	}

	return nil
}

// ReadProviderConfigFile reads a yaml-formatted provider config and unmarshals
// it into a ProviderConfig, which knows how to generate the serialized
// provider config that a terraform provider expects.
func ReadProviderConfigFile(path string) (ProviderConfig, error) {
	cfg := ProviderConfig{}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = yaml.UnmarshalStrict(content, &cfg)
	if err != nil {
		return cfg, err
	}

	if cfg.PluginDir == "" {
		cfg.PluginDir = VeryBadHardcodedPluginDirectory
	}
	return cfg, err
}

// NewProvider constructs a Provider, which is a container type, holding a
// terraform provider plugin grpc client, as well as metadata about this provider
// instance, eg its configuration and type.
func NewProvider(cfg ProviderConfig) (*Provider, error) {
	grpc, err := NewGRPCProvider(cfg.ProviderName, cfg.PluginDir)
	if err != nil {
		return nil, err
	}
	provider := &Provider{
		GRPCProvider: grpc,
		Config:       cfg,
	}
	err = provider.Configure()

	return provider, err
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

func GetProviderSchema(p *Provider) (*configschema.Block, error) {
	resp := p.GRPCProvider.GetSchema()
	if resp.Diagnostics.HasErrors() {
		return resp.Provider.Block, resp.Diagnostics.NonFatalErr()
	}
	return resp.Provider.Block, nil
}

func PopulateConfig(schema *configschema.Block, cfg ProviderConfig) cty.Value {
	merged := make(map[string]cty.Value)
	for key, attr := range schema.Attributes {
		if val, ok := cfg.Config[key]; ok {
			merged[key] = cty.StringVal(val)
		} else {
			switch attr.Type.FriendlyName() {
			case "string":
				merged[key] = cty.NullVal(cty.String)
				continue
			case "bool":
				merged[key] = cty.NullVal(cty.Bool)
				continue
			case "list of string":
				merged[key] = cty.ListValEmpty(cty.String)
			default:
				merged[key] = cty.NullVal(cty.EmptyObject)
			}
		}
	}
	/*
		for _, block := range cfgSchema.BlockTypes {
			for k2, v2 := range block.Attributes {
				fmt.Printf("- sub key %s = %s", k2, v2.Type.FriendlyName())
			}
		}
	*/
	batching := make(map[string]cty.Value)
	batching["enable_batching"] = cty.BoolVal(false)
	batching["send_after"] = cty.StringVal("3s")
	batchList := []cty.Value{cty.ObjectVal(batching)}
	merged["batching"] = cty.ListVal(batchList)

	val := cty.ObjectVal(merged)
	return val
}
