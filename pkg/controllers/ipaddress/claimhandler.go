package ipaddress

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api-ipam-provider-in-cluster/pkg/ipamutil"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/rikatz/capi-phpipam/api/v1alpha1"
	"github.com/rikatz/capi-phpipam/pkg/ipamclient"
)

// IPAddressClaimHandler reconciles an IPAddress Claim getting the right address from the right pool
type IPAddressClaimHandler struct {
	client.Client
	claim   *ipamv1.IPAddressClaim
	mask    int
	gateway string
	ipamcl  *ipamclient.IPAMClient
}

var _ ipamutil.ClaimHandler = &IPAddressClaimHandler{}

// FetchPool fetches the PHPIPAM Pool.
func (h *IPAddressClaimHandler) FetchPool(ctx context.Context) (client.Object, *ctrl.Result, error) {

	var err error
	phpipampool := &v1alpha1.PHPIPAMIPPool{}

	if err = h.Client.Get(ctx, types.NamespacedName{Namespace: h.claim.Namespace, Name: h.claim.Spec.PoolRef.Name}, phpipampool); err != nil {
		return nil, nil, errors.Wrap(err, "failed to fetch pool")
	}

	if phpipampool.Status.Mask == "" || phpipampool.Status.Gateway == "" || !v1alpha1.PoolHasReadyCondition(phpipampool.Status) {
		return nil, nil, fmt.Errorf("IPPool is not ready yet")
	}

	h.mask, err = strconv.Atoi(phpipampool.Status.Mask)
	if err != nil {
		return nil, nil, fmt.Errorf("pool contains invalid network mask")
	}
	h.gateway = phpipampool.Status.Gateway

	ipamcl, err := ipamclient.SpecToClient(&phpipampool.Spec)
	if err != nil {
		return nil, nil, err
	}
	h.ipamcl = ipamcl

	return phpipampool, nil, nil
}

// EnsureAddress ensures that the IPAddress contains a valid address.
func (h *IPAddressClaimHandler) EnsureAddress(ctx context.Context, address *ipamv1.IPAddress) (*ctrl.Result, error) {
	hostname := fmt.Sprintf("%s.%s", h.claim.GetName(), h.claim.GetNamespace())
	ipv4, err := h.ipamcl.GetAddress(hostname)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get an IP Address")
	}

	address.Spec.Address = ipv4
	address.Spec.Gateway = h.gateway
	address.Spec.Prefix = h.mask
	return nil, nil
}

// ReleaseAddress releases the ip address.
func (h *IPAddressClaimHandler) ReleaseAddress(ctx context.Context) (*ctrl.Result, error) {
	hostname := fmt.Sprintf("%s.%s", h.claim.GetName(), h.claim.GetNamespace())
	err := h.ipamcl.ReleaseAddress(hostname)
	return nil, err
}
