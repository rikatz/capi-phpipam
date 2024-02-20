package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PHPIPAMIPPoolConditionType is the type for status conditions on
// IPPool resources. This type should be used with the
// PHPIPamIPPoolStatus.Conditions field.
type PHPIPAMIPPoolConditionType string

// PHPIPAMIPPoolConditionReason defines the set of reasons that explain why a
// particular PHPIPamIPPool condition type has been raised.
type PHPIPAMIPPoolConditionReason string

var (
	ConditionTypeReady            PHPIPAMIPPoolConditionType   = "Ready"
	ConditionReasonInvalidPHPIPam PHPIPAMIPPoolConditionReason = "InvalidPHPIPamConfiguration"
	ConditionReasonInvalidCreds   PHPIPAMIPPoolConditionReason = "InvalidPHPIPamCredentials"
	ConditionReasonIsReady        PHPIPAMIPPoolConditionReason = "IPPoolReady"
)

// PHPIPAMPoolSpec defines the desired state of PHPIPAMPoolSpec
type PHPIPAMPoolSpec struct {
	// SubnetID specifies the subnet that should be used on phpIPAM
	SubnetID int `json:"subnetid,omitempty"`

	// Credentials contains the credentials to connect to phpIPAM
	Credentials *PHPIPAMCredentials `json:"credentials,omitempty"`
}

type PHPIPAMCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	AppID    string `json:"app_id,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

// PHPIPamIPPoolStatus defines the observed state of PHPIPamIPPool
type PHPIPAMIPPoolStatus struct {
	// Conditions defines a set of Conditions of this IPPool
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Gateway is the reported Gateway from the subnet
	Gateway string `json:"gateway,omitempty"`

	// Mask is the reported Mask from the IPPool
	Mask string `json:"mask,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// PHPIPAMIPPool is the Schema for the PHPIPAMIPPool API
type PHPIPAMIPPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PHPIPAMPoolSpec     `json:"spec,omitempty"`
	Status PHPIPAMIPPoolStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PHPIPAMIPPoolList contains a list of PHPIPAMIPPool
type PHPIPAMIPPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PHPIPAMIPPool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PHPIPAMIPPool{}, &PHPIPAMIPPoolList{})
}
