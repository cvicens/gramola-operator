package deployment

import (
	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewPersistentVolumeClaim returns a PersistenceVolumeClaim given name, namespace, size, etc.
func NewPersistentVolumeClaim(cr *gramolav1alpha1.AppService, name string, namespace string, pvcClaimSize string) *corev1.PersistentVolumeClaim {
	labels := GetAppServiceLabels(cr, name)

	accessModes := []corev1.PersistentVolumeAccessMode{
		corev1.ReadWriteOnce,
	}
	resources := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(pvcClaimSize),
		}}
	pvcSpec := corev1.PersistentVolumeClaimSpec{
		AccessModes: accessModes,
		Resources:   resources,
	}
	return &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: pvcSpec,
	}

}
