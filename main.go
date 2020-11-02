package main

import (
	"fmt"
	"os"
	"sort"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/crossplane-contrib/terraform-provider-gen/pkg/integration"
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/provider"
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/template"
	"github.com/crossplane-contrib/terraform-runtime/pkg/client"
)

var (
	gen = kingpin.New("terraform-provider-gen", "A cli for interacting with terraform providers.")

	updateFixturesCmd = gen.Command("update-fixtures", "update test fixtures based on current codegen output")
	repositoryRoot    = updateFixturesCmd.Flag("repo-root", "Path to root of repository so that the fixture generator can find paths").Required().String()

	schemaCmd    = gen.Command("schema", "subcommand for schema operations.")
	pluginPath   = gen.Flag("plugin-path", "Path to provider plugin binary.").Required().String()
	providerName = gen.Flag("providerName", "Terraform provider name. must match the value given to the 'provider' directive in a terraform config.").Required().String()

	generateSchemaCmd        = schemaCmd.Command("generate", "Use Provider.GetSchema() to generate crossplane types.")
	onlyGenerateResourceFlag = generateSchemaCmd.Flag("resource", "Limit generation to the single resource named by this flag.").String()
	outputDir                = generateSchemaCmd.Flag("output-dir", "output path").String()
	packagePath              = generateSchemaCmd.Flag("package-path", "base path for output packages, eg github.com/crossplane-contrib/provider-terraform-aws/generated/resources").Required().String()
	baseCrdVersion           = generateSchemaCmd.Flag("crd-version", "Base kind version for generated kubernete kinds, eg v1alpha1").Default("v1alpha1").String()
	repoRoot                 = generateSchemaCmd.Flag("repo-root", "path to the root of the terraform-provider-gen so the binary can find templates (defaults to PWD)").String()

	analyzeCmd      = gen.Command("analyze", "perform analysis on a provider's schemas")
	nestingCmd      = analyzeCmd.Command("nesting", "report on the different nesting paths and modes observed in a provider")
	nestingCmdStyle = nestingCmd.Flag("report-style", "Choose between summary (organized by nesting type and min/max), or dump (showing all nested values for all resources)").Default("dump").String()
	flatCmd         = analyzeCmd.Command("flat", "Find resources that do not use nesting at all")

	//renderCmd = schemaCmd.Command("render", "render crossplane types for the given provider.")
	//dumpSchemaCmd = schemaCmd.Command("dump", "Print schema to stdout.")
	//jsonDumpFlag  = dumpSchemaCmd.Flag("json", "Output schema formatted as a json object.").Bool()
)

func main() {
	gen.FatalIfError(run(), "Error while executing hiveworld command")
}

func run() error {
	switch kingpin.MustParse(gen.Parse(os.Args[1:])) {
	/*
		case dumpSchemaCmd.FullCommand():
			provider, err := client.NewProvider(*providerName, *pluginPath)
			defer provider.GRPCProvider.Close()
			if err != nil {
				return err
			}
			schema.Dump(provider, *jsonDumpFlag)
			return nil
	*/
	case updateFixturesCmd.FullCommand():
		err := integration.UpdateAllFixtures(*repositoryRoot)
		if err != nil {
			return err
		}
	case generateSchemaCmd.FullCommand():
		p, err := client.NewGRPCProvider(*providerName, *pluginPath)
		if err != nil {
			return err
		}
		cfg := &provider.SchemaTranslatorConfiguration{
			CRDVersion:   *baseCrdVersion,
			BasePath:     *outputDir,
			PackagePath:  *packagePath,
			ProviderName: *providerName,
		}
		tg := template.NewTemplateGetter(*repoRoot)
		st := provider.NewSchemaTranslator(cfg, p.GetSchema(), tg)
		return st.WriteAllTypeDefFiles()
	case nestingCmd.FullCommand():
		p, err := client.NewGRPCProvider(*providerName, *pluginPath)
		if err != nil {
			return err
		}
		unmm := make(integration.UniqueNestingModeMap)
		for name, s := range p.GetSchema().ResourceTypes {
			integration.VisitAllBlocks(unmm.Visitor, name, *s.Block)
		}
		switch *nestingCmdStyle {
		case "dump":
			inverted := make(map[string]integration.UniqueNestingMode)
			keys := make([]string, 0)
			for mode, paths := range unmm {
				for _, p := range paths {
					inverted[p] = mode
					keys = append(keys, p)
				}
			}
			sort.Strings(keys)
			for _, k := range keys {
				b := inverted[k]
				fmt.Printf("%s: %s (%d, %d)\n", k, b.Mode, b.MinItems, b.MaxItems)
			}
		default:
			return fmt.Errorf("report-style=%s not recognized", nestingCmdStyle)
		}
	case flatCmd.FullCommand():
		p, err := client.NewGRPCProvider(*providerName, *pluginPath)
		if err != nil {
			return err
		}
		frf := make(integration.FlatResourceFinder, 0)
		for name, s := range p.GetSchema().ResourceTypes {
			integration.VisitAllBlocks(frf.Visitor, name, *s.Block)
		}
		sort.Strings(frf)
		for _, r := range frf {
			fmt.Println(r)
		}
	}

	return nil
}
