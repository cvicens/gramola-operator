package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AppServiceSpec defines the desired state of AppService
type AppServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Flags if the the AppService object is enabled or not
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Enabled"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:booleanSwitch"
	Enabled bool `json:"enabled"`

	// Flags if the object has been initialized or not
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Initialized"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:booleanSwitch"
	Initialized bool `json:"initialized,omitempty"`

	// Different names for Gramola Service
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Alias"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text"
	// +kubebuilder:validation:Enum=Gramola;Gramophone;Phonograph
	Alias string `json:"alias,omitempty"`
}

// AppServiceConditionType defines the potential condition types
type AppServiceConditionType string

// AppServiceConditionTypes defined here
const (
	AppServiceConditionTypePromoted AppServiceConditionType = "Promoted"
)

// AppServiceConditionReason defines the potential condition reasons
type AppServiceConditionReason string

// AppServiceConditionReasons defined here
const (
	AppServiceConditionReasonInitialized AppServiceConditionReason = "Initialized"
	AppServiceConditionReasonWaiting     AppServiceConditionReason = "Waiting"
	AppServiceConditionReasonProgressing AppServiceConditionReason = "Progressing"
	AppServiceConditionReasonFinalising  AppServiceConditionReason = "Finalising"
	AppServiceConditionReasonSucceeded   AppServiceConditionReason = "Succeeded"
	AppServiceConditionReasonFailed      AppServiceConditionReason = "Failed"
)

// AppServiceConditionStatus defines the potential status
type AppServiceConditionStatus string

// AppServiceConditionStatuses defined here
const (
	AppServiceConditionStatusTrue    AppServiceConditionStatus = "True"
	AppServiceConditionStatusFalse   AppServiceConditionStatus = "False"
	AppServiceConditionStatusFailed  AppServiceConditionStatus = "Failed"
	AppServiceConditionStatusUnknown AppServiceConditionStatus = "Unknown"
)

// AppServiceCondition defines the desired state
type AppServiceCondition struct {
	// Type of replication controller condition.
	// +kubebuilder:validation:Enum=Promoted
	Type AppServiceConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=AppServiceConditionType"`
	// Status of the condition, one of True, False, Unknown.
	// +kubebuilder:validation:Enum=True;False;Unknown
	Status AppServiceConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,3,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition.
	// +optional
	// +kubebuilder:validation:Enum=Initialized;Waiting;Progressing;Finalising;Succeeded;Failed
	Reason AppServiceConditionReason `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// ReconcileStatus defines the reconciliation status
type ReconcileStatus struct {
	// Status shows the reconcile run
	// +kubebuilder:validation:Enum=Succeded;Progressing;Failed;True
	Status AppServiceConditionStatus `json:"status,omitempty"`
	// LastUpdate records the last time an update was regitered
	LastUpdate metav1.Time `json:"lastUpdate,omitempty"`
	// Reason for the update or change in status
	Reason string `json:"reason,omitempty"`
}

// ActionType defines the potential actions types
type ActionType string

// Action types defined here
const (
	BackupStarted ActionType = "BackupStarted"
	RequeueEvent  ActionType = "RequeueEvent"
	NoAction      ActionType = "NoAction"
)

// DatabaseUpdateStatus defines the potential status of a database update
type DatabaseUpdateStatus string

// DatabaseUpdateStatuses defined here
const (
	DatabaseUpdateStatusSucceeded DatabaseUpdateStatus = "Succeeded"
	DatabaseUpdateStatusFailed    DatabaseUpdateStatus = "Failed"
	DatabaseUpdateStatusUnknown   DatabaseUpdateStatus = "Unknown"
)

// DatabaseScriptRun logs script run and status
type DatabaseScriptRun struct {
	// Script
	Script string `json:"script"`

	// Status of the run of the Script
	// +kubebuilder:validation:Enum=Succeeded;Failed;Unknown
	Status DatabaseUpdateStatus `json:"eventsDatabaseUpdated,omitempty"`
}

// AppServiceStatus defines the observed state of AppService
type AppServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ReconcileStatus `json:",inline"`

	// Indicates if the Events Database has been updated or not
	// +kubebuilder:validation:Enum=Succeeded;Failed;Unknown
	EventsDatabaseUpdated DatabaseUpdateStatus `json:"eventsDatabaseUpdated,omitempty"`

	// List of Event Database Scripts Runs
	EventsDatabaseScriptRuns []DatabaseScriptRun `json:"eventsDatabaseScriptRuns,omitempty"`

	// Last Action run
	// +kubebuilder:validation:Enum=BackupStarted;NoAction;RequeueEvent
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="Last Action"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text"
	LastAction ActionType `json:"lastAction"`

	// Status Conditions
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="AppService Conditions"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes.conditions"
	Conditions []AppServiceCondition `json:"conditions,omitempty"` // Used to wait => kubectl wait canary/podinfo --for=condition=promoted
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppService is the Schema for the appservices API defines Gramola Backend Services
// +operator-sdk:gen-csv:customresourcedefinitions.displayName="AppService"
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=appservices,scope=Namespaced
type AppService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppServiceSpec   `json:"spec,omitempty"`
	Status AppServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppServiceList contains a list of AppService
type AppServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppService{}, &AppServiceList{})
}
