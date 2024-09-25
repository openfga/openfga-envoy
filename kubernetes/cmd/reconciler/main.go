package main

import (
	"context"
	"log"

	"k8c.io/reconciler/pkg/reconciling"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func namespaceReconcilerFactory() (name string, reconciler reconciling.NamespaceReconciler) {
	return "openfga-k8s", func(cm *corev1.Namespace) (*corev1.Namespace, error) {
		return cm, nil
	}
}

func main() {
	var client client.Client

	nsReconcilers := []reconciling.NamedNamespaceReconcilerFactory{namespaceReconcilerFactory}
	if err := reconciling.ReconcileNamespaces(context.Background(), nsReconcilers, "", client); err != nil {
		log.Fatalf("Failed to reconcile namespace: %v", err)
	}

	func (r *TenantReconciler) SetupWithManager(mgr ctrl.Manager) error {
		return ctrl.NewControllerManagedBy(mgr).
			For(&kubelbv1alpha1.Tenant{}).
			Complete(r)
	}
}
