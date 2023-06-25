package defintiions

import (
	"context"

	logging "sigs.k8s.io/controller-runtime/pkg/log"
)

type ProductUseCases map[string]func() error

func init() {

}

func getFnsFor(ctx context.Context, useCase string) {
	log := logging.FromContext(ctx)
	log.Info("getting reconcile functions for product use case", "useCase", useCase)

}
