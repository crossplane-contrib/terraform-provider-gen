module github.com/crossplane-contrib/terraform-provider-gen

go 1.13

replace github.com/crossplane-contrib/terraform-runtime => /Users/kasey/src/crossplane/terraform-provider-runtime

require (
	github.com/crossplane-contrib/terraform-runtime v0.0.0-00010101000000-000000000000
	github.com/crossplane/crossplane-runtime v0.10.0
	github.com/dave/jennifer v1.3.0
	github.com/hashicorp/terraform v0.13.5
	github.com/iancoleman/strcase v0.1.2
	github.com/pkg/errors v0.9.1
	github.com/zclconf/go-cty v1.7.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/apimachinery v0.18.8
)

replace (
	k8s.io/api => k8s.io/api v0.18.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.6
	k8s.io/client-go => k8s.io/client-go v0.18.6
)
