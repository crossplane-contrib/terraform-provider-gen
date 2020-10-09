package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/crossplane/terraform-provider-gen/cmd/schema"
	"github.com/crossplane/terraform-provider-runtime/pkg/client"
)

var (
	gen          = kingpin.New("terraform-provider-gen", "A cli for interacting with terraform providers.")
	pluginPath   = gen.Flag("pluginPath", "Path to provider plugin binary.").Required().String()
	providerName = gen.Flag("providerName", "Terraform provider name, ie the value given to the 'provider' directive in a terraform config.").Required().String()

	schemaCmd     = gen.Command("schema", "subcommand for schema operations.")
	dumpSchemaCmd = schemaCmd.Command("dump", "Print schema to stdout.")
	jsonDumpFlag  = dumpSchemaCmd.Flag("json", "Output schema formatted as a json object.").Bool()

	generateSchemaCmd        = schemaCmd.Command("generate", "Use Provider.GetSchema() to generate crossplane types.")
	onlyGenerateResourceFlag = generateSchemaCmd.Flag("resource", "Limit generation to the single resource named by this flag.").String()
)

func main() {
	gen.FatalIfError(run(), "Error while executing hiveworld command")
}

func run() error {
	switch kingpin.MustParse(gen.Parse(os.Args[1:])) {
	case dumpSchemaCmd.FullCommand():
		provider, err := client.NewProvider(*providerName, *pluginPath)
		defer provider.GRPCProvider.Close()
		if err != nil {
			return err
		}
		schema.Dump(provider, *jsonDumpFlag)
		return nil
	case generateSchemaCmd.FullCommand():
		provider, err := client.NewProvider(*providerName, *pluginPath)
		defer provider.GRPCProvider.Close()
		if err != nil {
			return err
		}
		err = schema.GenerateSchema(onlyGenerateResourceFlag, provider)
	}

	return nil
}
