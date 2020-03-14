package appservice

import (
	"context"
	"fmt"
	"math"
	"time"

	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	record "k8s.io/client-go/tools/record"

	// Import GateWay

	// For now... blank
	_ "github.com/redhat/gramola-operator/pkg/util"
)

// Operator Name
const operatorName = "KharonOperator"

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

	// Watch for changes to primary resource AppService
	err = c.Watch(&source.Kind{Type: &gramolav1alpha1.AppService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner AppService
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
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
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Validate the CR instance
	if ok, err := r.IsValid(instance); !ok {
		return r.ManageError(instance, err)
	}

	//////////////////////////
	// Gateway
	//////////////////////////
	if _, err := r.reconcileGateway(instance); err != nil {
		return r.ManageError(instance, err)
	}

	// Define a new Pod object
	pod := newPodForCR(instance)

	// Set AppService instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return r.ManageSuccess(instance, 0, gramolav1alpha1.NoAction)
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *gramolav1alpha1.AppService) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

// IsValid checks if our CR is valid or not
func (r *ReconcileAppService) IsValid(obj metav1.Object) (bool, error) {
	//log.Info(fmt.Sprintf("IsValid? %s", obj))

	instance, ok := obj.(*gramolav1alpha1.AppService)
	if !ok {
		err := errors.NewBadRequest(errorNotAppServiceObject)
		log.Error(err, errorNotAppServiceObject)
		return false, err
	}

	// Check Alias
	if len(instance.Spec.Alias) > 0 && instance.Spec.Alias != "Gramola" && instance.Spec.Alias != "Gramophone" {
		err := errors.NewBadRequest(errorAlias)
		log.Error(err, errorAlias)
		return false, err
	}

	return true, nil
}

// ManageError manages an error object, an instance of the CR is passed along
func (r *ReconcileAppService) ManageError(obj metav1.Object, issue error) (reconcile.Result, error) {
	log.Error(issue, "Error managed")
	runtimeObj, ok := (obj).(runtime.Object)
	if !ok {
		log.Error(errors.NewBadRequest("not a runtime.Object"), "passed object was not a runtime.Object", "object", obj)
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
			Status:     gramolav1alpha1.AppServiceConditionStatusFailure,
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
		log.Error(errors.NewBadRequest("not a runtime.Object"), "passed object was not a runtime.Object", "object", obj)
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
