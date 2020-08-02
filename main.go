package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/crossplane/terraform-provider-gen/cmd/resource"
	"github.com/crossplane/terraform-provider-gen/cmd/schema"
	"github.com/crossplane/terraform-provider-gen/generated/api/google"
	"github.com/crossplane/terraform-provider-runtime/pkg/client"
	"github.com/crossplane/terraform-provider-runtime/pkg/registry"
)

var (
	hiveworld  = kingpin.New("hiveworld", "A cli for interacting with terraform providers.")
	configPath = hiveworld.Flag("provider-config", "Path to provider configuration file (yaml).").String()

	schemaCmd     = hiveworld.Command("schema", "subcommand for schema operations.")
	dumpSchemaCmd = schemaCmd.Command("dump", "Print schema to stdout.")
	jsonDumpFlag  = dumpSchemaCmd.Flag("json", "Output schema formatted as a json object.").Bool()

	generateSchemaCmd        = schemaCmd.Command("generate", "Use Provider.GetSchema() to generate crossplane types.")
	onlyGenerateResourceFlag = generateSchemaCmd.Flag("resource", "Limit generation to the single resource named by this flag.").String()

	resourceCmd       = hiveworld.Command("resource", "subcommands operating on managed resources.")
	resourcePath      = resourceCmd.Flag("yaml-path", "Path to resource yaml on disk.").String()
	resourceReadCmd   = resourceCmd.Command("read", "Read metadata for managed resource described by on-disk yaml.")
	resourceCreateCmd = resourceCmd.Command("create", "Create managed resource described by on-disk yaml.")
	resourceUpdateCmd = resourceCmd.Command("update", "Update managed resource to state described by on-disk yaml.")
	resourceDeleteCmd = resourceCmd.Command("delete", "Update managed resource identified by on-disk yaml resource.")
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
	r := registry.NewRegistry()
	google.Register(r)
	switch kingpin.MustParse(hiveworld.Parse(os.Args[1:])) {
	case dumpSchemaCmd.FullCommand():
		provider, err := newProvider(*configPath)
		defer provider.GRPCProvider.Close()
		if err != nil {
			return err
		}
		schema.Dump(provider, *jsonDumpFlag)
		return nil
	case generateSchemaCmd.FullCommand():
		provider, err := newProvider(*configPath)
		defer provider.GRPCProvider.Close()
		if err != nil {
			return err
		}
		err = schema.GenerateSchema(onlyGenerateResourceFlag, provider)
	case resourceReadCmd.FullCommand():
		provider, err := newProvider(*configPath)
		defer provider.GRPCProvider.Close()
		if err != nil {
			return err
		}
		err = resource.ReadResource(*resourcePath, provider, r)
		if err != nil {
			return err
		}
	case resourceCreateCmd.FullCommand():
		provider, err := newProvider(*configPath)
		defer provider.GRPCProvider.Close()
		if err != nil {
			return err
		}
		err = resource.CreateResource(*resourcePath, provider, r)
		if err != nil {
			return err
		}
	case resourceUpdateCmd.FullCommand():
		provider, err := newProvider(*configPath)
		defer provider.GRPCProvider.Close()
		if err != nil {
			return err
		}
		err = resource.UpdateResource(*resourcePath, provider, r)
		if err != nil {
			return err
		}
	case resourceDeleteCmd.FullCommand():
		provider, err := newProvider(*configPath)
		defer provider.GRPCProvider.Close()
		if err != nil {
			return err
		}
		err = resource.DeleteResource(*resourcePath, provider, r)
		if err != nil {
			return err
		}
	}

	return nil
}
