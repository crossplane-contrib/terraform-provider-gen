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

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// +kubebuilder:object:root=true

// S3Bucket is a managed resource representing a resource mirrored in the cloud
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
type S3Bucket struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   S3BucketSpec   `json:"spec"`
	Status S3BucketStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// S3Bucket contains a list of S3BucketList
type S3BucketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []S3Bucket `json:"items"`
}

// A S3BucketSpec defines the desired state of a S3Bucket
type S3BucketSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       S3BucketParameters `json:"forProvider"`
}

// A S3BucketParameters defines the desired state of a S3Bucket
type S3BucketParameters struct {
	Bucket                            string                            `json:"bucket"`
	BucketPrefix                      string                            `json:"bucket_prefix"`
	HostedZoneId                      string                            `json:"hosted_zone_id"`
	Policy                            string                            `json:"policy"`
	WebsiteEndpoint                   string                            `json:"website_endpoint"`
	AccelerationStatus                string                            `json:"acceleration_status"`
	ForceDestroy                      bool                              `json:"force_destroy"`
	WebsiteDomain                     string                            `json:"website_domain"`
	Acl                               string                            `json:"acl"`
	Tags                              map[string]string                 `json:"tags,omitempty"`
	RequestPayer                      string                            `json:"request_payer"`
	Arn                               string                            `json:"arn"`
	ReplicationConfiguration          ReplicationConfiguration          `json:"replication_configuration"`
	CorsRule                          CorsRule                          `json:"cors_rule"`
	Logging                           Logging                           `json:"logging"`
	ObjectLockConfiguration           ObjectLockConfiguration           `json:"object_lock_configuration"`
	Versioning                        Versioning                        `json:"versioning"`
	Website                           Website                           `json:"website"`
	Grant                             Grant                             `json:"grant"`
	LifecycleRule                     LifecycleRule                     `json:"lifecycle_rule"`
	ServerSideEncryptionConfiguration ServerSideEncryptionConfiguration `json:"server_side_encryption_configuration"`
}

type ReplicationConfiguration struct {
	Role  string  `json:"role"`
	Rules []Rules `json:"rules"`
}

type Rules struct {
	Id                      string                  `json:"id"`
	Prefix                  string                  `json:"prefix"`
	Priority                int64                   `json:"priority"`
	Status                  string                  `json:"status"`
	Destination             Destination             `json:"destination"`
	Filter                  Filter                  `json:"filter"`
	SourceSelectionCriteria SourceSelectionCriteria `json:"source_selection_criteria"`
}

type Destination struct {
	AccountId                string                   `json:"account_id"`
	Bucket                   string                   `json:"bucket"`
	ReplicaKmsKeyId          string                   `json:"replica_kms_key_id"`
	StorageClass             string                   `json:"storage_class"`
	AccessControlTranslation AccessControlTranslation `json:"access_control_translation"`
}

type AccessControlTranslation struct {
	Owner string `json:"owner"`
}

type Filter struct {
	Prefix string            `json:"prefix"`
	Tags   map[string]string `json:"tags,omitempty"`
}

type SourceSelectionCriteria struct {
	SseKmsEncryptedObjects SseKmsEncryptedObjects `json:"sse_kms_encrypted_objects"`
}

type SseKmsEncryptedObjects struct {
	Enabled bool `json:"enabled"`
}

type CorsRule struct {
	MaxAgeSeconds  int64    `json:"max_age_seconds"`
	AllowedHeaders []string `json:"allowed_headers,omitempty"`
	AllowedMethods []string `json:"allowed_methods,omitempty"`
	AllowedOrigins []string `json:"allowed_origins,omitempty"`
	ExposeHeaders  []string `json:"expose_headers,omitempty"`
}

type Logging struct {
	TargetBucket string `json:"target_bucket"`
	TargetPrefix string `json:"target_prefix"`
}

type ObjectLockConfiguration struct {
	ObjectLockEnabled string `json:"object_lock_enabled"`
	Rule              Rule   `json:"rule"`
}

type Rule struct {
	DefaultRetention DefaultRetention `json:"default_retention"`
}

type DefaultRetention struct {
	Years int64  `json:"years"`
	Days  int64  `json:"days"`
	Mode  string `json:"mode"`
}

type Versioning struct {
	Enabled   bool `json:"enabled"`
	MfaDelete bool `json:"mfa_delete"`
}

type Website struct {
	ErrorDocument         string `json:"error_document"`
	IndexDocument         string `json:"index_document"`
	RedirectAllRequestsTo string `json:"redirect_all_requests_to"`
	RoutingRules          string `json:"routing_rules"`
}

type Grant struct {
	Id          string   `json:"id"`
	Permissions []string `json:"permissions,omitempty"`
	Type        string   `json:"type"`
	Uri         string   `json:"uri"`
}

type LifecycleRule struct {
	AbortIncompleteMultipartUploadDays int64                       `json:"abort_incomplete_multipart_upload_days"`
	Enabled                            bool                        `json:"enabled"`
	Id                                 string                      `json:"id"`
	Prefix                             string                      `json:"prefix"`
	Tags                               map[string]string           `json:"tags,omitempty"`
	NoncurrentVersionTransition        NoncurrentVersionTransition `json:"noncurrent_version_transition"`
	Transition                         Transition                  `json:"transition"`
	Expiration                         Expiration                  `json:"expiration"`
	NoncurrentVersionExpiration        NoncurrentVersionExpiration `json:"noncurrent_version_expiration"`
}

type NoncurrentVersionTransition struct {
	StorageClass string `json:"storage_class"`
	Days         int64  `json:"days"`
}

type Transition struct {
	Date         string `json:"date"`
	Days         int64  `json:"days"`
	StorageClass string `json:"storage_class"`
}

type Expiration struct {
	Date                      string `json:"date"`
	Days                      int64  `json:"days"`
	ExpiredObjectDeleteMarker bool   `json:"expired_object_delete_marker"`
}

type NoncurrentVersionExpiration struct {
	Days int64 `json:"days"`
}

type ServerSideEncryptionConfiguration struct {
	Rule Rule `json:"rule"`
}

type Rule struct {
	ApplyServerSideEncryptionByDefault ApplyServerSideEncryptionByDefault `json:"apply_server_side_encryption_by_default"`
}

type ApplyServerSideEncryptionByDefault struct {
	KmsMasterKeyId string `json:"kms_master_key_id"`
	SseAlgorithm   string `json:"sse_algorithm"`
}

// A S3BucketStatus defines the observed state of a S3Bucket
type S3BucketStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          S3BucketObservation `json:"atProvider"`
}

// A S3BucketObservation records the observed state of a S3Bucket
type S3BucketObservation struct {
	Region                   string `json:"region"`
	BucketRegionalDomainName string `json:"bucket_regional_domain_name"`
	BucketDomainName         string `json:"bucket_domain_name"`
}
