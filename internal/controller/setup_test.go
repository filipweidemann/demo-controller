package controller

import (
	"log"
	"os"
	"testing"

	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var testEnv *envtest.Environment
var cfg *rest.Config
var k8sclient client.Client
var adminClientSet *clientset.Clientset
var podController *PodReconciler

var testContext context.Context
var testContextCancel context.CancelFunc

func packageSetup() {
	log.Print("PACKAGE SETUP")
	testContext, testContextCancel = context.WithCancel(context.Background())
	testEnv = &envtest.Environment{

		ErrorIfCRDPathMissing: false,
	}

	cfg, err := testEnv.Start()
	if err != nil {
		log.Fatalf("Failed to setup envtest, %v", err)
	}

	err = corev1.AddToScheme(scheme.Scheme)
	if err != nil {
		log.Fatalf("Couldn't add kind to scheme")
	}

	k8sclient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		log.Fatalf("Could not create k8s client: %v", err)
	}

	userOptions := envtest.User{Name: "integration-tests", Groups: []string{"system:masters"}}
	_, err = testEnv.ControlPlane.AddUser(userOptions, cfg)
	if err != nil {
		log.Fatalf("Could not create control plane user")
	}
	adminClientSet = clientset.NewForConfigOrDie(cfg)

	mgrOptions := ControllerManagerOptions{
		Scheme: scheme.Scheme,
	}

	mgr, err := CreateControllerManager(&mgrOptions)
	ctrl := &PodReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}

	if err := ctrl.SetupWithManager(mgr); err != nil {
		log.Fatalf(err.Error(), "unable to create controller", "controller", "Pod")
	}

	go mgr.Start(testContext)
}

func packageCleanup() {
	log.Print("PACKAGE CLEANUP")
	testContextCancel()
	testEnv.Stop()
}

func TestMain(m *testing.M) {
	packageSetup()
	rc := m.Run()
	packageCleanup()
	os.Exit(rc)
}
