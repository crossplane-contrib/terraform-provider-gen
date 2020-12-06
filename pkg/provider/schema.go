package provider

import (
	"fmt"

	"github.com/crossplane-contrib/terraform-provider-gen/pkg/template"
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/translate"
	"github.com/hashicorp/terraform/providers"
)

type SchemaTranslatorConfiguration struct {
	CRDVersion   string
	PackagePath  string
	ProviderName string
}

type SchemaTranslator struct {
	cfg      Config
	schema   providers.GetSchemaResponse
	renamer  StringTransformer
	tg       template.TemplateGetter
	basePath string
}

func (st *SchemaTranslator) WriteAllGeneratedResourceFiles() error {
	for name, s := range st.schema.ResourceTypes {
		if st.cfg.IsExcluded(name) {
			fmt.Printf("Skipping resource %s", name)
			continue
		}
		namer := NewTerraformResourceNamer(st.cfg.Name, name)
		pt := NewPackageTranslator(s, namer, st.basePath, st.cfg, st.tg)
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

		err = pt.WriteDocFile(st.cfg.BaseCRDVersion, pt.namer.ManagedResourceName(), st.cfg.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewSchemaTranslator(cfg Config, basePath string, schema providers.GetSchemaResponse, tg template.TemplateGetter) *SchemaTranslator {
	return &SchemaTranslator{
		basePath: basePath,
		cfg:      cfg,
		schema:   schema,
		tg:       tg,
	}
}
