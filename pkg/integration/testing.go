package integration

import (
	"fmt"
	"os"

	tfplugin "github.com/hashicorp/terraform/plugin"

	"github.com/crossplane-contrib/terraform-provider-gen/pkg/template"
	"github.com/crossplane-contrib/terraform-runtime/pkg/client"
)

var (
	EnvVarProviderName = "CPTF_PLUGIN_PROVIDER_NAME"
	EnvVarPluginPath   = "CPTF_PLUGIN_PATH"
)

func getEnvOrError(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("Could not retrieve value for environment variable named %s", name)
	}
	return val, nil
}

func getProvider(itc *IntegrationTestConfig) (*tfplugin.GRPCProvider, error) {
	providerName, err := itc.ProviderName()
	if err != nil {
		return nil, err
	}
	pluginPath, err := itc.PluginPath()
	if err != nil {
		return nil, err
	}
	return client.NewGRPCProvider(providerName, pluginPath)
}

type IntegrationTestConfig struct {
	providerName string
	pluginPath   string
	repoRoot     string
}

func (itc *IntegrationTestConfig) ProviderName() (string, error) {
	if itc.providerName != "" {
		return itc.providerName, nil
	}
	return getEnvOrError(EnvVarProviderName)
}

func (itc *IntegrationTestConfig) PluginPath() (string, error) {
	if itc.pluginPath != "" {
		return itc.pluginPath, nil
	}
	return getEnvOrError(EnvVarPluginPath)
}

func (itc *IntegrationTestConfig) RepoRoot() (string, error) {
	if itc.repoRoot != "" {
		return itc.repoRoot, nil
	}
	return os.Getwd()
}

func (itc *IntegrationTestConfig) TemplateGetter() (template.TemplateGetter, error) {
	p, err := itc.RepoRoot()
	if err != nil {
		return nil, err
	}
	return template.NewTemplateGetter(p), nil
}

type TestConfigOption func(*IntegrationTestConfig)

func WithRepoRoot(repoRoot string) TestConfigOption {
	return func(itc *IntegrationTestConfig) {
		itc.repoRoot = repoRoot
	}
}

func WithPluginPath(pluginPath string) TestConfigOption {
	return func(itc *IntegrationTestConfig) {
		itc.pluginPath = pluginPath
	}
}

func WithProvidername(providerName string) TestConfigOption {
	return func(itc *IntegrationTestConfig) {
		itc.providerName = providerName
	}
}

func NewIntegrationTestConfig(opts ...TestConfigOption) *IntegrationTestConfig {
	itc := &IntegrationTestConfig{}
	for _, opt := range opts {
		opt(itc)
	}
	return itc
}
