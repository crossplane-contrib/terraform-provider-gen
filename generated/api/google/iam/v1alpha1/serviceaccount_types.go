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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
)

// ServiceAccountParameters defines parameters for a desired IAM ServiceAccount
// https://cloud.google.com/iam/docs/reference/rest/v1/projects.serviceAccounts
// The name of the service account (ie the `accountId` parameter of the Create
// call) is determined by the value of the `crossplane.io/external-name`
// annotation. Unless overridden by the user, this annotation is automatically
// populated with the value of the `metadata.name` attribute.
type ServiceAccountParameters struct {
	// DisplayName is an optional user-specified name for the service account.
	// Must be less than or equal to 100 characters.
	// +optional
	DisplayName *string `json:"displayName,omitempty"`

	// Description is an optional user-specified opaque description of the
	// service account. Must be less than or equal to 256 characters.
	// +optional
	Description *string `json:"description,omitempty"`

	// Account ID specifies the part of the user name before the @
	AccountID string `json:"account_id,omitempty"`
}

// ServiceAccountObservation is used to show the observed state of the
// ServiceAccount resource on GCP. All fields in this structure should only
// be populated from GCP responses; any changes made to the k8s resource outside
// of the crossplane gcp controller will be ignored and overwritten.
type ServiceAccountObservation struct {
	// Name is the "relative resource name" of the service account in the following format:
	// projects/{PROJECT_ID}/serviceAccounts/{external-name}.
	// part of https://godoc.org/google.golang.org/genproto/googleapis/iam/admin/v1#ServiceAccount
	// not to be confused with CreateServiceAccountRequest.Name aka ServiceAccountParameters.ProjectName
	Name string `json:"name,omitempty"`

	// ProjectID is the id of the project that owns the service account.
	Project string `json:"project,omitempty"`

	//The unique and stable id of the service account.
	UniqueID string `json:"unique_id,omitempty"`

	// Email is the the email address of the service account.
	// This matches the EMAIL field you would see using `gcloud iam service-accounts list`
	Email string `json:"email,omitempty"`
}

// ServiceAccountSpec defines the desired state of a
// ServiceAccount.
type ServiceAccountSpec struct {
	runtimev1alpha1.ResourceSpec `json:",inline"`
	ForProvider                  ServiceAccountParameters `json:"forProvider"`
}

// ServiceAccountStatus represents the observed state of a
// ServiceAccount.
type ServiceAccountStatus struct {
	runtimev1alpha1.ResourceStatus `json:",inline"`
	AtProvider                     ServiceAccountObservation `json:"atProvider,omitempty"`
}

type ServiceAccountTerraformConfig struct {
	Timeouts ServiceAccountTerraformConfigTimeout `json:"timeouts,omitempty"`
}

// +kubebuilder:object:root=true

// ServiceAccount is a managed resource that represents a Google IAM Service Account.
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="DISPLAYNAME",type="string",JSONPath=".spec.forProvider.displayName"
// +kubebuilder:printcolumn:name="EMAIL",type="string",JSONPath=".status.atProvider.email"
// +kubebuilder:printcolumn:name="DISABLED",type="boolean",JSONPath=".status.atProvider.disabled"
// +kubebuilder:resource:scope=Cluster
type ServiceAccount struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec            ServiceAccountSpec            `json:"spec"`
	Status          ServiceAccountStatus          `json:"status,omitempty"`
	TerraformConfig ServiceAccountTerraformConfig `json:"config,omitempty"`
}

// +kubebuilder:object:root=true

// ServiceAccountList contains a list of ServiceAccount types
type ServiceAccountList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceAccount `json:"items"`
}

type ServiceAccountTerraformConfigTimeout struct {
	Create *string `json:"create,omitempty"`
}

// TODO: discuss the fact that this had to be
// defined as a separate struct in order for code generation
// to not throw a segfault
/*
$ go generate ./...
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x20 pc=0x16bb20c]

goroutine 1 [running]:
sigs.k8s.io/controller-tools/pkg/crd.structToSchema(0xc0051f3680, 0xc000fccde0, 0x0)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/schema.go:322 +0x9c
sigs.k8s.io/controller-tools/pkg/crd.typeToSchema(0xc0051f3680, 0x1a0a040, 0xc000fccde0, 0x1f)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/schema.go:164 +0x374
sigs.k8s.io/controller-tools/pkg/crd.structToSchema(0xc000cc2120, 0xc000fcce20, 0x0)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/schema.go:375 +0x484
sigs.k8s.io/controller-tools/pkg/crd.typeToSchema(0xc000cc2120, 0x1a0a040, 0xc000fcce20, 0x0)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/schema.go:164 +0x374
sigs.k8s.io/controller-tools/pkg/crd.infoToSchema(...)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/schema.go:107
sigs.k8s.io/controller-tools/pkg/crd.(*Parser).NeedSchemaFor(0xc004434550, 0xc00068ff60, 0xc0002c8240, 0x1d)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/parser.go:174 +0x321
sigs.k8s.io/controller-tools/pkg/crd.(*schemaContext).requestSchema(0xc0051f35f0, 0x0, 0x0, 0xc0002c8240, 0x1d)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/schema.go:99 +0x9e
sigs.k8s.io/controller-tools/pkg/crd.localNamedToSchema(0xc0051f35f0, 0xc000fcd1a0, 0x0)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/schema.go:220 +0x1b4
sigs.k8s.io/controller-tools/pkg/crd.typeToSchema(0xc0051f35f0, 0x1a09c00, 0xc000fcd1a0, 0x1f)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/schema.go:154 +0x194
sigs.k8s.io/controller-tools/pkg/crd.structToSchema(0xc000cc2b18, 0xc000fcd1e0, 0x0)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/schema.go:375 +0x484
sigs.k8s.io/controller-tools/pkg/crd.typeToSchema(0xc000cc2b18, 0x1a0a040, 0xc000fcd1e0, 0x0)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/schema.go:164 +0x374
sigs.k8s.io/controller-tools/pkg/crd.infoToSchema(...)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/schema.go:107
sigs.k8s.io/controller-tools/pkg/crd.(*Parser).NeedSchemaFor(0xc004434550, 0xc00068ff60, 0xc0002c61d0, 0xe)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/parser.go:174 +0x321
sigs.k8s.io/controller-tools/pkg/crd.(*Parser).NeedFlattenedSchemaFor(0xc004434550, 0xc00068ff60, 0xc0002c61d0, 0xe)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/parser.go:186 +0xd8
sigs.k8s.io/controller-tools/pkg/crd.(*Parser).NeedCRDFor(0xc004434550, 0xc00043408e, 0x26, 0xc0002c61d0, 0xe, 0x0)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/spec.go:85 +0x60e
sigs.k8s.io/controller-tools/pkg/crd.Generator.Generate(0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc004434500, 0x16de4e1, 0xc0000cae4b)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/crd/gen.go:108 +0x336
sigs.k8s.io/controller-tools/pkg/genall.(*Runtime).Run(0xc0010c4400, 0xc000269540)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/pkg/genall/genall.go:171 +0x15e
main.main.func1(0xc000235b80, 0xc000269540, 0x4, 0x4, 0x0, 0x0)
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/cmd/controller-gen/main.go:176 +0xa6
github.com/spf13/cobra.(*Command).execute(0xc000235b80, 0xc0000a6060, 0x4, 0x4, 0xc000235b80, 0xc0000a6060)
    /Users/kasey/go/pkg/mod/github.com/spf13/cobra@v0.0.5/command.go:826 +0x453
github.com/spf13/cobra.(*Command).ExecuteC(0xc000235b80, 0xc00028e320, 0x4, 0x0)
    /Users/kasey/go/pkg/mod/github.com/spf13/cobra@v0.0.5/command.go:914 +0x2fb
github.com/spf13/cobra.(*Command).Execute(...)
    /Users/kasey/go/pkg/mod/github.com/spf13/cobra@v0.0.5/command.go:864
main.main()
    /Users/kasey/go/pkg/mod/sigs.k8s.io/controller-tools@v0.2.4/cmd/controller-gen/main.go:200 +0x34a
exit status 2
generated/api/google/generate.go:26: running "go": exit status 1
*/
