/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"

	tbdv1 "github.com/curly-parakeet-example/api/v1"
)

// InfrastructureReconciler reconciles a Infrastructure object
type InfrastructureReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tbd.github.com,resources=infrastructures,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tbd.github.com,resources=infrastructures/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tbd.github.com,resources=infrastructures/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Infrastructure object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *InfrastructureReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx)
	log.Info("reconciling infra", "name", req.Name)

	// get actual infra CR
	var infra tbdv1.Infrastructure
	if err := r.Get(ctx, req.NamespacedName, &infra); err != nil {
		log.Error(err, "unable to fetch infra")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("should create the relevant child Infra CRs - like S3", "service type", infra.Spec.ServiceType)
	/*
		if err := r.Client.Create(ctx, &tbdv1.S3{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: fmt.Sprintf("%s-", req.Name),
				Namespace:    req.Namespace,
			},
			Spec: tbdv1.S3Spec{
				BucketName: "mycustombucket",
			},
		}); err != nil {
			log.Error(err, "could not create S3 custom resource")
		}
	*/

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InfrastructureReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tbdv1.Infrastructure{}).
		Complete(r)
}
