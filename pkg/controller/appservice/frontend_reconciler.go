package appservice

import (
	"context"
	"fmt"

	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
	_deployment "github.com/redhat/gramola-operator/pkg/deployment"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	frontendServiceName = "frontend"
)

// Reconciling Frontend
func (r *ReconcileAppService) reconcileFrontend(instance *gramolav1alpha1.AppService) (reconcile.Result, error) {

	if result, err := r.addFrontend(instance); err != nil {
		return result, err
	}

	// Success
	return reconcile.Result{}, nil
}

func (r *ReconcileAppService) addFrontend(instance *gramolav1alpha1.AppService) (reconcile.Result, error) {
	deployment := _deployment.NewFrontendDeployment(instance, frontendServiceName, instance.Namespace, []string{0: "events"})
	if err := controllerutil.SetControllerReference(instance, deployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.client.Create(context.TODO(), deployment); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		log.Info(fmt.Sprintf("Created %s Deployment", deployment.Name))
		r.recorder.Eventf(instance, "Normal", "Deployment Created", "Created %s Deployment", deployment.Name)
	}

	service := _deployment.NewService(instance, frontendServiceName, instance.Namespace, []string{"http"}, []int32{8080})
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.client.Create(context.TODO(), service); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		log.Info(fmt.Sprintf("Created %s Service", service.Name))
		r.recorder.Eventf(instance, "Normal", "Service Created", "Created %s Service", service.Name)
	}

	route := _deployment.NewRoute(instance, frontendServiceName, instance.Namespace, frontendServiceName, 8080)
	if err := controllerutil.SetControllerReference(instance, route, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.client.Create(context.TODO(), route); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		log.Info(fmt.Sprintf("Created %s Route", route.Name))
		r.recorder.Eventf(instance, "Normal", "Route Created", "Created %s Route", route.Name)
	}

	//Success
	return reconcile.Result{}, nil
}
