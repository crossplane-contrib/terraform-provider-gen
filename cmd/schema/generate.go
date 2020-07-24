package schema

import (
	"github.com/crossplane/hiveworld/pkg/generator"
	"github.com/crossplane/provider-terraform-plugin/pkg/client"
)

func GenerateSchema(onlyGenerateResourceFlag *string, provider *client.Provider) error {
	gen, err := generator.NewSchemaGenerator(provider, generator.WithResourceName(onlyGenerateResourceFlag))
	if err != nil {
		return err
	}
	return gen.Generate()
}
