package generator

func Types() string {
	return "/*\n\tCopyright 2019 The Crossplane Authors.\n\n\tLicensed under the Apache License, Version 2.0 (the \"License\");\n\tyou may not use this file except in compliance with the License.\n\tYou may obtain a copy of the License at\n\n\t    http://www.apache.org/licenses/LICENSE-2.0\n\n\tUnless required by applicable law or agreed to in writing, software\n\tdistributed under the License is distributed on an \"AS IS\" BASIS,\n\tWITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n\tSee the License for the specific language governing permissions and\n\tlimitations under the License.\n*/\n\npackage v1alpha1\n\nimport (\n\tmetav1 \"k8s.io/apimachinery/pkg/apis/meta/v1\"\n\n\txpv1 \"github.com/crossplane/crossplane-runtime/apis/common/v1\"\n)\n{{- .TypeDefs}}"
}
