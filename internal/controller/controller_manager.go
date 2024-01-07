package controller

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

type ControllerManagerOptions struct {
	K8sConfig       *rest.Config
	Scheme          *runtime.Scheme
	MetricsBindAddr string
	MetricsBindPort int
	ProbeBindAddr   string
	LeaderElection  bool
}

func CreateControllerManager(opts *ControllerManagerOptions) (ctrl.Manager, error) {
	ctrlOptions := ctrl.Options{
		Scheme:                  opts.Scheme,
		Metrics:                 server.Options{BindAddress: opts.MetricsBindAddr},
		HealthProbeBindAddress:  opts.ProbeBindAddr,
		LeaderElection:          true,
		LeaderElectionID:        "pod-labeller",
		LeaderElectionNamespace: "kube-system",
	}

	// When testing, we inject the config from the envtest cluster
	if opts.K8sConfig == nil {
		println("Config is null, generate default one")
		opts.K8sConfig = ctrl.GetConfigOrDie()
	}

	mgr, err := ctrl.NewManager(opts.K8sConfig, ctrlOptions)

	return mgr, err
}
