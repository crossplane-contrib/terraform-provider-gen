package v1alpha1

import (
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/provider-terraform-plugin/pkg/controller/terraform"
	ctrl "sigs.k8s.io/controller-runtime"
)

// SetupMyType adds a controller that reconciles MyType managed resources.
func ConfigureReconciler(mgr ctrl.Manager, l logging.Logger) error {
	name := managed.ControllerName(ServiceAccountGroupKind)
	r := managed.NewReconciler(mgr,
		resource.ManagedKind(ServiceAccountGroupVersionKind),
		managed.WithExternalConnecter(&terraform.Connector{KubeClient: mgr.GetClient()}),
		managed.WithInitializers(managed.NewNameAsExternalName(mgr.GetClient())),
		managed.WithLogger(l.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&ServiceAccount{}).
		Complete(r)
}
