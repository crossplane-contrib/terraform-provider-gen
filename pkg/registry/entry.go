package registry

import (
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

type Entry struct {
	TerraformResourceName     string
	GVK                       k8schema.GroupVersionKind
	UnmarshalResourceCallback ResourceUnmarshalFunc
	EncodeCtyCallback         CtyEncodeFunc
	DecodeCtyCallback         CtyDecodeFunc
	SchemeBuilder             *scheme.Builder
	YamlEncodeCallback        YAMLEncodeFunc
}
