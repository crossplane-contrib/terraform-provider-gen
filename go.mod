module github.com/crossplane/hiveworld

go 1.13

replace github.com/crossplane/terraform-provider-runtime => /Users/kasey/src/crossplane/terraform-provider-runtime

require (
	cloud.google.com/go/logging v1.0.0
	github.com/alecthomas/kong v0.2.9
	github.com/crossplane/crossplane-runtime v0.9.0
	github.com/crossplane/crossplane-tools v0.0.0-20200412230150-efd0edd4565b
	github.com/crossplane/terraform-provider-runtime v0.0.0-00010101000000-000000000000
	github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/hashicorp/go-hclog v0.0.0-20181001195459-61d530d6c27f
	github.com/hashicorp/go-plugin v1.3.0
	github.com/hashicorp/terraform v0.12.29
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/vmihailenco/msgpack v4.0.1+incompatible
	github.com/zclconf/go-cty v1.5.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/controller-tools v0.2.4
	sigs.k8s.io/yaml v1.2.0
)
