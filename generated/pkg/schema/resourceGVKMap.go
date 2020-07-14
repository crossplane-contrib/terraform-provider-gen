package schema

import (
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
)

var ResourceGVKMap map[string]k8schema.GroupVersionKind = map[string]k8schema.GroupVersionKind{
	"google_service_account": k8schema.FromAPIVersionAndKind("iam.gcp.crossplane.io/v1alpha1", "google_service_account"),
}
