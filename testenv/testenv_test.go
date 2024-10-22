package testenv_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var testEnv *envtest.Environment
var k8sClient client.Client

func TestMain(m *testing.M) {
	// Set up a logger for controller-runtime
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Set up the test environment
	testEnv = &envtest.Environment{}

	// Start the environment
	cfg, err := testEnv.Start()
	if err != nil {
		fmt.Printf("Failed to start test environment: %v", err)
		return
	}

	// Add Kubernetes schemes for your test
	err = corev1.AddToScheme(scheme.Scheme)
	if err != nil {
		fmt.Printf("Failed to add CoreV1 scheme: %v", err)
		return
	}

	// Create a Kubernetes client to interact with the fake cluster
	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		fmt.Printf("Failed to create client: %v", err)
		return
	}

	// Run the tests
	code := m.Run()

	// Stop the test environment after tests are done
	err = testEnv.Stop()
	if err != nil {
		fmt.Printf("Failed to stop test environment: %v", err)
		return
	}

	// Exit the tests
	fmt.Printf("Test run completed with code %d\n", code)
}

func TestK8sInteraction(t *testing.T) {
	// Example: Create a namespace in the test Kubernetes environment
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-namespace",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := k8sClient.Create(ctx, ns); err != nil {
		t.Fatalf("Failed to create namespace: %v", err)
	}

	// Check if the namespace was created successfully
	createdNs := &corev1.Namespace{}
	err := k8sClient.Get(ctx, client.ObjectKey{Name: "test-namespace"}, createdNs)
	if err != nil {
		t.Fatalf("Failed to get namespace: %v", err)
	}

	t.Logf("Namespace %s created successfully", createdNs.Name)
}
