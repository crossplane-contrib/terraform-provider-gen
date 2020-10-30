package integration

import (
	"fmt"
	"os"

	tfplugin "github.com/hashicorp/terraform/plugin"

	"github.com/crossplane-contrib/terraform-runtime/pkg/client"
)

var (
	envVarProviderName = "CPTF_PLUGIN_PROVIDER_NAME"
	envVarPluginPath   = "CPTF_PLUGIN_PATH"
)

func getEnvOrError(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("Could not retrieve value for environment variable named %s", name)
	}
	return val, nil
}

func getProvider() (*tfplugin.GRPCProvider, error) {
	providerName, err := getEnvOrError(envVarProviderName)
	if err != nil {
		return nil, err
	}
	pluginPath, err := getEnvOrError(envVarPluginPath)
	if err != nil {
		return nil, err
	}
	return client.NewGRPCProvider(providerName, pluginPath)
}
