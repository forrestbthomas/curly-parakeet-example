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

// S3Reconciler reconciles a S3 object
type S3Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tbd.github.com,resources=s3s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tbd.github.com,resources=s3s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tbd.github.com,resources=s3s/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the S3 object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *S3Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx)
	log.Info("reconciling s3", "name", req.Name)

	// get actual product CR
	var s3 tbdv1.S3
	if err := r.Get(ctx, req.NamespacedName, &s3); err != nil {
		log.Error(err, "unable to fetch product")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if s3.Status.Condition != "created" {
		log.Info("creating S3 bucket", "name", s3.Spec.BucketName)
	}
	log.Info("bucket exists")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *S3Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tbdv1.S3{}).
		Complete(r)
}
