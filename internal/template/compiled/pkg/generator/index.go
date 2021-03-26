package generator

func Index() string {
	return "/*\nCopyright 2019 The Crossplane Authors.\n\nLicensed under the Apache License, Version 2.0 (the \"License\");\nyou may not use this file except in compliance with the License.\nYou may obtain a copy of the License at\n\n    http://www.apache.org/licenses/LICENSE-2.0\n\nUnless required by applicable law or agreed to in writing, software\ndistributed under the License is distributed on an \"AS IS\" BASIS,\nWITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\nSee the License for the specific language governing permissions and\nlimitations under the License.\n*/\n\npackage v1alpha1\n\nimport (\n\t\"github.com/crossplane-contrib/terraform-runtime/pkg/plugin\"\n\t\"k8s.io/apimachinery/pkg/runtime/schema\"\n\t\"sigs.k8s.io/controller-runtime/pkg/scheme\"\n)\n\n// Package type metadata.\nconst (\n\tGroup   = \"{{ .APIGroup }}\"\n\tVersion = \"{{ .APIVersion }}\"\n)\n\nvar (\n\t// SchemeGroupVersion is group version used to register these objects\n\tSchemeGroupVersion = schema.GroupVersion{Group: Group, Version: Version}\n)\n\nvar (\n\tKind                  = \"{{ .ManagedResourceName }}\"\n\tGroupKind             = schema.GroupKind{Group: Group, Kind: Kind}.String()\n\tKindAPIVersion        = Kind + \".\" + SchemeGroupVersion.String()\n\tGroupVersionKind      = SchemeGroupVersion.WithKind(Kind)\n\tTerraformResourceName = \"{{ .TerraformResourceName }}\"\n)\n\nfunc Implementation() *plugin.Implementation {\n\t// SchemeBuilder is used to add go types to the GroupVersionKind scheme\n\tschemeBuilder := &scheme.Builder{GroupVersion: SchemeGroupVersion}\n\tschemeBuilder.Register(&{{ .ManagedResourceName }}{}, &{{ .ManagedResourceListName }}{})\n\treturn &plugin.Implementation{\n\t\tGVK:                      GroupVersionKind,\n\t\tTerraformResourceName:    TerraformResourceName,\n\t\tSchemeBuilder:            schemeBuilder,\n\t\tReconcilerConfigurer:     &reconcilerConfigurer{},\n\t\tResourceMerger:           &resourceMerger{},\n\t\tCtyEncoder:               &ctyEncoder{},\n\t\tCtyDecoder:               &ctyDecoder{},\n\t}\n}\n"
}
