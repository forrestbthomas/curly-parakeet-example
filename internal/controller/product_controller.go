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
	"github.com/curly-parakeet-example/pkg/tasks"
	"github.com/curly-parakeet-example/pkg/util"
)

// ProductReconciler reconciles a Product object
type ProductReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	req    ctrl.Request
	ctx    context.Context
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
	r.ctx = ctx
	r.req = req

	// get actual product CR
	var product tbdv1.Product
	if err := r.Get(ctx, req.NamespacedName, &product); err != nil {
		log.Error(err, "unable to fetch product")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// when would updates occur for a Product CR? the only spec field is useCase and that wouldn't change
	// there will also be a hash of dependent specs, and that can change
	// that means that we want the product controller to get the spec of the downstream objects and hash it and then compare that to the hash in the product spec
	// the issue with that is that the controller will now be updating the spec, and that's not good
	// i think what this means is that the only thing in the product spec should be the use case - if the use case changes, then that should be rejected, the only operation on a Product should be create/delete
	// downstream specs should be hashed and bubbled as a status field in the product CR
	// this way a human operator is required to delete a product, which would in turn delete downstream objects, via a finalizer

	// Product Controller Responsibility - create CRs that own the low-level objects for products (delete only via Finalizers), create/update Product Status
	// check if the object is a product, or an infra or a container cr
	// if product:
	// - if new, create then update status and return (it is new if the Status section is nil)
	// - if exists, check statuses,top level fields of downstream object name, the value being an object with fields of Phase and Hash,
	//   - if hash has changed, update and return
	//   - if Phase has changed, update and return (might be different depending on what the Phase is)

	/*

			DAG

			                                   isNew
		                                     │
		                               Yes   │ No
		      Out    ◄────── SetStatus  ◄───────┴──────►  CheckStatus
		                                                 │    │
		                                                 │    │
		                                                 │    │
		                                                 │    │
		                                    ProductStatus│    │DownstreamStatus
		                                                 │    │
		                           No                    │    │                 No
		                   Out   ◄────────hasChanged  ◄──┘    └──►  hasChanged ───────►  Out
		                                      │                          │
		                                      │                          │
		                                  Yes │                          │
		                                      │                          │
		                       │ ◄────────────┤            ◄───────┬─────┴─────┬─────────►
		                       │   Creating   │                    │           │
		                       │              │              Infra │           │Container
		            setStatus  │              │                    │           │
		Out     ◄──────────────┤ ◄────────────┤                    │           │
		                       │  Completed   │         Creating   │           │Creating
		                       │              │       │    ◄───────┤           ├────────►  │
		                       │              │       │            │           │           │
		                       │ ◄────────────┘       │ Completed  │           │Completed  │
		                          Failed              │    ◄───────┤           ├────────►  │
		                            │                 │            │           │           │
		                            │                 │ Failed     │           │Failed     │ setStatus
		                            │                 │    ◄───────┤           ├────────►  │
		                            │                 │            │           │           │
		                            │                 │            │           │           │
		                            ▼                 │            │           │           └───────────►   Out
		                        Alert         setStatus            │           │
		                                              │            │           │
		                                              │            ▼           ▼
		                                              │
		                                              │          Alert        Alert
		                                              │
		                           Out    ◄───────────┘

	*/

	// TODO: provide this as a default `fold` in the lib, but allow providing custom - hint: interfaces are hot
	foldFn := func(acc util.ReconcileAccumulator, nextTask util.ReconcilerTask) util.ReconcileAccumulator {
		var result ctrl.Result
		var err error
		var exit bool
		ok, err := nextTask.Predicate(ctx, r.Client, req)
		if err != nil {
			acc.Err = err
			return acc
		}

		if ok {
			result, err, exit = nextTask.Fn(ctx, r.Client, req)
			acc.Result = result
			acc.Err = err
			acc.Exit = exit
		}

		return acc
	}

	// folding instead of mapping so that we can track state in the accumulator
	result := util.ReconcilerFoldl(
		tasks.Pipeline,
		foldFn,
		util.ReconcileAccumulator{},
	)

	return result.Result, result.Err
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProductReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tbdv1.Product{}).
		Complete(r)
}
