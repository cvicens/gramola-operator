package appservice

import (
	"context"
	"fmt"

	routev1 "github.com/openshift/api/route/v1"
	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
	_deployment "github.com/redhat/gramola-operator/pkg/deployment"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	gatewayServiceName = "gateway"
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
	if gatewayDeployment, err := _deployment.NewGatewayDeployment(instance, r.scheme); err == nil {
		if err := r.client.Create(context.TODO(), gatewayDeployment); err != nil {
			if errors.IsAlreadyExists(err) {
				from := &appsv1.Deployment{}
				if err = r.client.Get(context.TODO(), types.NamespacedName{Name: gatewayDeployment.Name, Namespace: gatewayDeployment.Namespace}, from); err == nil {
					patch := _deployment.NewGatewayDeploymentPatch(from)
					if err := r.client.Patch(context.TODO(), from, patch); err != nil {
						return reconcile.Result{}, err
					}
				}
			} else {
				return reconcile.Result{}, err
			}
		}
		// Events Database Deployment created/updated successfully
		log.Info(fmt.Sprintf("Created/Updated %s Deployment", gatewayDeployment.Name))
		r.recorder.Eventf(instance, "Normal", "Deployment Created/Updated", "Created/Updated %s Deployment", gatewayDeployment.Name)
	} else {
		return reconcile.Result{}, err
	}

	if gatewayService, err := _deployment.NewGatewayService(instance, r.scheme); err == nil {
		if err := r.client.Create(context.TODO(), gatewayService); err != nil {
			if errors.IsAlreadyExists(err) {
				from := &corev1.Service{}
				if err = r.client.Get(context.TODO(), types.NamespacedName{Name: gatewayService.Name, Namespace: gatewayService.Namespace}, from); err == nil {
					patch := _deployment.NewGatewayServicePatch(from)
					if err := r.client.Patch(context.TODO(), from, patch); err != nil {
						return reconcile.Result{}, err
					}
				}
			} else {
				return reconcile.Result{}, err
			}
		}
		// Events Database Deployment created/updated successfully
		log.Info(fmt.Sprintf("Created/Updated %s Service", gatewayService.Name))
		r.recorder.Eventf(instance, "Normal", "Service Created/Updated", "Created/Updated %s Service", gatewayService.Name)
	} else {
		return reconcile.Result{}, err
	}

	//service := _deployment.NewService(instance, gatewayServiceName, instance.Namespace, []string{"http"}, []int32{8080})
	//if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
	//	return reconcile.Result{}, err
	//}
	//if err := r.client.Create(context.TODO(), service); err != nil && !errors.IsAlreadyExists(err) {
	//	return reconcile.Result{}, err
	//} else if err == nil {
	//	log.Info(fmt.Sprintf("Created %s Service", service.Name))
	//	r.recorder.Eventf(instance, "Normal", "Service Created", "Created %s Service", service.Name)
	//}

	if gatewayRoute, err := _deployment.NewGatewayRoute(instance, r.scheme); err == nil {
		if err := r.client.Create(context.TODO(), gatewayRoute); err != nil {
			if errors.IsAlreadyExists(err) {
				from := &routev1.Route{}
				if err = r.client.Get(context.TODO(), types.NamespacedName{Name: gatewayRoute.Name, Namespace: gatewayRoute.Namespace}, from); err == nil {
					patch := _deployment.NewGatewayRoutePatch(from)
					if err := r.client.Patch(context.TODO(), from, patch); err != nil {
						return reconcile.Result{}, err
					}
				}
			} else {
				return reconcile.Result{}, err
			}
		}
		// Events Database Deployment created/updated successfully
		log.Info(fmt.Sprintf("Created/Updated %s Route", gatewayRoute.Name))
		r.recorder.Eventf(instance, "Normal", "Route Created/Updated", "Created/Updated %s Route", gatewayRoute.Name)
	} else {
		return reconcile.Result{}, err
	}

	//route := _deployment.NewRoute(instance, gatewayServiceName, instance.Namespace, gatewayServiceName, 8080)
	//if err := controllerutil.SetControllerReference(instance, route, r.scheme); err != nil {
	//	return reconcile.Result{}, err
	//}
	//if err := r.client.Create(context.TODO(), route); err != nil && !errors.IsAlreadyExists(err) {
	//	return reconcile.Result{}, err
	//} else if err == nil {
	//	log.Info(fmt.Sprintf("Created %s Route", route.Name))
	//	r.recorder.Eventf(instance, "Normal", "Route Created", "Created %s Route", route.Name)
	//}

	//Success
	return reconcile.Result{}, nil
}
