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
	"crypto"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	logging "sigs.k8s.io/controller-runtime/pkg/log"

	tbdv1 "github.com/curly-parakeet-example/api/v1"
)

// ProductReconciler reconciles a Product object
type ProductReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tbd.github.com,resources=products,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tbd.github.com,resources=products/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tbd.github.com,resources=products/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Product object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *ProductReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx)
	log.Info("reconciling product", "name", req.Name)

	// get actual product CR
	var product tbdv1.Product
	if err := r.Get(ctx, req.NamespacedName, &product); err != nil {
		log.Error(err, "unable to fetch product")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.Client.Create(ctx, &tbdv1.Infrastructure{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", req.Name),
			Namespace:    req.Namespace,
		},
		Spec:   tbdv1.InfrastructureSpec{ServiceType: "objectStorage"},
		Status: tbdv1.InfrastructureStatus{},
	}); err != nil {
		log.Error(err, "could not create S3 custom resource")
	}

	return r.reconcileProduct(ctx, product)

}

func (r *ProductReconciler) reconcileProduct(ctx context.Context, product tbdv1.Product) (ctrl.Result, error) {
	err := errors.AggregateGoroutines(r.reconcileContainerOrchestration(ctx, product), r.reconcileInfrastructure(ctx, product))
	return ctrl.Result{}, err

}

func (r *ProductReconciler) reconcileContainerOrchestration(ctx context.Context, product tbdv1.Product) func() error {
	logging := log.FromContext(ctx)
	return func() (err error) {
		if product.Status.ObservedContainerOrchestration == "" {
			// create Container Orchestration CR and return nil or err
			logging.Info("would create a container orchestration CR")
			return
		}

		desiredContainerOrchestration, err := r.hash(product)
		if err != nil {
			return err
		}

		if product.Status.ObservedContainerOrchestration != desiredContainerOrchestration {
			// update product Status and return nil or err
			logging.Info("would update product status and return")
			return
		}

		return
	}
}

func (r *ProductReconciler) reconcileInfrastructure(ctx context.Context, product tbdv1.Product) func() error {
	logging := log.FromContext(ctx)
	return func() (err error) {
		if product.Status.ObservedInfrasturcture == "" {
			// create Infrastructure CR and return nil or err
			logging.Info("would create an infra CR")
		}

		desiredInfrastructure, err := r.hash(product)
		if err != nil {
			return err
		}

		if product.Status.ObservedInfrasturcture != desiredInfrastructure {
			// update product Status and return nil or err
			logging.Info("would update product status and return")
		}

		return
	}
}

func (r *ProductReconciler) hash(product tbdv1.Product) (s string, err error) {
	statusJson, err := json.Marshal(product)
	if err != nil {
		return "", err
	}

	hash := crypto.MD5.New()
	_, err = hash.Write(statusJson)
	if err != nil {
		return "", err
	}

	return string(hash.Sum(nil)), nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProductReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tbdv1.Product{}).
		Complete(r)
}
