package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/crossplane/hiveworld/cmd/resource"
	"github.com/crossplane/hiveworld/cmd/schema"
	"github.com/crossplane/hiveworld/pkg/client"
)

var (
	hiveworld  = kingpin.New("hiveworld", "A cli for interacting with terraform providers.")
	configPath = hiveworld.Flag("provider-config", "Path to provider configuration file (yaml).").String()

	schemaCmd     = hiveworld.Command("schema", "subcommand for schema operations.")
	dumpSchemaCmd = schemaCmd.Command("dump", "Print schema to stdout.")
	jsonDumpFlag  = dumpSchemaCmd.Flag("json", "Output schema formatted as a json object.").Bool()

	resourceCmd      = hiveworld.Command("resource", "subcommands operating on managed resources.")
	resourceReadCmd  = resourceCmd.Command("read", "Read metadata for managed resource described by on-disk yaml.")
	resourceReadPath = resourceReadCmd.Arg("yaml-path", "Path to resource yaml on disk.").String()
)

func main() {
	hiveworld.FatalIfError(run(), "Error while executing hiveworld command")
}

func newProvider(path string) (*client.Provider, error) {
	cfg, err := client.ReadProviderConfigFile(path)
	if err != nil {
		return nil, err
	}
	return client.NewProvider(cfg)
}

func run() error {
	switch kingpin.MustParse(hiveworld.Parse(os.Args[1:])) {
	case dumpSchemaCmd.FullCommand():
		provider, err := newProvider(*configPath)
		defer provider.GRPCProvider.Close()
		if err != nil {
			return err
		}
		schema.Dump(provider, *jsonDumpFlag)
		return nil
	case resourceReadCmd.FullCommand():
		provider, err := newProvider(*configPath)
		defer provider.GRPCProvider.Close()
		if err != nil {
			return err
		}
		err = resource.ReadResource(*resourceReadPath, provider)
		if err != nil {
			return err
		}
	}

	return nil
}
