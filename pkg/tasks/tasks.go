package tasks

import (
	"context"
	"time"

	"github.com/curly-parakeet-example/pkg/predicates"
	"github.com/curly-parakeet-example/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ctrl "sigs.k8s.io/controller-runtime"
	logging "sigs.k8s.io/controller-runtime/pkg/log"

	tbdv1 "github.com/curly-parakeet-example/api/v1"
)

var (
	Pipeline = map[util.ReconcilerTaskName]util.ReconcilerTask{

		util.Input: {
			Name:      util.Input,
			Fn:        SetStatusWith(tbdv1.Creating),
			Predicate: predicates.IsNew,
			Next: []util.ReconcilerTaskName{
				util.CheckStatus,
			},
		},

		util.CheckStatus: {
			Name:      util.CheckStatus,
			Fn:        CheckStatus,
			Predicate: predicates.PassThrough,
		},
	}
)

func SetStatusWith(status tbdv1.Condition) func(ctx context.Context, client client.Client, req ctrl.Request) (ctrl.Result, error, bool) {

	return func(ctx context.Context, client client.Client, req ctrl.Request) (ctrl.Result, error, bool) {
		log := logging.FromContext(ctx)
		product, err := util.GetProduct(ctx, client, req)
		if err != nil {
			log.Error(err, "error getting product before setting status")
			return ctrl.Result{}, err, true
		}

		product.Status = tbdv1.ProductStatus{
			Condition: status,
		}

		err = client.Status().Update(ctx, &product)
		if err != nil {
			log.Error(err, "error setting status", "product", product)
			return ctrl.Result{
				RequeueAfter: 1 * time.Minute,
			}, err, true
		}

		log.Info("product status was set successfully")

		return ctrl.Result{}, nil, false
	}
}

func CheckStatus(ctx context.Context, client client.Client, req ctrl.Request) (ctrl.Result, error, bool) {
	log := logging.FromContext(ctx)

	product, err := util.GetProduct(ctx, client, req)
	if err != nil {
		log.Error(err, "error getting product before setting status")
		return ctrl.Result{}, err, true
	}

	log.Info("retrieved product status", "status", product.Status.Condition)
	return ctrl.Result{}, nil, true

}
