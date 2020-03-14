package appservice

import (
	"context"
	"fmt"

	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
	_deployment "github.com/redhat/gramola-operator/pkg/deployment"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	serviceName = "gateway"
)

// Reconciling Gateway
func (r *ReconcileAppService) reconcileGateway(instance *gramolav1alpha1.AppService) (reconcile.Result, error) {

	if result, err := r.addGateway(instance); err != nil {
		return result, err
	}

	// Success
	return reconcile.Result{}, nil
}

func (r *ReconcileAppService) addGateway(instance *gramolav1alpha1.AppService) (reconcile.Result, error) {
	deployment := _deployment.NewGatewayDeployment(instance, serviceName, instance.Namespace)
	if err := r.client.Create(context.TODO(), deployment); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		log.Info(fmt.Sprintf("Created %s Deployment", deployment.Name))
	}

	etherpadService := _deployment.NewService(instance, serviceName, instance.Namespace, []string{"http"}, []int32{9001})
	if err := r.client.Create(context.TODO(), etherpadService); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		log.Info(fmt.Sprintf("Created %s Service", etherpadService.Name))
	}

	etherpadRoute := _deployment.NewRoute(instance, "etherpad", instance.Namespace, "etherpad", 9001)
	if err := r.client.Create(context.TODO(), etherpadRoute); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		log.Info(fmt.Sprintf("Created %s Route", etherpadRoute.Name))
	}

	//Success
	return reconcile.Result{}, nil
}
