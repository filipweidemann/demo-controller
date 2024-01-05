package controller

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ControllerManagerOptions struct {
	Scheme          *runtime.Scheme
	MetricsBindAddr string
	MetricsBindPort int
	ProbeBindAddr   string
	LeaderElection  bool
}

func CreateControllerManager(opts *ControllerManagerOptions) (ctrl.Manager, error) {
	leaderElectionNamespace := ""
	if testing.Testing() {
		leaderElectionNamespace = "kube-system"
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                  opts.Scheme,
		HealthProbeBindAddress:  opts.ProbeBindAddr,
		LeaderElection:          opts.LeaderElection,
		LeaderElectionID:        "pod-label-controller",
		LeaderElectionNamespace: leaderElectionNamespace,
	})

	return mgr, err
}
