package util

import (
	"context"

	tbdv1 "github.com/curly-parakeet-example/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
)

type ReconcileAccumulator struct {
	Result ctrl.Result
	Err    error
	Exit   bool
}

type ReconcilerFn func() (ctrl.Result, error, bool)

func ReconcilerFoldl(fns []ReconcilerFn, foldFn func(acc ReconcileAccumulator, next ReconcilerFn) (ctrl.Result, error, bool), init ReconcileAccumulator) ReconcileAccumulator {
	for _, fn := range fns {
		res, err, exit := foldFn(init, fn)
		init.Result = res
		init.Err = err
		init.Exit = exit
	}
	return init
}

func GetProduct(ctx context.Context, client client.Client, req ctrl.Request) (tbdv1.Product, error) {
	log := logging.FromContext(ctx)
	var product tbdv1.Product
	err := client.Get(
		ctx,
		req.NamespacedName,
		&product,
	)
	if err == nil {
		log.Info("got product")
	}
	return product, err

}
