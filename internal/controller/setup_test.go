package controller

import (
	"context"
	"log"
	"os"
	"testing"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var testEnv *envtest.Environment
var cfg *rest.Config
var userConfig *rest.Config
var userClient client.Client
var k8sclient client.Client
var adminClientSet *clientset.Clientset

var testContext context.Context
var testContextCancel context.CancelFunc

func packageSetup() {
	testContext, testContextCancel = context.WithCancel(context.Background())
	testEnv = &envtest.Environment{}

	cfg, err := testEnv.Start()
	if err != nil {
		log.Fatalf("Failed to setup envtest, %v", err)
	}

	k8sclient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		log.Fatalf("Could not create k8s client: %v", err)
	}

	userOptions := envtest.User{Name: "integration-tests", Groups: []string{"system:masters"}}
	user, err := testEnv.ControlPlane.AddUser(userOptions, cfg)
	if err != nil {
		log.Fatalf("Could not create control plane user")
	}

	userConfig := user.Config()
	userClient, err = client.New(userConfig, client.Options{})
	if err != nil {
		log.Fatalf("Could not create userClient")
	}

	adminClientSet = clientset.NewForConfigOrDie(cfg)

	mgrOptions := ControllerManagerOptions{
		Scheme:    scheme.Scheme,
		K8sConfig: cfg,
	}

	mgr, err := CreateControllerManager(&mgrOptions)
	controller := &PodReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}

	if err := controller.SetupWithManager(mgr); err != nil {
		log.Fatalf(err.Error(), "unable to create controller", "controller", "Pod")
	}

	go func() {
		err := mgr.Start(testContext)
		if err != nil {
			println("An error occured inside the controller: ", err.Error())
		}
	}()
}

func packageCleanup() {
	testContextCancel()
	testEnv.Stop()
}

func TestMain(m *testing.M) {
	packageSetup()
	rc := m.Run()
	packageCleanup()
	os.Exit(rc)
}
