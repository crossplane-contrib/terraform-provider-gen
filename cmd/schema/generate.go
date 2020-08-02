package schema

import (
	"github.com/crossplane/hiveworld/pkg/generator"
	"github.com/crossplane/terraform-provider-runtime/pkg/client"
)

func GenerateSchema(onlyGenerateResourceFlag *string, provider *client.Provider) error {
	gen, err := generator.NewSchemaGenerator(provider, generator.WithResourceName(onlyGenerateResourceFlag))
	if err != nil {
		return err
	}
	return gen.Generate()
}
