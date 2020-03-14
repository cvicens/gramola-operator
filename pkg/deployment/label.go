package deployment

import (
	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
)

// GetAppServiceLabels returns a map with the labels we want for all AppService assets
func GetAppServiceLabels(cr *gramolav1alpha1.AppService, component string) (labels map[string]string) {
	labels = map[string]string{"app": "gramola", "component": component}
	return labels
}
