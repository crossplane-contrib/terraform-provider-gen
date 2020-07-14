package api

import (
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/hiveworld/pkg/client"
	"github.com/crossplane/hiveworld/pkg/registry"
	"github.com/hashicorp/terraform/providers"
	"github.com/zclconf/go-cty/cty"
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// Delete deletes the given resource from the provider
// In terraform slang this is expressed as asking the provider
// to act on a Nil planned state.
func Delete(p *client.Provider, res resource.Managed, gvk k8schema.GroupVersionKind) error {
	s, err := SchemaForGVK(gvk, p)
	if err != nil {
		return err
	}
	ctyEncoder, err := registry.GetCtyEncoder(gvk)
	if err != nil {
		return err
	}
	encoded, err := ctyEncoder(res, s)
	if err != nil {
		return err
	}
	tfName, err := registry.GetTerraformNameForGVK(gvk)
	if err != nil {
		return err
	}

	// TODO: research how/if the major providers are using Config
	// same goes for the private state blobs that are shuffled around
	req := providers.ApplyResourceChangeRequest{
		TypeName:   tfName,
		PriorState: encoded,
		// TODO: For the purposes of Create, I am assuming that it's fine for
		// Config and PlannedState to be the same
		Config:       cty.NullVal(s.Block.ImpliedType()),
		PlannedState: cty.NullVal(s.Block.ImpliedType()),
	}
	resp := p.GRPCProvider.ApplyResourceChange(req)
	if resp.Diagnostics.HasErrors() {
		return resp.Diagnostics.NonFatalErr()
	}
	return nil
}
