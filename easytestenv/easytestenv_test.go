package easytestenv_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomkillen/easytestenv/easytestenv"
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
		It("should shutdown when finished", func() {
			env.Shutdown()
		})
	})
})
