/*
Copyright 2019 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"context"
	"reflect"

	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/terraform-provider-runtime/pkg/client"
	"github.com/crossplane/terraform-provider-runtime/pkg/registry"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	kubeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// Package type metadata.
const (
	Group                         = "gcp.terraform-plugin.crossplane.io"
	Version                       = "v1alpha1"
	errProviderNotRetrieved       = "provider could not be retrieved"
	errProviderSecretNotRetrieved = "secret referred in provider could not be retrieved"
	errProviderSecretNil          = "cannot find Secret reference on Provider"
	ProviderName                  = "google"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: Group, Version: Version}
	// Provider type metadata.
	ProviderKind             = reflect.TypeOf(Provider{}).Name()
	ProviderGroupKind        = schema.GroupKind{Group: Group, Kind: ProviderKind}.String()
	ProviderKindAPIVersion   = ProviderKind + "." + SchemeGroupVersion.String()
	ProviderGroupVersionKind = SchemeGroupVersion.WithKind(ProviderKind)
)

func initializeProvider(ctx context.Context, mr resource.Managed, kube kubeclient.Client) (*client.Provider, error) {
	provider := &Provider{}
	nn := meta.NamespacedNameOf(mr.GetProviderReference())
	if err := kube.Get(ctx, nn, provider); err != nil {
		return nil, errors.Wrap(err, errProviderNotRetrieved)
	}

	if provider.GetCredentialsSecretReference() == nil {
		return nil, errors.New(errProviderSecretNil)
	}

	secret := &v1.Secret{}
	n := types.NamespacedName{Namespace: provider.Spec.CredentialsSecretRef.Namespace, Name: provider.Spec.CredentialsSecretRef.Name}
	if err := kube.Get(ctx, n, secret); err != nil {
		return nil, errors.Wrap(err, errProviderSecretNotRetrieved)
	}
	credentialString := string(secret.Data[provider.Spec.CredentialsSecretRef.Key])
	cfg := populateConfig(provider, credentialString)

	return client.NewProvider(ProviderName, cfg)
}

// Note that this config still needs to have null values filled in with the correct structure
func populateConfig(p *Provider, credentials string) map[string]cty.Value {
	merged := make(map[string]cty.Value)
	merged["project"] = cty.StringVal(p.Spec.Project)
	merged["region"] = cty.StringVal(p.Spec.Region)
	merged["zone"] = cty.StringVal(p.Spec.Zone)
	merged["credentials"] = cty.StringVal(credentials)

	batching := make(map[string]cty.Value)
	batching["enable_batching"] = cty.BoolVal(false)
	batching["send_after"] = cty.StringVal("3s")
	batchList := []cty.Value{cty.ObjectVal(batching)}
	merged["batching"] = cty.ListVal(batchList)

	return merged
}

func GetProviderEntry() *registry.ProviderEntry {
	schemeBuilder := &scheme.Builder{GroupVersion: SchemeGroupVersion}
	schemeBuilder.Register(&Provider{}, &ProviderList{})
	return &registry.ProviderEntry{
		SchemeBuilder: schemeBuilder,
		Initializer:   initializeProvider,
	}
}
