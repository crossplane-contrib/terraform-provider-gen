package schema

import (
	"github.com/crossplane/hiveworld/pkg/client"
	"github.com/crossplane/hiveworld/pkg/generator"
)

func GenerateSchema(onlyGenerateResourceFlag *string, provider *client.Provider) error {
	gen, err := generator.NewSchemaGenerator(provider, generator.WithResourceName(onlyGenerateResourceFlag))
	if err != nil {
		return err
	}
	return gen.Generate()
}
