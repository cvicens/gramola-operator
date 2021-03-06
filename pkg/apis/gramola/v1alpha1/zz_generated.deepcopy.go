// +build !ignore_autogenerated

// You can add comments here...
// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppService) DeepCopyInto(out *AppService) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppService.
func (in *AppService) DeepCopy() *AppService {
	if in == nil {
		return nil
	}
	out := new(AppService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AppService) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppServiceCondition) DeepCopyInto(out *AppServiceCondition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppServiceCondition.
func (in *AppServiceCondition) DeepCopy() *AppServiceCondition {
	if in == nil {
		return nil
	}
	out := new(AppServiceCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppServiceList) DeepCopyInto(out *AppServiceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AppService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppServiceList.
func (in *AppServiceList) DeepCopy() *AppServiceList {
	if in == nil {
		return nil
	}
	out := new(AppServiceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AppServiceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppServiceSpec) DeepCopyInto(out *AppServiceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppServiceSpec.
func (in *AppServiceSpec) DeepCopy() *AppServiceSpec {
	if in == nil {
		return nil
	}
	out := new(AppServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppServiceStatus) DeepCopyInto(out *AppServiceStatus) {
	*out = *in
	in.ReconcileStatus.DeepCopyInto(&out.ReconcileStatus)
	if in.EventsDatabaseScriptRuns != nil {
		in, out := &in.EventsDatabaseScriptRuns, &out.EventsDatabaseScriptRuns
		*out = make([]DatabaseScriptRun, len(*in))
		copy(*out, *in)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]AppServiceCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppServiceStatus.
func (in *AppServiceStatus) DeepCopy() *AppServiceStatus {
	if in == nil {
		return nil
	}
	out := new(AppServiceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseScriptRun) DeepCopyInto(out *DatabaseScriptRun) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseScriptRun.
func (in *DatabaseScriptRun) DeepCopy() *DatabaseScriptRun {
	if in == nil {
		return nil
	}
	out := new(DatabaseScriptRun)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReconcileStatus) DeepCopyInto(out *ReconcileStatus) {
	*out = *in
	in.LastUpdate.DeepCopyInto(&out.LastUpdate)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReconcileStatus.
func (in *ReconcileStatus) DeepCopy() *ReconcileStatus {
	if in == nil {
		return nil
	}
	out := new(ReconcileStatus)
	in.DeepCopyInto(out)
	return out
}
