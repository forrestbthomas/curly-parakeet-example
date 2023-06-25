package util

import (
	"context"

	tbdv1 "github.com/curly-parakeet-example/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Input       ReconcilerTaskName = "Input"
	CheckStatus ReconcilerTaskName = "CheckStatus"
)

type ReconcileAccumulator struct {
	TaskToRun   ReconcilerTaskName
	Result      ctrl.Result
	Err         error
	Exit        bool
	InputHasRun bool
}

type ReconcilerTask struct {
	Fn        ReconcilerFn
	Predicate ReconcilerPredicate
	Next      []ReconcilerTaskName
	Name      ReconcilerTaskName
}

type ReconcilerFn func() (ctrl.Result, error, bool)
type ReconcilerPredicate func() (bool, error)
type ReconcilerTaskName string

func ReconcilerFoldl(tasks map[ReconcilerTaskName]ReconcilerTask, foldFn func(ReconcileAccumulator, ReconcilerTask) ReconcileAccumulator, init ReconcileAccumulator) ReconcileAccumulator {
	// run input on every step, the Predicate will skip this step when needed: TODO - do this better
	init = foldFn(init, tasks[Input])
	for _, task := range tasks[Input].Next {
		init.TaskToRun = task
		init = foldFn(init, tasks[task])
		if init.Exit {
			return init
		}
	}
	return init
}

func GetProduct(ctx context.Context, client client.Client, req ctrl.Request) (tbdv1.Product, error) {
	var product tbdv1.Product
	err := client.Get(
		ctx,
		req.NamespacedName,
		&product,
	)
	return product, err

}
