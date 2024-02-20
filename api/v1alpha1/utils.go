package v1alpha1

import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func PoolHasReadyCondition(status PHPIPAMIPPoolStatus) bool {
	for i := range status.Conditions {
		if status.Conditions[i].Type == string(ConditionTypeReady) && status.Conditions[i].Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}
