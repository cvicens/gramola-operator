package appservice

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
	_deployment "github.com/redhat/gramola-operator/pkg/deployment"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Events services names
const (
	EventsServiceName         = "events"
	EventsDatabaseServiceName = EventsServiceName + "-database"
)

// Constants to locate the scripts to update the database
const (
	DbScriptsBaseEnvVarName = "DB_SCRIPTS_BASE_DIR"
	DbUpdateScriptName      = "events-database-update-0.0.1.sql"
	DbScriptsMountPoint     = "/operator/scripts"
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
	databaseCredentials := map[string]string{
		"database-name":     "eventsdb",
		"database-password": "secret",
		"database-user":     "luke",
	}
	databaseSecret := _deployment.NewSecretFromStringData(instance, EventsDatabaseServiceName, instance.Namespace, databaseCredentials)
	if err := controllerutil.SetControllerReference(instance, databaseSecret, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.client.Create(context.TODO(), databaseSecret); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		log.Info(fmt.Sprintf("Created %s Secret", databaseSecret.Name))
		r.recorder.Eventf(instance, "Normal", "Secret Created", "Created %s Secret", databaseSecret.Name)
	}

	scripts := make(map[string]string)
	if dbUpdateScriptData, err := readFile(DbUpdateScriptName); err == nil {
		dbUpdateScriptDataReplaced := strings.Replace(dbUpdateScriptData, "{{DB_USERNAME}}", databaseCredentials["database-user"], -1)
		scripts[DbUpdateScriptName] = dbUpdateScriptDataReplaced
	}

	//log.Info(fmt.Sprintf("scripts %s", scripts))

	databaseConfigMap := _deployment.NewConfigMapFromData(instance, EventsDatabaseServiceName+"-scripts", instance.Namespace, scripts)
	if err := controllerutil.SetControllerReference(instance, databaseConfigMap, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.client.Create(context.TODO(), databaseConfigMap); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		log.Info(fmt.Sprintf("Created %s ConfigMap", databaseConfigMap.Name))
		r.recorder.Eventf(instance, "Normal", "ConfigMap Created", "Created %s ConfigMap", databaseConfigMap.Name)
	}

	databasePersistentVolumeClaim := _deployment.NewPersistentVolumeClaim(instance, EventsDatabaseServiceName, instance.Namespace, "512Mi")
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
	databaseDeployment := _deployment.NewEventsDatabaseDeployment(instance, EventsDatabaseServiceName, instance.Namespace, databaseSecret.Name, databaseConfigMap.Name, DbScriptsMountPoint)
	if err := controllerutil.SetControllerReference(instance, databaseDeployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.client.Create(context.TODO(), databaseDeployment); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		log.Info(fmt.Sprintf("Created %s Database", databaseDeployment.Name))
		r.recorder.Eventf(instance, "Normal", "Deployment Created", "Created %s Database", databaseDeployment.Name)
	}

	databaseService := _deployment.NewService(instance, EventsDatabaseServiceName, instance.Namespace, []string{"postgresql"}, []int32{5432})
	if err := controllerutil.SetControllerReference(instance, databaseService, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.client.Create(context.TODO(), databaseService); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		log.Info(fmt.Sprintf("Created %s Service", databaseService.Name))
		r.recorder.Eventf(instance, "Normal", "Service Created", "Created %s Service", databaseService.Name)
	}

	deployment := _deployment.NewEventsDeployment(instance, EventsServiceName, instance.Namespace, EventsDatabaseServiceName, EventsDatabaseServiceName, "5432")
	if err := controllerutil.SetControllerReference(instance, deployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.client.Create(context.TODO(), deployment); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		log.Info(fmt.Sprintf("Created %s Deployment", deployment.Name))
		r.recorder.Eventf(instance, "Normal", "Deployment Created", "Created %s Deployment", deployment.Name)
	}

	service := _deployment.NewService(instance, EventsServiceName, instance.Namespace, []string{"http"}, []int32{8080})
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.client.Create(context.TODO(), service); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		log.Info(fmt.Sprintf("Created %s Service", service.Name))
		r.recorder.Eventf(instance, "Normal", "Service Created", "Created %s Service", service.Name)
	}

	route := _deployment.NewRoute(instance, EventsServiceName, instance.Namespace, EventsServiceName, 8080)
	if err := controllerutil.SetControllerReference(instance, route, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.client.Create(context.TODO(), route); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		log.Info(fmt.Sprintf("Created %s Route", route.Name))
		r.recorder.Eventf(instance, "Normal", "Route Created", "Created %s Route", route.Name)
	}

	//Success
	return reconcile.Result{}, nil
}

func readFile(fileName string) (string, error) {
	filePath := DbScriptsBasePath + "/" + fileName
	log.Info(fmt.Sprintf("Reading file %s", fileName))
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("File reading error", err)
		return "", err
	}
	return string(data), nil
}
