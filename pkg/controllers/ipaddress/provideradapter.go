package ipaddress

import (
	"context"

	"github.com/rikatz/capi-phpipam/api/v1alpha1"
	"github.com/rikatz/capi-phpipam/pkg/ipamclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api-ipam-provider-in-cluster/pkg/ipamutil"
	ipampredicates "sigs.k8s.io/cluster-api-ipam-provider-in-cluster/pkg/predicates"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
)

// PHPIPAMProviderAdapter is used as middle layer for provider integration.
type PHPIPAMProviderAdapter struct {
	Client     client.Client
	IPAMClient *ipamclient.IPAMClient
}

var _ ipamutil.ProviderAdapter = &PHPIPAMProviderAdapter{}

// SetupWithManager sets up the controller with the Manager.
func (v *PHPIPAMProviderAdapter) SetupWithManager(_ context.Context, b *ctrl.Builder) error {
	b.
		For(&ipamv1.IPAddressClaim{}, builder.WithPredicates(
			ipampredicates.ClaimReferencesPoolKind(metav1.GroupKind{
				Group: v1alpha1.GroupVersion.Group,
				Kind:  v1alpha1.PHPIPAMPoolKind,
			}),
		)).
		WithOptions(controller.Options{
			// To avoid race conditions when allocating IP Addresses, we explicitly set this to 1
			MaxConcurrentReconciles: 1,
		}).
		Watches(
			&v1alpha1.PHPIPAMIPPool{},
			handler.EnqueueRequestsFromMapFunc(v.IPPoolToIPClaims()),
			builder.WithPredicates(resourceTransitionedToUnpaused()),
		).
		Owns(&ipamv1.IPAddress{}, builder.WithPredicates(
			ipampredicates.AddressReferencesPoolKind(metav1.GroupKind{
				Group: v1alpha1.GroupVersion.Group,
				Kind:  v1alpha1.PHPIPAMPoolKind,
			}),
		))
	return nil
}

// ClaimHandlerFor returns a claim handler for a specific claim.
func (v *PHPIPAMProviderAdapter) ClaimHandlerFor(_ client.Client, claim *ipamv1.IPAddressClaim) ipamutil.ClaimHandler {
	return &IPAddressClaimHandler{
		Client: v.Client,
		claim:  claim,
	}
}
