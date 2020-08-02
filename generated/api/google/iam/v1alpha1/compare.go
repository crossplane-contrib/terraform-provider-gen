package v1alpha1

import (
	"fmt"

	xpresource "github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/terraform-provider-runtime/pkg/registry"
)

func forProviderDiffs(kube xpresource.Managed, prov xpresource.Managed) []string {
	diffs := make([]string, 0)
	k := kube.(*ServiceAccount)
	p := prov.(*ServiceAccount)
	if *k.Spec.ForProvider.Description != *p.Spec.ForProvider.Description {
		diffs = append(diffs, "description")
	}
	if *k.Spec.ForProvider.DisplayName != *p.Spec.ForProvider.DisplayName {
		diffs = append(diffs, "display_name")
	}
	return diffs
}

func atProviderDiffs(kube xpresource.Managed, prov xpresource.Managed) []string {
	diffs := make([]string, 0)
	k := kube.(*ServiceAccount)
	p := prov.(*ServiceAccount)
	if k.Status.AtProvider.Email != p.Status.AtProvider.Email {
		diffs = append(diffs, "email")
	}
	if k.Status.AtProvider.Name != p.Status.AtProvider.Name {
		diffs = append(diffs, "name")
	}
	if k.Status.AtProvider.Project != p.Status.AtProvider.Project {
		diffs = append(diffs, "project")
	}
	if k.Status.AtProvider.UniqueID != p.Status.AtProvider.UniqueID {
		diffs = append(diffs, "unique_id")
	}
	return diffs
}

func mergeResources(kube xpresource.Managed, prov xpresource.Managed) (xpresource.Managed, error) {
	k := kube.(*ServiceAccount)
	p := prov.(*ServiceAccount)
	k = k.DeepCopy()
	k.Status.AtProvider.Email = p.Status.AtProvider.Email
	k.Status.AtProvider.Name = p.Status.AtProvider.Name
	k.Status.AtProvider.Project = p.Status.AtProvider.Project
	k.Status.AtProvider.UniqueID = p.Status.AtProvider.UniqueID
	return k, nil
}

func diffIniter(kube xpresource.Managed, prov xpresource.Managed) (registry.ResourceDiff, error) {
	diffs := registry.ResourceDiff{}
	if _, ok := kube.(*ServiceAccount); !ok {
		return diffs, fmt.Errorf("v1alpha1.diffLister received a resource.Managed (left) value which is not a ServiceAccount.")
	}
	diffs.KubeResource = kube
	if _, ok := prov.(*ServiceAccount); !ok {
		return diffs, fmt.Errorf("v1alpha1.diffLister received a resource.Managed (left) value which is not a ServiceAccount.")
	}
	diffs.ProviderResource = prov
	diffs.AtProviderDiffCallback = atProviderDiffs
	diffs.ForProviderDiffCallback = forProviderDiffs
	diffs.MergeFunc = mergeResources
	return diffs, nil
}
