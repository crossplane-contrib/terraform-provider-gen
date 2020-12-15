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
	cfg             Config
	schema          providers.GetSchemaResponse
	renamer         StringTransformer
	tg              template.TemplateGetter
	basePath        string
	overlayBasePath string
}

func (st *SchemaTranslator) WriteAllGeneratedResourceFiles() error {
	for name, s := range st.schema.ResourceTypes {
		if st.cfg.IsExcluded(name) {
			fmt.Printf("Skipping resource %s", name)
			continue
		}
		namer := NewTerraformResourceNamer(st.cfg.Name, name, st.cfg.BaseCRDVersion)
		pt := NewPackageTranslator(s, namer, st.basePath, st.overlayBasePath, st.cfg, st.tg)
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
		err = pt.WriteDecodeFile(mr)
		if err != nil {
			return err
		}
		err = pt.WriteCompareFile(mr)
		if err != nil {
			return err
		}

		err = pt.WriteConfigureFile()
		if err != nil {
			return err
		}
		err = pt.WriteDocFile()
		if err != nil {
			return err
		}
		err = pt.WriteIndexFile()
		if err != nil {
			return err
		}
	}

	return nil
}

func NewSchemaTranslator(cfg Config, basePath, overlayBasePath string, schema providers.GetSchemaResponse, tg template.TemplateGetter) *SchemaTranslator {
	return &SchemaTranslator{
		overlayBasePath: overlayBasePath,
		basePath:        basePath,
		cfg:             cfg,
		schema:          schema,
		tg:              tg,
	}
}
