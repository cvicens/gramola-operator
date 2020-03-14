package deployment

import (
	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewSecretFromStringData returns a ConfigMap given a stringData
func NewSecretFromStringData(cr *gramolav1alpha1.AppService, name string, namespace string, stringData map[string]string) *corev1.Secret {
	labels := GetAppServiceLabels(cr, name)
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		StringData: stringData,
	}
}

// NewSecretFromCrt return  a secret given a certificate
func NewSecretFromCrt(cr *gramolav1alpha1.AppService, name string, namespace string, crt []byte) *corev1.Secret {
	labels := GetAppServiceLabels(cr, name)
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: map[string][]byte{
			"ca.crt": crt,
		},
	}
}
