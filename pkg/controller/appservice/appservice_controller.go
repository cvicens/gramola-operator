package appservice

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
	_deployment "github.com/redhat/gramola-operator/pkg/deployment"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"k8s.io/client-go/kubernetes/scheme"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/remotecommand"

	// Route
	routev1 "github.com/openshift/api/route/v1"

	errors "github.com/pkg/errors"

	// For now... blank
	_ "github.com/redhat/gramola-operator/pkg/util"
)

// Operator Name
const operatorName = "gramola-operator"

// Best practices
const controllerName = "controller-appservice"

const (
	errorAlias                    = "Not a proper AppService object because TargetRef is not Deployment or DeploymentConfig"
	errorNotAppServiceObject      = "Not a AppService object"
	errorAppServiceObjectNotValid = "Not a valid AppService object"
	errorUnableToUpdateInstance   = "Unable to update instance"
	errorUnableToUpdateStatus     = "Unable to update status"
	errorUnexpected               = "Unexpected error"
)

var log = logf.Log.WithName(controllerName)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new AppService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	//return &ReconcileAppService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
	// Best practices
	return &ReconcileAppService{client: mgr.GetClient(), scheme: mgr.GetScheme(), recorder: mgr.GetEventRecorderFor(controllerName)}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// register OpenShift Routes in the scheme
	if err := routev1.AddToScheme(mgr.GetScheme()); err != nil {
		return err
	}

	appServicePredicate := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			log.Info("AppService (predicate->UpdateEvent) " + e.MetaNew.GetName())
			// Check that new and old objects are the expected type
			_, ok := e.ObjectOld.(*gramolav1alpha1.AppService)
			if !ok {
				log.Error(nil, "Update event has no old proper runtime object to update", "event", e)
				return false
			}
			newServiceConfig, ok := e.ObjectNew.(*gramolav1alpha1.AppService)
			if !ok {
				log.Error(nil, "Update event has no proper new runtime object for update", "event", e)
				return false
			}
			if !newServiceConfig.Spec.Enabled {
				log.Error(nil, "Runtime object is not enabled", "event", e)
				return false
			}

			// Also check if no change in ResourceGeneration to return false
			if e.MetaOld == nil {
				log.Error(nil, "Update event has no old metadata", "event", e)
				return false
			}
			if e.MetaNew == nil {
				log.Error(nil, "Update event has no new metadata", "event", e)
				return false
			}
			if e.MetaNew.GetGeneration() == e.MetaOld.GetGeneration() {
				return false
			}

			return true
		},
		CreateFunc: func(e event.CreateEvent) bool {
			log.Info("AppService (predicate->CreateFunc) " + e.Meta.GetName())
			_, ok := e.Object.(*gramolav1alpha1.AppService)
			if !ok {
				return false
			}

			return true
		},
	}

	// Watch for changes to primary resource AppService
	err = c.Watch(&source.Kind{Type: &gramolav1alpha1.AppService{}}, &handler.EnqueueRequestForObject{}, appServicePredicate)
	if err != nil {
		return err
	}

	podPredicate := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			log.Info("Pod (predicate->UpdateEvent) " + e.MetaNew.GetName())

			// Ignore if not events-database-*
			if !strings.Contains(e.MetaNew.GetName(), "events-database") {
				log.Info("Pod is not events-database - [IGNORED]")
				return false
			}

			// Check that new and old objects are the expected type
			_, ok := e.ObjectOld.(*corev1.Pod)
			if !ok {
				log.Error(nil, "Update event has no old proper runtime object to update", "event", e)
				return false
			}
			newPod, ok := e.ObjectNew.(*corev1.Pod)
			if !ok {
				log.Error(nil, "Update event has no proper new runtime object for update", "event", e)
				return false
			}
			if newPod.Status.Phase != corev1.PodRunning && newPod.Status.Phase != corev1.PodSucceeded {
				log.Info("Pod is not Running - [IGNORED]", "event", e.MetaNew.GetName())
				return false
			}

			return true
		},
		CreateFunc: func(e event.CreateEvent) bool {
			log.Info("Pod (predicate->CreateFunc) " + e.Meta.GetName() + "- [IGNORED]")
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			log.Info("Pod (predicate->DeleteEvent) " + e.Meta.GetName() + "- [IGNORED]")
			return false
		},
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner AppService
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &gramolav1alpha1.AppService{},
	}, podPredicate)
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner AppService
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &gramolav1alpha1.AppService{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileAppService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAppService{}

// ReconcileAppService reconciles a AppService object
type ReconcileAppService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	// Best practices...
	recorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a AppService object and makes changes based on the state read
// and what is in the AppService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAppService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AppService")

	// Fetch the AppService instance
	instance := &gramolav1alpha1.AppService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if k8s_errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	log.Info(fmt.Sprintf("Status %s", instance.Status))

	// Validate the CR instance
	if ok, err := r.isValid(instance); !ok {
		return r.ManageError(instance, err)
	}

	// Now that we have a target let's initialize the CR instance. Updates Spec.Initilized in `instance`
	if initialized, err := r.isInitialized(instance); err == nil && !initialized {
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			log.Error(err, errorUnableToUpdateInstance, "instance", instance)
			return r.ManageError(instance, err)
		}
		return reconcile.Result{}, nil
	} else {
		if err != nil {
			return r.ManageError(instance, err)
		}
	}

	//////////////////////////
	// Events
	//////////////////////////
	if _, err := r.reconcileEvents(instance); err != nil {
		return r.ManageError(instance, err)
	}

	//////////////////////////
	// Gateway
	//////////////////////////
	if _, err := r.reconcileGateway(instance); err != nil {
		return r.ManageError(instance, err)
	}

	//////////////////////////
	// Frontend
	//////////////////////////
	if _, err := r.reconcileFrontend(instance); err != nil {
		return r.ManageError(instance, err)
	}

	//////////////////////////
	// Update Events DataBase
	//////////////////////////
	log.Info(fmt.Sprintf("Status before Database Update %s", instance.Status))
	if !(instance.Status.EventsDatabaseUpdated == gramolav1alpha1.DatabaseUpdateStatusSucceeded) {
		if dataBaseUpdated, err := r.UpdateEventsDatabase(request); err != nil {
			log.Error(err, "error DB update", "instance", instance)
			// Update Status
			instance.Status.EventsDatabaseUpdated = gramolav1alpha1.DatabaseUpdateStatusFailed
			return r.ManageError(instance, err)
		} else {
			if dataBaseUpdated {
				log.Info(fmt.Sprintf("dataBaseUpdated ====> %s", instance.Status))
				// Update Status
				instance.Status.EventsDatabaseUpdated = gramolav1alpha1.DatabaseUpdateStatusSucceeded
				err := r.client.Update(context.Background(), instance)
				if err != nil {
					log.Error(err, errorUnableToUpdateInstance, "instance", instance)
					return r.ManageError(instance, err)
				}
			}
		}
	}

	// Nothing else to do
	return r.ManageSuccess(instance, 0, gramolav1alpha1.NoAction)
}

// isValid checks if our CR is valid or not
func (r *ReconcileAppService) isValid(obj metav1.Object) (bool, error) {
	//log.Info(fmt.Sprintf("isValid? %s", obj))

	instance, ok := obj.(*gramolav1alpha1.AppService)
	if !ok {
		err := k8s_errors.NewBadRequest(errorNotAppServiceObject)
		log.Error(err, errorNotAppServiceObject)
		return false, err
	}

	// Check Alias
	if len(instance.Spec.Alias) > 0 && instance.Spec.Alias != "Gramola" && instance.Spec.Alias != "Gramophone" {
		err := k8s_errors.NewBadRequest(errorAlias)
		log.Error(err, errorAlias)
		return false, err
	}

	return true, nil
}

// IsInitialized checks if our CR has been initialized or not
func (r *ReconcileAppService) isInitialized(obj metav1.Object) (bool, error) {
	instance, ok := obj.(*gramolav1alpha1.AppService)
	if !ok {
		err := k8s_errors.NewBadRequest(errorNotAppServiceObject)
		log.Error(err, errorNotAppServiceObject)
		return false, err
	}
	if instance.Spec.Initialized {
		return true, nil
	}

	// TODO add a Finalizer...
	// util.AddFinalizer(mycrd, controllerName)
	instance.Spec.Initialized = true
	return false, nil
}

// ManageError manages an error object, an instance of the CR is passed along
func (r *ReconcileAppService) ManageError(obj metav1.Object, issue error) (reconcile.Result, error) {
	log.Error(issue, "Error managed")
	runtimeObj, ok := (obj).(runtime.Object)
	if !ok {
		err := k8s_errors.NewBadRequest("not a runtime.Object")
		log.Error(err, "passed object was not a runtime.Object", "object", obj)
		r.recorder.Event(runtimeObj, "Error", "ProcessingError", err.Error())
		return reconcile.Result{}, nil
	}
	var retryInterval time.Duration
	r.recorder.Event(runtimeObj, "Warning", "ProcessingError", issue.Error())
	if instance, ok := (obj).(*gramolav1alpha1.AppService); ok {
		lastUpdate := instance.Status.LastUpdate
		lastStatus := instance.Status.Status
		status := gramolav1alpha1.ReconcileStatus{
			LastUpdate: metav1.Now(),
			Reason:     issue.Error(),
			Status:     gramolav1alpha1.AppServiceConditionStatusFailed,
		}
		instance.Status.ReconcileStatus = status
		err := r.client.Status().Update(context.Background(), runtimeObj)
		if err != nil {
			log.Error(err, errorUnableToUpdateStatus)
			return reconcile.Result{
				RequeueAfter: time.Second,
				Requeue:      true,
			}, nil
		}
		if lastUpdate.IsZero() || lastStatus == "Success" {
			retryInterval = time.Second
		} else {
			retryInterval = status.LastUpdate.Sub(lastUpdate.Time.Round(time.Second))
		}
	} else {
		log.Info("object is not RecocileStatusAware, not setting status")
		retryInterval = time.Second
	}
	return reconcile.Result{
		RequeueAfter: time.Duration(math.Min(float64(retryInterval.Nanoseconds()*2), float64(time.Hour.Nanoseconds()*6))),
		Requeue:      true,
	}, nil
}

// ManageSuccess manages a success and updates status accordingly, an instance of the CR is passed along
func (r *ReconcileAppService) ManageSuccess(obj metav1.Object, requeueAfter time.Duration, action gramolav1alpha1.ActionType) (reconcile.Result, error) {
	log.Info(fmt.Sprintf("===> ManageSuccess with requeueAfter: %d from: %s", requeueAfter, action))
	runtimeObj, ok := (obj).(runtime.Object)
	if !ok {
		log.Error(k8s_errors.NewBadRequest("not a runtime.Object"), "passed object was not a runtime.Object", "object", obj)
		return reconcile.Result{}, nil
	}
	if instance, ok := (obj).(*gramolav1alpha1.AppService); ok {
		status := gramolav1alpha1.ReconcileStatus{
			LastUpdate: metav1.Now(),
			Reason:     "",
			Status:     gramolav1alpha1.AppServiceConditionStatusTrue,
		}
		instance.Status.ReconcileStatus = status
		instance.Status.LastAction = action

		err := r.client.Status().Update(context.Background(), runtimeObj)
		if err != nil {
			log.Error(err, "Unable to update status")
			r.recorder.Event(runtimeObj, "Warning", "ProcessingError", "Unable to update status")
			return reconcile.Result{
				RequeueAfter: time.Second,
				Requeue:      true,
			}, nil
		}
		//if instance.Status.IsCanaryRunning {
		//	r.recorder.Event(runtimeObj, "Normal", "StatusUpdate", fmt.Sprintf("AppService in progress %d%%", instance.Status.CanaryWeight))
		//}
	} else {
		log.Info("object is not AppService, not setting status")
		r.recorder.Event(runtimeObj, "Warning", "ProcessingError", "Object is not AppService, not setting status")
	}

	if requeueAfter > 0 {
		return reconcile.Result{
			RequeueAfter: requeueAfter,
			Requeue:      true,
		}, nil
	}
	return reconcile.Result{}, nil
}

// UpdateEventsDatabase runs a script in the first 'Events' database pod found (and ready) returns true if the script was run succesfully
func (r *ReconcileAppService) UpdateEventsDatabase(request reconcile.Request) (bool, error) {
	// List all pods of the Events Database
	podList := &corev1.PodList{}
	lbs := map[string]string{
		"component": _deployment.EventsDatabaseServiceName,
	}
	labelSelector := labels.SelectorFromSet(lbs)
	listOps := &client.ListOptions{Namespace: request.Namespace, LabelSelector: labelSelector}
	if err := r.client.List(context.TODO(), podList, listOps); err != nil {
		return false, err
	}

	//log.Info(fmt.Sprintf("podList: %s", podList))

	// Count the pods that are pending or running as available
	var ready []corev1.Pod
	for _, pod := range podList.Items {
		log.Info(fmt.Sprintf("pod: %s phase: %s", pod.Name, pod.Status.Phase))
		if pod.Status.Phase == corev1.PodRunning {
			for _, containerStatus := range pod.Status.ContainerStatuses {
				if containerStatus.Name == _deployment.EventsDatabaseServiceName && containerStatus.Ready { // TODO constant for "postgresql"
					ready = append(ready, pod)
					break
				}
			}
		}
	}

	if len(ready) > 0 {
		filePath := _deployment.EventsDatabaseScriptsMountPath + "/" + _deployment.EventsDatabaseUpdateScriptName
		if _out, _err, err := r.ExecuteRemoteCommand(&ready[0], "psql -U $POSTGRESQL_USER $POSTGRESQL_DATABASE -f "+filePath); err != nil {
			return false, err
		} else {
			log.Info(fmt.Sprintf("stdout: %s\nstderr: %s", _out, _err))
			if len(_err) > 0 {
				return false, errors.Wrapf(err, "Failed executing script %s on %s", filePath, _deployment.EventsDatabaseServiceName)
			} else {
				return true, nil
			}
		}
	}

	return false, nil
}

// ExecuteRemoteCommand executes a remote shell command on the given pod
// returns the output from stdout and stderr
func (r *ReconcileAppService) ExecuteRemoteCommand(pod *corev1.Pod, command string) (string, string, error) {
	kubeCfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	restCfg, err := kubeCfg.ClientConfig()
	if err != nil {
		return "", "", err
	}
	coreClient, err := corev1client.NewForConfig(restCfg)
	if err != nil {
		return "", "", err
	}

	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	request := coreClient.RESTClient().
		Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: []string{"/bin/bash", "-c", command},
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", request.URL())
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: errBuf,
	})
	if err != nil {
		return "", "", errors.Wrapf(err, "Failed executing command %s on %v/%v", command, pod.Namespace, pod.Name)
	}

	return buf.String(), errBuf.String(), nil
}
