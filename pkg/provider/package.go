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

func (pt *PackageTranslator) EnsureOutputLocation() error {
	fmt.Printf("creating basepath=%s\n", pt.cfg.BasePath)
	err := os.MkdirAll(pt.outputDir(), 0700)
	if err != nil {
		return err
	}
	return nil
}

func (pt *PackageTranslator) WriteTypeDefFile(mr *generator.ManagedResource) error {
	fmt.Printf("Writing typedefs for %s to %s\n", pt.namer.ManagedResourceName(), pt.typesPath())
	fh, err := os.OpenFile(pt.typesPath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer fh.Close()
	if err != nil {
		return err
	}
	renderer := generator.NewManagedResourceTypeDefRenderer(mr, pt.tg)
	rendered, err := renderer.Render()
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString(rendered)
	_, err = io.Copy(fh, buf)
	return err
}

func (pt *PackageTranslator) WriteEncoderFile(mr *generator.ManagedResource) error {
	fmt.Printf("Writing encoder for %s to %s\n", pt.namer.ManagedResourceName(), pt.encoderPath())
	fh, err := os.OpenFile(pt.encoderPath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer fh.Close()
	if err != nil {
		return err
	}
	generated, err := translate.GenerateEncoders(mr, pt.tg)
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString(generated)
	_, err = io.Copy(fh, buf)
	return err
}

func (pt *PackageTranslator) typesPath() string {
	return path.Join(pt.outputDir(), "types.go")
}

func (pt *PackageTranslator) encoderPath() string {
	return path.Join(pt.outputDir(), "encode.go")
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
