package provider

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/crossplane-contrib/terraform-provider-gen/pkg/generator"
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/template"
	"github.com/crossplane-contrib/terraform-provider-gen/pkg/translate"
	"github.com/hashicorp/terraform/providers"
)

type PackageTranslator struct {
	namer          TerraformResourceNamer
	resourceSchema providers.Schema
	cfg            *SchemaTranslatorConfiguration
	tg             template.TemplateGetter
}

func (pt *PackageTranslator) WriteTypeDefFile() error {
	fmt.Printf("basepath=%s\n", pt.cfg.BasePath)
	fmt.Printf("Writing %s to %s\n", pt.namer.ManagedResourceName(), pt.typesPath())
	err := os.MkdirAll(pt.outputDir(), 0700)
	if err != nil {
		return err
	}
	fh, err := os.OpenFile(pt.typesPath(), os.O_RDWR|os.O_CREATE, 0755)
	defer fh.Close()
	if err != nil {
		return err
	}
	mr := translate.SchemaToManagedResource(pt.namer.ManagedResourceName(), pt.cfg.PackagePath, pt.resourceSchema)
	renderer := generator.NewManagedResourceTypeDefRenderer(mr, pt.tg)
	rendered, err := renderer.Render()
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString(rendered)
	_, err = io.Copy(fh, buf)
	return err
}

func (pt *PackageTranslator) typesPath() string {
	return path.Join(pt.outputDir(), "types.go")
}

func (pt *PackageTranslator) outputDir() string {
	return path.Join(pt.cfg.BasePath, pt.namer.PackageName(), pt.cfg.CRDVersion)
}

func NewPackageTranslator(s providers.Schema, namer TerraformResourceNamer, cfg *SchemaTranslatorConfiguration, tg template.TemplateGetter) *PackageTranslator {
	return &PackageTranslator{
		namer:          namer,
		resourceSchema: s,
		cfg:            cfg,
		tg:             tg,
	}
}
