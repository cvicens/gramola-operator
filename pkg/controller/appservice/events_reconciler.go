package appservice

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
	_deployment "github.com/redhat/gramola-operator/pkg/deployment"

	routev1 "github.com/openshift/api/route/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Constants to locate the scripts to update the database
const (
	DbScriptsBaseEnvVarName = "DB_SCRIPTS_BASE_DIR"
	//DbUpdateScriptName      = "events-database-update-0.0.2.sql"
	//DbScriptsMountPoint = "/operator/scripts"
)

// DbScriptsBasePath point to the directory where the scripts to update the database should be
var DbScriptsBasePath = os.Getenv(DbScriptsBaseEnvVarName) + "/db"

// Reconciling Events
func (r *ReconcileAppService) reconcileEvents(instance *gramolav1alpha1.AppService) (reconcile.Result, error) {

	if result, err := r.addEvents(instance); err != nil {
		return result, err
	}

	// Success
	return reconcile.Result{}, nil
}

func (r *ReconcileAppService) addEvents(instance *gramolav1alpha1.AppService) (reconcile.Result, error) {
	if databaseSecret, err := _deployment.NewEventsDatabaseCredentialsSecret(instance, r.scheme); err == nil {
		if err := r.client.Create(context.TODO(), databaseSecret); err != nil {
			if errors.IsAlreadyExists(err) {
				from := &corev1.Secret{}
				if err = r.client.Get(context.TODO(), types.NamespacedName{Name: databaseSecret.Name, Namespace: databaseSecret.Namespace}, from); err == nil {
					patch := _deployment.NewEventsDatabaseCredentialsSecretPatch(from)
					if err := r.client.Patch(context.TODO(), from, patch); err != nil {
						return reconcile.Result{}, err
					}
				}
			} else {
				return reconcile.Result{}, err
			}
		}
		// Secret created/updated successfully
		log.Info(fmt.Sprintf("Created/Updated %s Secret", databaseSecret.Name))
		r.recorder.Eventf(instance, "Normal", "Secret Created/Updated", "Created/Updated %s Secret", databaseSecret.Name)
	} else {
		return reconcile.Result{}, err
	}

	// Create Events Database Script ConfigMap
	if databaseScriptsConfigMap, err := _deployment.NewEventsDatabaseScriptsConfigMap(instance, r.scheme); err == nil {
		if err := r.client.Create(context.TODO(), databaseScriptsConfigMap); err != nil {
			if errors.IsAlreadyExists(err) {
				from := &corev1.ConfigMap{}
				if err = r.client.Get(context.TODO(), types.NamespacedName{Name: databaseScriptsConfigMap.Name, Namespace: databaseScriptsConfigMap.Namespace}, from); err == nil {
					patch := _deployment.NewEventsDatabaseScriptsConfigMapPatch(from)
					if err := r.client.Patch(context.TODO(), from, patch); err != nil {
						return reconcile.Result{}, err
					}
				}
			} else {
				return reconcile.Result{}, err
			}
		}
		// ConfigMap created/updated successfully
		log.Info(fmt.Sprintf("Created/Updated %s ConfigMap", databaseScriptsConfigMap.Name))
		r.recorder.Eventf(instance, "Normal", "ConfigMap Created/Updated", "Created/Updated %s ConfigMap", databaseScriptsConfigMap.Name)
	} else {
		return reconcile.Result{}, err
	}

	// PVC for Events Database
	databasePersistentVolumeClaim := _deployment.NewPersistentVolumeClaim(instance, _deployment.EventsDatabaseServiceName, instance.Namespace, "512Mi")
	if err := controllerutil.SetControllerReference(instance, databasePersistentVolumeClaim, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.client.Create(context.TODO(), databasePersistentVolumeClaim); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		log.Info(fmt.Sprintf("Created %s Persistent Volume Claim", databasePersistentVolumeClaim.Name))
		r.recorder.Eventf(instance, "Normal", "PVC Created", "Created %s Persistent Volume Claim", databasePersistentVolumeClaim.Name)
	}

	// Adds environment variables from the secret values passed and also mounts a volume with the configmap also passed in
	if databaseDeployment, err := _deployment.NewEventsDatabaseDeployment(instance, r.scheme); err == nil {
		if err := r.client.Create(context.TODO(), databaseDeployment); err != nil {
			if errors.IsAlreadyExists(err) {
				from := &appsv1.Deployment{}
				if err = r.client.Get(context.TODO(), types.NamespacedName{Name: databaseDeployment.Name, Namespace: databaseDeployment.Namespace}, from); err == nil {
					patch := _deployment.NewEventsDatabaseDeploymentPatch(from)
					if err := r.client.Patch(context.TODO(), from, patch); err != nil {
						return reconcile.Result{}, err
					}
				}
			} else {
				return reconcile.Result{}, err
			}
		}
		// Events Database Deployment created/updated successfully
		log.Info(fmt.Sprintf("Created/Updated %s Deployment", databaseDeployment.Name))
		r.recorder.Eventf(instance, "Normal", "Deployment Created/Updated", "Created/Updated %s Deployment", databaseDeployment.Name)
	} else {
		return reconcile.Result{}, err
	}

	if databaseService, err := _deployment.NewEventsDatabaseService(instance, r.scheme); err == nil {
		if err := r.client.Create(context.TODO(), databaseService); err != nil {
			if errors.IsAlreadyExists(err) {
				from := &corev1.Service{}
				if err = r.client.Get(context.TODO(), types.NamespacedName{Name: databaseService.Name, Namespace: databaseService.Namespace}, from); err == nil {
					patch := _deployment.NewEventsDatabaseServicePatch(from)
					if err := r.client.Patch(context.TODO(), from, patch); err != nil {
						return reconcile.Result{}, err
					}
				}
			} else {
				return reconcile.Result{}, err
			}
		}
		// Events Database Service created/updated successfully
		log.Info(fmt.Sprintf("Created/Updated %s Service", databaseService.Name))
		r.recorder.Eventf(instance, "Normal", "Service Created/Updated", "Created/Updated %s Service", databaseService.Name)
	} else {
		return reconcile.Result{}, err
	}

	if eventsDeployment, err := _deployment.NewEventsDeployment(instance, r.scheme); err == nil {
		if err := r.client.Create(context.TODO(), eventsDeployment); err != nil {
			if errors.IsAlreadyExists(err) {
				from := &appsv1.Deployment{}
				if err = r.client.Get(context.TODO(), types.NamespacedName{Name: eventsDeployment.Name, Namespace: eventsDeployment.Namespace}, from); err == nil {
					patch := _deployment.NewEventsDeploymentPatch(from)
					if err := r.client.Patch(context.TODO(), from, patch); err != nil {
						return reconcile.Result{}, err
					}
				}
			} else {
				return reconcile.Result{}, err
			}
		}
		// Events Database Deployment created/updated successfully
		log.Info(fmt.Sprintf("Created/Updated %s Deployment", eventsDeployment.Name))
		r.recorder.Eventf(instance, "Normal", "Deployment Created/Updated", "Created/Updated %s Deployment", eventsDeployment.Name)
	} else {
		return reconcile.Result{}, err
	}

	if eventsService, err := _deployment.NewEventsService(instance, r.scheme); err == nil {
		if err := r.client.Create(context.TODO(), eventsService); err != nil {
			if errors.IsAlreadyExists(err) {
				from := &corev1.Service{}
				if err = r.client.Get(context.TODO(), types.NamespacedName{Name: eventsService.Name, Namespace: eventsService.Namespace}, from); err == nil {
					patch := _deployment.NewEventsServicePatch(from)
					if err := r.client.Patch(context.TODO(), from, patch); err != nil {
						return reconcile.Result{}, err
					}
				}
			} else {
				return reconcile.Result{}, err
			}
		}
		// Events Database Deployment created/updated successfully
		log.Info(fmt.Sprintf("Created/Updated %s Service", eventsService.Name))
		r.recorder.Eventf(instance, "Normal", "Service Created/Updated", "Created/Updated %s Service", eventsService.Name)
	} else {
		return reconcile.Result{}, err
	}

	if eventsRoute, err := _deployment.NewEventsRoute(instance, r.scheme); err == nil {
		if err := r.client.Create(context.TODO(), eventsRoute); err != nil {
			if errors.IsAlreadyExists(err) {
				from := &routev1.Route{}
				if err = r.client.Get(context.TODO(), types.NamespacedName{Name: eventsRoute.Name, Namespace: eventsRoute.Namespace}, from); err == nil {
					patch := _deployment.NewEventsRoutePatch(from)
					if err := r.client.Patch(context.TODO(), from, patch); err != nil {
						return reconcile.Result{}, err
					}
				}
			} else {
				return reconcile.Result{}, err
			}
		}
		// Events Database Deployment created/updated successfully
		log.Info(fmt.Sprintf("Created/Updated %s Route", eventsRoute.Name))
		r.recorder.Eventf(instance, "Normal", "Route Created/Updated", "Created/Updated %s Route", eventsRoute.Name)
	} else {
		return reconcile.Result{}, err
	}

	//Success
	return reconcile.Result{}, nil
}

func readFile(fileName string) (string, error) {
	filePath := _deployment.DbScriptsBasePath + "/" + fileName
	log.Info(fmt.Sprintf("Reading file %s", fileName))
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("File reading error", err)
		return "", err
	}
	return string(data), nil
}
