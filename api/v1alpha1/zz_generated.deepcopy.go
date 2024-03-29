//go:build !ignore_autogenerated

/*
I am too lazy for a copyright :P
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PHPIPAMCredentials) DeepCopyInto(out *PHPIPAMCredentials) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PHPIPAMCredentials.
func (in *PHPIPAMCredentials) DeepCopy() *PHPIPAMCredentials {
	if in == nil {
		return nil
	}
	out := new(PHPIPAMCredentials)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PHPIPAMIPPool) DeepCopyInto(out *PHPIPAMIPPool) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PHPIPAMIPPool.
func (in *PHPIPAMIPPool) DeepCopy() *PHPIPAMIPPool {
	if in == nil {
		return nil
	}
	out := new(PHPIPAMIPPool)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PHPIPAMIPPool) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PHPIPAMIPPoolList) DeepCopyInto(out *PHPIPAMIPPoolList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PHPIPAMIPPool, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PHPIPAMIPPoolList.
func (in *PHPIPAMIPPoolList) DeepCopy() *PHPIPAMIPPoolList {
	if in == nil {
		return nil
	}
	out := new(PHPIPAMIPPoolList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PHPIPAMIPPoolList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PHPIPAMIPPoolStatus) DeepCopyInto(out *PHPIPAMIPPoolStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PHPIPAMIPPoolStatus.
func (in *PHPIPAMIPPoolStatus) DeepCopy() *PHPIPAMIPPoolStatus {
	if in == nil {
		return nil
	}
	out := new(PHPIPAMIPPoolStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PHPIPAMPoolSpec) DeepCopyInto(out *PHPIPAMPoolSpec) {
	*out = *in
	if in.Credentials != nil {
		in, out := &in.Credentials, &out.Credentials
		*out = new(PHPIPAMCredentials)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PHPIPAMPoolSpec.
func (in *PHPIPAMPoolSpec) DeepCopy() *PHPIPAMPoolSpec {
	if in == nil {
		return nil
	}
	out := new(PHPIPAMPoolSpec)
	in.DeepCopyInto(out)
	return out
}
