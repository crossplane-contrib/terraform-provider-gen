package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

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
	providerName = gen.Flag("providerName", "Terraform provider name. must match the value given to the 'provider' directive in a terraform config.").String()

	generateSchemaCmd = schemaCmd.Command("generate", "Use Provider.GetSchema() to generate crossplane types.")
	outputDir         = generateSchemaCmd.Flag("output-dir", "output path").String()
	cfgPath           = generateSchemaCmd.Flag("cfg-path", "path to schema generation config yaml").String()
	//packagePath              = generateSchemaCmd.Flag("package-path", "base path for output packages, eg github.com/crossplane-contrib/provider-terraform-aws/generated/resources").Required().String()
	//baseCrdVersion           = generateSchemaCmd.Flag("crd-version", "Base kind version for generated kubernete kinds, eg v1alpha1").Default("v1alpha1").String()
	repoRoot = generateSchemaCmd.Flag("repo-root", "path to the root of the terraform-provider-gen so the binary can find templates (defaults to PWD)").String()

	analyzeCmd       = gen.Command("analyze", "perform analysis on a provider's schemas")
	nestingCmd       = analyzeCmd.Command("nesting", "report on the different nesting paths and modes observed in a provider")
	nestingCmdStyle  = nestingCmd.Flag("report-style", "Choose between summary (organized by nesting type and min/max), or dump (showing all nested values for all resources)").Default("dump").String()
	flatCmd          = analyzeCmd.Command("flat", "Find resources that do not use nesting at all")
	typesIndexCmd    = analyzeCmd.Command("cty-index", "Build an index showing which resources use which cty types")
	excludeTypesList = typesIndexCmd.Flag("exclude-types", "comma separated list of types to ignore (mutually exclusive with include-types)").String()
	includeTypesList = typesIndexCmd.Flag("include-types", "comma separated list of types to include (mutually exclusive with ignore-types)").String()
	listTypes        = typesIndexCmd.Flag("list-types", "Only list the types, leave out the breakdown of where they can be found").Bool()

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
		opts := []integration.TestConfigOption{
			integration.WithPluginPath(*pluginPath),
			integration.WithProvidername(*providerName),
		}
		if repositoryRoot != nil {
			opts = append(opts, integration.WithRepoRoot(*repositoryRoot))
		}
		itc := integration.NewIntegrationTestConfig(opts...)
		err := integration.UpdateAllFixtures(itc)
		if err != nil {
			return err
		}
	case generateSchemaCmd.FullCommand():
		cfg, err := provider.ConfigFromFile(*cfgPath)
		if err != nil {
			return err
		}
		if len(cfg.ExcludeResources) > 0 {
			fmt.Println("Excluding the following resources from codegen:")
			for _, p := range cfg.ExcludeResources {
				fmt.Println(p)
			}
		}

		tg := template.NewTemplateGetter(*repoRoot)
		p, err := client.NewGRPCProvider(cfg.Name, *pluginPath)
		if err != nil {
			return err
		}
		st := provider.NewSchemaTranslator(cfg, *outputDir, p.GetSchema(), tg)
		return st.WriteAllGeneratedResourceFiles()
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
				fmt.Printf("%s: %s (%d, %d, %t)\n", k, b.Mode, b.MinItems, b.MaxItems, b.IsRequired)
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
	case typesIndexCmd.FullCommand():
		skipType, err := skipTypeFunc(*includeTypesList, *excludeTypesList)
		if err != nil {
			return err
		}
		cti := make(integration.CtyTypeIndexer)
		err = doBlockVisit(cti.Visitor)
		if err != nil {
			return err
		}
		for t, l := range cti {
			if skipType(t) {
				continue
			}
			fmt.Printf("%s:\n", t)
			if *listTypes {
				continue
			}
			for _, path := range l {
				fmt.Printf("\t%s\n", path)
			}
		}
	}

	return nil
}

type filterFunc func(t string) bool

func skipTypeFunc(incTypes, exclTypes string) (filterFunc, error) {
	if incTypes != "" && exclTypes != "" {
		return nil, fmt.Errorf("--include-types and --exclude-types flags are mutually exclusive")
	}
	if incTypes == "" && exclTypes == "" {
		return func(t string) bool {
			return false
		}, nil
	}

	if exclTypes != "" {
		ignoreTypes := make(map[string]bool)
		for _, tn := range strings.Split(exclTypes, ",") {
			ignoreTypes[tn] = true
		}
		return func(t string) bool {
			if _, ignored := ignoreTypes[t]; ignored {
				return true
			}
			return false
		}, nil
	}

	includeTypes := make(map[string]bool)
	for _, tn := range strings.Split(incTypes, ",") {
		includeTypes[tn] = true
	}
	return func(t string) bool {
		if _, included := includeTypes[t]; included {
			return false
		}
		return true
	}, nil
}

func doBlockVisit(visitor integration.Visitor) error {
	p, err := client.NewGRPCProvider(*providerName, *pluginPath)
	if err != nil {
		return err
	}
	for name, s := range p.GetSchema().ResourceTypes {
		integration.VisitAllBlocks(visitor, name, *s.Block)
	}
	return nil
}
