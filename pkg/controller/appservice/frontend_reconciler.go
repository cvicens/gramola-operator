package appservice

import (
	"context"
	"fmt"

	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
	_deployment "github.com/redhat/gramola-operator/pkg/deployment"

	routev1 "github.com/openshift/api/route/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

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
	if frontendDeployment, err := _deployment.NewFrontendDeployment(instance, r.scheme); err == nil {
		if err := r.client.Create(context.TODO(), frontendDeployment); err != nil {
			if errors.IsAlreadyExists(err) {
				from := &appsv1.Deployment{}
				if err = r.client.Get(context.TODO(), types.NamespacedName{Name: frontendDeployment.Name, Namespace: frontendDeployment.Namespace}, from); err == nil {
					patch := _deployment.NewFrontendDeploymentPatch(from)
					if err := r.client.Patch(context.TODO(), from, patch); err != nil {
						return reconcile.Result{}, err
					}
				}
			} else {
				return reconcile.Result{}, err
			}
		}
		// Events Database Deployment created/updated successfully
		log.Info(fmt.Sprintf("Created/Updated %s Deployment", frontendDeployment.Name))
		r.recorder.Eventf(instance, "Normal", "Deployment Created/Updated", "Created/Updated %s Deployment", frontendDeployment.Name)
	} else {
		return reconcile.Result{}, err
	}

	if frontendService, err := _deployment.NewFrontendService(instance, r.scheme); err == nil {
		if err := r.client.Create(context.TODO(), frontendService); err != nil {
			if errors.IsAlreadyExists(err) {
				from := &corev1.Service{}
				if err = r.client.Get(context.TODO(), types.NamespacedName{Name: frontendService.Name, Namespace: frontendService.Namespace}, from); err == nil {
					patch := _deployment.NewFrontendServicePatch(from)
					if err := r.client.Patch(context.TODO(), from, patch); err != nil {
						return reconcile.Result{}, err
					}
				}
			} else {
				return reconcile.Result{}, err
			}
		}
		// Events Database Deployment created/updated successfully
		log.Info(fmt.Sprintf("Created/Updated %s Service", frontendService.Name))
		r.recorder.Eventf(instance, "Normal", "Service Created/Updated", "Created/Updated %s Service", frontendService.Name)
	} else {
		return reconcile.Result{}, err
	}

	if frontendRoute, err := _deployment.NewFrontendRoute(instance, r.scheme); err == nil {
		if err := r.client.Create(context.TODO(), frontendRoute); err != nil {
			if errors.IsAlreadyExists(err) {
				from := &routev1.Route{}
				if err = r.client.Get(context.TODO(), types.NamespacedName{Name: frontendRoute.Name, Namespace: frontendRoute.Namespace}, from); err == nil {
					patch := _deployment.NewFrontendRoutePatch(from)
					if err := r.client.Patch(context.TODO(), from, patch); err != nil {
						return reconcile.Result{}, err
					}
				}
			} else {
				return reconcile.Result{}, err
			}
		}
		// Events Database Deployment created/updated successfully
		log.Info(fmt.Sprintf("Created/Updated %s Route", frontendRoute.Name))
		r.recorder.Eventf(instance, "Normal", "Route Created/Updated", "Created/Updated %s Route", frontendRoute.Name)
	} else {
		return reconcile.Result{}, err
	}

	//Success
	return reconcile.Result{}, nil
}
