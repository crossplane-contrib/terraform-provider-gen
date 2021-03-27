package provider

import (
	"bytes"
	"fmt"
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/optimize"
	"io"
	"os"
	"path"

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

func (st *SchemaTranslator) WriteGeneratedTypes() error {
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
		mr, err = optimize.Deduplicate(mr)
		if err != nil {
			return err
		}
		err = pt.WriteTypeDefFile(mr)
		if err != nil {
			return err
		}

		err = pt.WriteDocFile()
		if err != nil {
			return err
		}
	}
	return nil
}

func (st *SchemaTranslator) WriteGeneratedRuntime() error {
	pis := make([]PackageImport, 0)
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
		mr, err = optimize.Deduplicate(mr)
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
		err = pt.WriteIndexFile()
		if err != nil {
			return err
		}
		pis = append(pis, pt.PackageImport())
	}
	return st.writeResourceImplementationIndex(pis)
}

func (st *SchemaTranslator) writeResourceImplementationIndex(pis []PackageImport) error {
	dir := path.Dir(st.basePath)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}
	tpl, err := st.tg.Get(RESOURCE_IMPLEMENTATIONS_PATH)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	values := struct {
		Config
		PackageImports []PackageImport
	}{
		Config:         st.cfg,
		PackageImports: pis,
	}
	err = tpl.Execute(buf, values)
	if err != nil {
		return err
	}

	outPath := path.Join(dir, "index_resources.go")
	fh, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer fh.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(fh, buf)
	return err
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
