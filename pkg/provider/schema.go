package provider

import (
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/template"
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/translate"
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

func (st *SchemaTranslator) WriteAllGeneratedResourceFiles() error {
	for name, s := range st.schema.ResourceTypes {
		namer := NewTerraformResourceNamer(st.cfg.ProviderName, name)
		pt := NewPackageTranslator(s, namer, st.cfg, st.tg)
		err := pt.EnsureOutputLocation()
		if err != nil {
			return err
		}
		mr := translate.SchemaToManagedResource(pt.namer.ManagedResourceName(), pt.cfg.PackagePath, pt.resourceSchema)
		err = pt.WriteTypeDefFile(mr)
		if err != nil {
			return err
		}
		err = pt.WriteEncoderFile(mr)
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
