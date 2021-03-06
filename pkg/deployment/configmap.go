package deployment

import (
	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewConfigMapFromData returns a ConfigMap given a data object
func NewConfigMapFromData(cr *gramolav1alpha1.AppService, name string, namespace string, data map[string]string) *corev1.ConfigMap {
	labels := GetAppServiceLabels(cr, name)
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: data,
	}
}

// NewConfigMapPatchFromData returns a Patch given a ConfigMap and a KV map to add to th
//func NewConfigMapPatchFromCurrentAndKVMap(current *corev1.ConfigMap, data map[string]string) *client.Patch {
//	patch := client.MergeFrom(current.DeepCopy())
//	for k, v := range data {
//		current.Data[k] = v
//	}
//
//	return &patch
//}
