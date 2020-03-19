package deployment

import (
	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
)

// Label Consts
const (
	repo = "https://github.com/cvicens/gramola"
	ref  = "master"
)

// GetEventsAnnotations returns a map with the annotations for Events
func GetEventsAnnotations(cr *gramolav1alpha1.AppService) (labels map[string]string) {
	annotations := map[string]string{
		"app.openshift.io/connects-to": "events-database",
		"app.openshift.io/vcs-ref":     ref,
		"app.openshift.io/vcs-uri":     repo,
	}
	return annotations
}

// GetGatewayAnnotations returns a map with the annotations for Gateway
func GetGatewayAnnotations(cr *gramolav1alpha1.AppService) (labels map[string]string) {
	annotations := map[string]string{
		"app.openshift.io/connects-to": "events",
		"app.openshift.io/vcs-ref":     ref,
		"app.openshift.io/vcs-uri":     repo,
	}
	return annotations
}

// GetFrontendAnnotations returns a map with the annotations for Gateway
func GetFrontendAnnotations(cr *gramolav1alpha1.AppService) (labels map[string]string) {
	annotations := map[string]string{
		"app.openshift.io/connects-to": "gateway",
		"app.openshift.io/vcs-ref":     ref,
		"app.openshift.io/vcs-uri":     repo,
	}
	return annotations
}
