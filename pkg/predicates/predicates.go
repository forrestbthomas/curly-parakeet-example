package predicates

import (
	"context"

	"github.com/curly-parakeet-example/pkg/util"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
)

func PassThrough(_ context.Context, _ client.Client, _ ctrl.Request) (bool, error) {
	return true, nil
}

func IsNew(ctx context.Context, client client.Client, req ctrl.Request) (bool, error) {
	log := logging.FromContext(ctx)
	product, err := util.GetProduct(ctx, client, req)
	if err != nil {
		log.Error(err, "error getting product")
		return false, err
	}

	if product.Status.Condition == "" {
		log.Info("product is new will set status to Creating", "product", product.Name)
		return true, nil
	}
	log.Info("pre-existing product, no-op", "product", product.Name)

	return false, err
}
