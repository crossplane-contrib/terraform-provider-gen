package google

import (
	iam "github.com/crossplane/hiveworld/generated/api/google/iam/v1alpha1"
	"github.com/crossplane/hiveworld/pkg/registry"
)

// TODO: out this list somewhere in the codegen pipeline
var RegistryEntries = []*registry.Entry{
	iam.ServiceAccountRegistryEntry(),
}
