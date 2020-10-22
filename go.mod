module github.com/crossplane/terraform-provider-gen

go 1.13

replace github.com/crossplane/terraform-provider-runtime => /Users/kasey/src/crossplane/terraform-provider-runtime

require (
	github.com/crossplane/terraform-provider-runtime v0.0.0-00010101000000-000000000000
	github.com/dave/jennifer v1.3.0
	github.com/hashicorp/terraform v0.12.29
	github.com/pkg/errors v0.9.1
	github.com/zclconf/go-cty v1.5.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)
