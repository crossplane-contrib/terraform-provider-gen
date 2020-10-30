package provider

import (
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/template"
	"github.com/hashicorp/terraform/providers"
)

type SchemaTranslatorConfiguration struct {
	CRDVersion   string
	BasePath     string
	PackagePath  string
	ProviderName string
}

type SchemaTranslator struct {
	cfg     *SchemaTranslatorConfiguration
	schema  providers.GetSchemaResponse
	renamer StringTransformer
	tg      template.TemplateGetter
}

func (st *SchemaTranslator) WriteAllTypeDefFiles() error {
	for name, s := range st.schema.ResourceTypes {
		namer := NewTerraformResourceNamer(st.cfg.ProviderName, name)
		pt := NewPackageTranslator(s, namer, st.cfg, st.tg)
		err := pt.WriteTypeDefFile()
		if err != nil {
			return err
		}
	}

	return nil
}

func NewSchemaTranslator(cfg *SchemaTranslatorConfiguration, schema providers.GetSchemaResponse, tg template.TemplateGetter) *SchemaTranslator {
	return &SchemaTranslator{
		cfg:    cfg,
		schema: schema,
		tg:     tg,
	}
}
