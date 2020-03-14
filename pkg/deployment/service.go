package deployment

import (
	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// NewService return a Service object given name, namespace, etc.
func NewService(cr *gramolav1alpha1.AppService, name string, namespace string, portName []string, portNumber []int32) *corev1.Service {
	labels := GetAppServiceLabels(cr, name)
	ports := []corev1.ServicePort{}
	for i := range portName {
		port := corev1.ServicePort{
			Name:     portName[i],
			Port:     portNumber[i],
			Protocol: "TCP",
		}
		ports = append(ports, port)
	}
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports:    ports,
			Selector: labels,
		},
	}
}

// NewCustomService returns a custom Service objec...
func NewCustomService(cr *gramolav1alpha1.AppService, name string, namespace string, portName []string, portNumber []int32, targetPortNumber []intstr.IntOrString) *corev1.Service {
	labels := GetAppServiceLabels(cr, name)
	ports := []corev1.ServicePort{}
	for i := range portName {
		port := corev1.ServicePort{
			Name:       portName[i],
			Port:       portNumber[i],
			TargetPort: targetPortNumber[i],
			Protocol:   "TCP",
		}
		ports = append(ports, port)
	}
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: ports,
		},
	}
}
