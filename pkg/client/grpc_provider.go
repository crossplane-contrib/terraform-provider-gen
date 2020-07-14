package client

import (
	"fmt"

	plugin "github.com/hashicorp/go-plugin"
	tfplugin "github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/plugin/discovery"
)

const ProviderPluginType = "provider"

// NewGRPCProvider creates a new GRPCClient instance.
func NewGRPCProvider(providerName, pluginDir string) (*tfplugin.GRPCProvider, error) {
	// 1. find plugins in the filesystem
	// name and version are just parsed out of the provider file name ({name}_v{version})
	pluginMetaSet := discovery.FindPlugins(ProviderPluginType, []string{pluginDir}).WithName(providerName)
	if pluginMetaSet.Count() < 1 {
		return nil, fmt.Errorf("Failed to find plugin: %s. Plugin binary was not found in the plugin directory(%s)", providerName, pluginDir)
	}
	// this is just comparing semvers to find the highest one
	pluginMeta := pluginMetaSet.Newest()

	// plugin.NewClient returns a client that knows how to spawn a provider subprocess and set up the grpc connection
	cfg := tfplugin.ClientConfig(pluginMeta)
	// this discards the noisy debug logs that we get back from go-plugin
	// note that if we want to add options to display these logs w/ verbosity config, setting Output: nil
	// will use the default stdout writer and hclog Levels line up with the standard logging lib.
	/*
		cfg.Logger = hclog.New(&hclog.LoggerOptions{
			Output: ioutil.Discard,
			Level:  hclog.Trace,
			Name:   "plugin",
		})
	*/
	pluginClient := plugin.NewClient(cfg)
	// 2. Spawn the chosen plugin binary as a subprocess, connect to its stdout and parse the grpc connection configuration
	// the connection to the client is set up at this point in the negotiation process, but the protobuf client
	// code hasn't been initialized
	client, err := pluginClient.Client()
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize GRPC plugin: %s", err)
	}

	// 3. Dispense returns a client type that implements the terraform provider grpc interface.
	// It is type `interface{}`` for 2 reasons. First because go-plugin has backwards support for net/rpc
	// in addition to grpc. Second, because up to this point, go-plugin just has a generic grpc *connection*.
	// Dispense is where the connection is actually handed off to the grpc generated code.
	// GRPCProvider is actually one additional level of indirection on top of the GRPC generated code.
	raw, err := client.Dispense(tfplugin.ProviderPluginName)
	if err != nil {
		return nil, fmt.Errorf("Failed to dispense GRPC plugin: %s", err)
	}

	// 4. Finally we type assert that we received a GRPCProvider so that we can rely on
	// TODO: add to our list of issues that handling all the different flavors of provider negotiation adds complexity to our system
	// can we draw a line and only support post .12 plugins which support the modern plugin protocol?
	provider, ok := raw.(*tfplugin.GRPCProvider)
	if !ok {
		return nil, fmt.Errorf("Did not get a plugin provider of type GRPCProvider from client.Dispense")
	}
	provider.PluginClient = pluginClient

	return provider, nil
}
