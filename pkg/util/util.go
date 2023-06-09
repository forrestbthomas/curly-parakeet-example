package util

import (
	"context"
	"fmt"

	tbdv1 "github.com/curly-parakeet-example/api/v1"
	errors "k8s.io/apimachinery/pkg/util/errors"
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

type ReconcilerFn func(ctx context.Context, client client.Client, req ctrl.Request) (ctrl.Result, error, bool)
type ReconcilerPredicate func(ctx context.Context, client client.Client, req ctrl.Request) (bool, error)
type ReconcilerPipeline map[ReconcilerTaskName]ReconcilerTask
type ReconcilerTaskName string

func ReconcilerFoldl(tasks ReconcilerPipeline, foldFn func(ReconcileAccumulator, ReconcilerTask) ReconcileAccumulator, init ReconcileAccumulator) ReconcileAccumulator {
	fmt.Println("*****")
	fmt.Println("FOLDING")

	nextTaskName := init.TaskToRun
	if nextTaskName == "" {
		nextTaskName = Input
	}

	nextTask := tasks[nextTaskName]
	init = foldFn(init, nextTask)
	out := init

	if len(nextTask.Next) > 0 {
		initCh := make(chan ReconcileAccumulator, len(nextTask.Next))
		for _, taskName := range nextTask.Next {
			go func(ch chan ReconcileAccumulator, t ReconcilerTaskName, ts ReconcilerPipeline) {
				init.TaskToRun = t
				ch <- foldFn(init, ts[t])
				close(ch)
			}(initCh, taskName, tasks)
		}

		for next := range initCh {
			out = mergeAccumulators(next, out)
		}
	}

	fmt.Println("FOLDED")
	fmt.Println("*****")
	return out
}

// mergeAccumulators will take the concurrent tasks and merge their accs, preferencing early exits, then errors, then result
func mergeAccumulators(prev, out ReconcileAccumulator) (res ReconcileAccumulator) {
	if prev.Exit || out.Exit {
		res.Exit = true
	}

	if prev.Err != nil || out.Err != nil {
		aggregateErr := errors.NewAggregate([]error{prev.Err, out.Err})
		res.Err = aggregateErr
	}

	if prev.Result.RequeueAfter > out.Result.RequeueAfter {
		res.Result = prev.Result
	} else {
		res.Result = out.Result
	}

	return
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
