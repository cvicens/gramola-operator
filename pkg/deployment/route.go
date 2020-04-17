package deployment

import (
	routev1 "github.com/openshift/api/route/v1"
	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// NewRoute returns an OpenShift Route object
func NewRoute(cr *gramolav1alpha1.AppService, name string, namespace string, serviceName string, port int32) *routev1.Route {
	labels := GetAppServiceLabels(cr, name)
	targetPort := intstr.IntOrString{
		Type:   intstr.Int,
		IntVal: int32(port),
	}
	return &routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: routev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: serviceName,
			},
			Port: &routev1.RoutePort{
				TargetPort: targetPort,
			},
		},
	}
}
