package easytestenv_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEasytestenv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Easytestenv Suite")
}
