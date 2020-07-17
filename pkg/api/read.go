package api

import (
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/hiveworld/pkg/client"
	"github.com/crossplane/provider-terraform-plugin/pkg/registry"
	"github.com/hashicorp/terraform/providers"
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// Read returns an up-to-date version of the resource
// TODO: If `id` is unset for a new resource, how do we figure out
// what value needs to be used as the id?
func Read(p *client.Provider, r *registry.Registry, res resource.Managed, gvk k8schema.GroupVersionKind) (resource.Managed, error) {
	// read resource yaml
	// lookup terraform schema from type name
	// TODO: registry of types from indexing provider getschema
	// - we should get schemas from provider, but use codgen to map schema back to go types
	// serialization function should traverse provider schema and then use codegen'd method to look up values for each field
	//
	// traverse provider schema and
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
	req := providers.ReadResourceRequest{
		TypeName:   tfName,
		PriorState: encoded,
		Private:    nil,
	}
	resp := p.GRPCProvider.ReadResource(req)
	if resp.Diagnostics.HasErrors() {
		return res, resp.Diagnostics.NonFatalErr()
	}
	// should we persist resp.Private in a blob in the resource to use on the next call?
	// Risky since size is unbounded, but we might be matching core behavior more carefully
	ctyDecoder, err := r.GetCtyDecoder(gvk)
	if err != nil {
		return nil, err
	}
	return ctyDecoder(res, resp.NewState, s)
}
