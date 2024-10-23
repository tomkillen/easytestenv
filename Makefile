.DEFAULT_GOAL: ;

### Config

ENVTEST_K8S_VERSION ?= 1.28.3

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: fmt vet setup-envtest ginkgo
	KUBEBUILDER_ASSETS="$(PWD)/$(shell $(SETUP_ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir bin -p path)" $(GINKGO) ./...

### Tools

# === ginkgo ===
GINKGO ?= bin/ginkgo
.PHONY: ginkgo
ginkgo: $(GINKGO)
$(GINKGO):
	GOBIN=$(PWD)/bin go install github.com/onsi/ginkgo/v2/ginkgo@v2.19.0

# === setup-envtest ===
SETUP_ENVTEST ?= bin/setup-envtest
.PHONY: setup-envtest
setup-envtest: $(SETUP_ENVTEST)
$(SETUP_ENVTEST):
	GOBIN=$(PWD)/bin go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest