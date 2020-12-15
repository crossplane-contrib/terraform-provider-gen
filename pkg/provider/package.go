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

func (pt *PackageTranslator) WriteTypeDefFile(mr *generator.ManagedResource) error {
	outputPath := pt.outputPath("types.go")
	fmt.Printf("Writing typedefs for %s to %s\n", pt.namer.ManagedResourceName(), outputPath)
	fh, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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
	outputPath := pt.outputPath("encode.go")
	fmt.Printf("Writing encoder for %s to %s\n", pt.namer.ManagedResourceName(), outputPath)
	fh, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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

func (pt *PackageTranslator) WriteDecodeFile(mr *generator.ManagedResource) error {
	outputPath := pt.outputPath("decode.go")
	fmt.Printf("Writing decoder for %s to %s\n", pt.namer.ManagedResourceName(), outputPath)
	fh, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer fh.Close()
	if err != nil {
		return err
	}
	generated, err := translate.GenerateDecoders(mr, pt.tg)
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString(generated)
	_, err = io.Copy(fh, buf)
	return err
}

func (pt *PackageTranslator) WriteCompareFile(mr *generator.ManagedResource) error {
	outputPath := pt.outputPath("compare.go")
	fmt.Printf("Writing merger for %s to %s\n", pt.namer.ManagedResourceName(), outputPath)
	fh, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer fh.Close()
	if err != nil {
		return err
	}
	generated, err := translate.GenerateMergers(mr, pt.tg)
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString(generated)
	_, err = io.Copy(fh, buf)
	return err
}

func (pt *PackageTranslator) WriteConfigureFile() error {
	return pt.renderWithNamer("configure.go")
}

func (pt *PackageTranslator) WriteDocFile() error {
	return pt.renderWithNamer("doc.go")
}

func (pt *PackageTranslator) WriteIndexFile() error {
	return pt.renderWithNamer("index.go")
}

func (pt *PackageTranslator) renderWithNamer(filename string) error {
	overlaid, err := pt.overlaid(filename)
	if err != nil {
		return err
	}
	if overlaid {
		return nil
	}
	outputPath := pt.outputPath(filename)
	fmt.Printf("Writing %s for %s to %s\n", filename, pt.namer.ManagedResourceName(), outputPath)
	fh, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer fh.Close()
	if err != nil {
		return err
	}
	ttpl, err := pt.tg.Get(pt.templatePath(filename))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = ttpl.Execute(buf, pt.namer)
	if err != nil {
		return err
	}

	_, err = io.Copy(fh, buf)
	return err
}

func (pt *PackageTranslator) overlaid(filename string) (bool, error) {
	overlayPath := pt.overlayPath(filename)
	if _, err := os.Stat(overlayPath); os.IsNotExist(err) {
		return false, nil
	}
	in, err := os.Open(overlayPath)
	if err != nil {
		return false, err
	}
	defer in.Close()
	outputPath := pt.outputPath(filename)
	out, err := os.Create(outputPath)
	if err != nil {
		return false, err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	fmt.Printf("Overlayed %s onto %s\n", overlayPath, outputPath)
	return true, err
}

func (pt *PackageTranslator) templatePath(filename string) string {
	return fmt.Sprintf("hack/template/pkg/generator/%s.tmpl", filename)
}

func (pt *PackageTranslator) overlayPath(filename string) string {
	// we treat overlay files as .txt so they don't confuse the compiler
	ftxt := fmt.Sprintf("%s.txt", filename)
	return path.Join(pt.overlayDir(), ftxt)
}

func (pt *PackageTranslator) overlayDir() string {
	return pt.resourcePath(pt.overlayBasePath)
}
func (pt *PackageTranslator) outputPath(filename string) string {
	return path.Join(pt.outputDir(), filename)
}

func (pt *PackageTranslator) outputDir() string {
	return pt.resourcePath(pt.basePath)
}

func (pt *PackageTranslator) resourcePath(base string) string {
	return path.Join(base, pt.namer.PackageName(), pt.namer.APIVersion())
}

func (pt *PackageTranslator) EnsureOutputLocation() error {
	fmt.Printf("creating basepath=%s\n", pt.basePath)
	err := os.MkdirAll(pt.outputDir(), 0700)
	if err != nil {
		return err
	}
	return nil
}

type PackageTranslator struct {
	namer           TerraformResourceNamer
	resourceSchema  providers.Schema
	cfg             Config
	tg              template.TemplateGetter
	basePath        string
	overlayBasePath string
}

func NewPackageTranslator(s providers.Schema, namer TerraformResourceNamer, basePath, overlayBasePath string, cfg Config, tg template.TemplateGetter) *PackageTranslator {
	return &PackageTranslator{
		namer:           namer,
		resourceSchema:  s,
		cfg:             cfg,
		tg:              tg,
		basePath:        basePath,
		overlayBasePath: overlayBasePath,
	}
}
