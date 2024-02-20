package phpipampool

import (
	"context"
	"errors"
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	ipamv1alpha1 "github.com/rikatz/capi-phpipam/api/v1alpha1"
	"github.com/rikatz/capi-phpipam/pkg/ipamclient"
)

var (
	ippoollogger logr.Logger
)

// PHPIPAMIPPoolReconciler reconciles a PHPIPAMIPPool object
type PHPIPAMIPPoolReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	ipamcl *ipamclient.IPAMClient
}

//+kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=phpipamippools,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=phpipamippools/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=phpipamippools/finalizers,verbs=update

// Reconcile the IPPool and set it as ready
func (r *PHPIPAMIPPoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	ippoollogger = log.FromContext(ctx).WithName("PHPIPAM ippool")
	ippoollogger.Info("received reconciliation", "request", req.NamespacedName.String())

	var ippool ipamv1alpha1.PHPIPAMIPPool
	if err := r.Get(ctx, req.NamespacedName, &ippool); err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		ippoollogger.Error(err, "unable to get ippool")
		return ctrl.Result{}, err
	}

	ipamcl, err := ipamclient.SpecToClient(&ippool.Spec)
	if err != nil {
		return r.ConditionsWithErrors(ctx, req, &ippool, ipamv1alpha1.ConditionReasonInvalidPHPIPam, "PHPIPAMconfig configuration is invalid: "+err.Error(), true)

	}

	subnetCfg, err := ipamcl.GetSubnetConfig()
	if err != nil {
		return r.ConditionsWithErrors(ctx, req, &ippool, ipamv1alpha1.ConditionReasonInvalidCreds, "failed to login to phpipam: "+err.Error(), false)
	}
	ippool.Status.Gateway = subnetCfg.Gateway.IPAddress
	ippool.Status.Mask = subnetCfg.Mask
	r.ipamcl = ipamcl
	return r.SetReady(ctx, req, &ippool)
}

// SetupWithManager sets up the controller with the Manager.
func (r *PHPIPAMIPPoolReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ipamv1alpha1.PHPIPAMIPPool{}).
		Complete(r)
}

func (r *PHPIPAMIPPoolReconciler) ConditionsWithErrors(ctx context.Context, req ctrl.Request, pool *ipamv1alpha1.PHPIPAMIPPool, reason ipamv1alpha1.PHPIPAMIPPoolConditionReason, errorMsg string, terminal bool) (ctrl.Result, error) {
	condition := metav1.Condition{
		Type:    string(ipamv1alpha1.ConditionTypeReady),
		Status:  metav1.ConditionFalse,
		Reason:  string(reason),
		Message: errorMsg,
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		var updatedPool ipamv1alpha1.PHPIPAMIPPool
		if err := r.Get(ctx, req.NamespacedName, &updatedPool); err != nil {
			return err
		}

		if updatedPool.Status.Conditions == nil {
			updatedPool.Status.Conditions = make([]metav1.Condition, 0)
		}
		meta.SetStatusCondition(&updatedPool.Status.Conditions, condition)
		return r.Status().Update(ctx, &updatedPool)
	})

	err := fmt.Errorf(errorMsg)
	if retryErr != nil {
		err = errors.Join(err, retryErr)
	}

	if terminal {
		err = reconcile.TerminalError(err)
	}
	return ctrl.Result{}, err

}

func (r *PHPIPAMIPPoolReconciler) SetReady(ctx context.Context, req ctrl.Request, pool *ipamv1alpha1.PHPIPAMIPPool) (ctrl.Result, error) {
	condition := metav1.Condition{
		Type:    string(ipamv1alpha1.ConditionTypeReady),
		Status:  metav1.ConditionTrue,
		Reason:  string(ipamv1alpha1.ConditionReasonIsReady),
		Message: "IPPool is ready",
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		var updatedPool ipamv1alpha1.PHPIPAMIPPool
		if err := r.Get(ctx, req.NamespacedName, &updatedPool); err != nil {
			return err
		}

		if updatedPool.Status.Conditions == nil {
			updatedPool.Status.Conditions = make([]metav1.Condition, 0)
		}
		updatedPool.Status.Gateway = pool.Status.Gateway
		updatedPool.Status.Mask = pool.Status.Mask
		meta.SetStatusCondition(&updatedPool.Status.Conditions, condition)
		return r.Status().Update(ctx, &updatedPool)
	})
	return ctrl.Result{}, retryErr
}
