package provider

import (
	"bytes"
	"io"
	"os"
	"path"

	"github.com/crossplane-contrib/terraform-provider-gen/pkg/template"
	"github.com/hashicorp/terraform/providers"
)

const (
	PROVIDER_MAINGO_PATH = "provider/cmd/provider/main.go.tpl"

	PROVIDER_INDEX_PATH           = "provider/generated/index.go.tpl"
	PROVIDERCONFIG_INIT_PATH      = "provider/generated/index_provider.go.tpl"
	RESOURCE_IMPLEMENTATIONS_PATH = "provider/generated/index_resources.go.tpl"

	PROVIDERCONFIG_DOC_PATH   = "provider/generated/provider/v1alpha1/doc.go.tpl"
	PROVIDERCONFIG_TYPES_PATH = "provider/generated/provider/v1alpha1/types.go.tpl"
	PROVIDERCONFIG_INDEX_PATH = "provider/generated/provider/v1alpha1/index.go.tpl"
)

type Bootstrapper struct {
	cfg    Config
	tg     template.TemplateGetter
	schema providers.GetSchemaResponse
}

func (bs *Bootstrapper) Bootstrap() error {
	if err := bs.WriteMainGo(); err != nil {
		return err
	}
	if err := bs.WriteIndexGo(); err != nil {
		return err
	}
	if err := bs.WriteProviderDoc(); err != nil {
		return err
	}
	if err := bs.WriteProviderTypes(); err != nil {
		return err
	}
	if err := bs.WriteProviderIndex(); err != nil {
		return err
	}
	return nil
}

func (bs *Bootstrapper) WriteMainGo() error {
	path := path.Join(bs.cfg.BasePath, "cmd", "provider", "main.go")
	return bs.writeExecutedConfigTemplate(PROVIDER_MAINGO_PATH, path)
}

func (bs *Bootstrapper) WriteIndexGo() error {
	path := path.Join(bs.cfg.BasePath, "generated", "index.go")
	return bs.writeExecutedConfigTemplate(PROVIDER_INDEX_PATH, path)
}

func (bs *Bootstrapper) WriteProviderInitGo() error {
	path := path.Join(bs.cfg.BasePath, "generated", "index_provider.go")
	return bs.writeExecutedConfigTemplate(PROVIDERCONFIG_INIT_PATH, path)
}

func (bs *Bootstrapper) WriteProviderDoc() error {
	path := path.Join(bs.cfg.BasePath, "generated", "provider", bs.cfg.ProviderConfigVersion, "doc.go")
	return bs.writeExecutedConfigTemplate(PROVIDERCONFIG_DOC_PATH, path)
}

func (bs *Bootstrapper) WriteProviderTypes() error {
	path := path.Join(bs.cfg.BasePath, "generated", "provider", bs.cfg.ProviderConfigVersion, "types.go")
	return bs.writeExecutedConfigTemplate(PROVIDERCONFIG_TYPES_PATH, path)
}

func (bs *Bootstrapper) WriteProviderIndex() error {
	path := path.Join(bs.cfg.BasePath, "generated", "provider", bs.cfg.ProviderConfigVersion, "index.go")
	return bs.writeExecutedConfigTemplate(PROVIDERCONFIG_INDEX_PATH, path)
}

func (bs *Bootstrapper) writeExecutedConfigTemplate(tplPath, outPath string) error {
	dir := path.Dir(outPath)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}
	tpl, err := bs.tg.Get(tplPath)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, bs.cfg)
	if err != nil {
		return err
	}

	fh, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer fh.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(fh, buf)
	return err
}

func NewBootstrapper(cfg Config, tg template.TemplateGetter, schema providers.GetSchemaResponse) *Bootstrapper {
	return &Bootstrapper{
		cfg:    cfg,
		tg:     tg,
		schema: schema,
	}
}
