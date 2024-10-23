package easytestenv_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomkillen/easytestenv/easytestenv"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("create test env", func() {
	When("using a test environment", func() {
		env, err := easytestenv.New()
		It("should start successfully", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(env).ToNot(BeNil())
		})
		It("should create deployment", func() {
			Expect(env.ApplyResources(filepath.Join("..", "deploy", "deployment"))).To(Succeed())
		})
		It("should have created the namespace test-namespace", func() {
			var ns corev1.Namespace
			Expect(env.Client.Get(env.Context, types.NamespacedName{Name: "test-namespace"}, &ns)).To(Succeed())
			Expect(ns.Name).To(Equal("test-namespace"))
		})
		It("should have created the deployment test-namespace/nginx-deployment", func() {
			var d appsv1.Deployment
			Expect(env.Client.Get(env.Context, types.NamespacedName{Namespace: "test-namespace", Name: "nginx-deployment"}, &d)).To(Succeed())
			Expect(d.Name).To(Equal("nginx-deployment"))
		})
		It("should shutdown when finished", func() {
			env.Shutdown()
		})
	})
})
