package api

import (
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/hiveworld/pkg/client"
	"github.com/crossplane/provider-terraform-plugin/pkg/registry"
	"github.com/hashicorp/terraform/providers"
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// Update syncs with an existing resource and modifies mutable values
func Update(p *client.Provider, r *registry.Registry, res resource.Managed, gvk k8schema.GroupVersionKind) (resource.Managed, error) {
	s, err := SchemaForGVK(gvk, p)
	if err != nil {
		return nil, err
	}
	ctyEncoder, err := r.GetCtyEncoder(gvk)
	if err != nil {
		return nil, err
	}
	encoded, err := ctyEncoder(res, s)
	if err != nil {
		return nil, err
	}
	tfName, err := r.GetTerraformNameForGVK(gvk)
	if err != nil {
		return nil, err
	}

	prior, err := Read(p, r, res, gvk)
	if err != nil {
		return nil, err
	}
	priorEncoded, err := ctyEncoder(prior, s)
	if err != nil {
		return nil, err
	}

	// TODO: research how/if the major providers are using Config
	// same goes for the private state blobs that are shuffled around
	req := providers.ApplyResourceChangeRequest{
		TypeName:   tfName,
		PriorState: priorEncoded,
		// TODO: For the purposes of Create, I am assuming that it's fine for
		// Config and PlannedState to be the same
		Config:       encoded,
		PlannedState: encoded,
	}
	resp := p.GRPCProvider.ApplyResourceChange(req)
	if resp.Diagnostics.HasErrors() {
		return res, resp.Diagnostics.NonFatalErr()
	}
	ctyDecoder, err := r.GetCtyDecoder(gvk)
	if err != nil {
		return nil, err
	}
	return ctyDecoder(res, resp.NewState, s)
}
