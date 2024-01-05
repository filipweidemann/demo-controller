package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PodReconciler reconciles a Pod object
type PodReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	ANNOTATION = "demo-controller/label-pod"
	LABEL      = "demo-controller/podname"
)

func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, retErr error) {
	l := log.FromContext(ctx)

	pod := corev1.Pod{}
	err := r.Get(ctx, req.NamespacedName, &pod)

	// When Pod is deleted, we can silently let the reconciliation fail
	if apierrors.IsNotFound(err) {
		return res, nil
	}

	if err != nil {
		l.Error(err, "Unable to fetch pod")
		return res, err
	}

	labelShouldBePresent := pod.Annotations[ANNOTATION] == "true"
	labelIsPresent := pod.Labels[LABEL] == pod.Name

	if labelShouldBePresent == labelIsPresent {
		// State matches expectations.
		l.Info("No update required")
		return ctrl.Result{}, nil
	}

	if labelShouldBePresent {
		// If the label should be set but is not, set it.
		if pod.Labels == nil {
			pod.Labels = make(map[string]string)
		}
		pod.Labels[LABEL] = pod.Name
		l.Info("Adding label")
	} else {
		// If the label should not be set but is, remove it.
		delete(pod.Labels, LABEL)
		l.Info("Removing label")
	}

	err = r.Update(ctx, &pod)
	if apierrors.IsConflict(err) || apierrors.IsNotFound(err) {
		return ctrl.Result{Requeue: true}, nil
	}

	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.
		NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(r)
}
