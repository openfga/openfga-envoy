PACKAGE_VERSION ?= dev
BUILD_OS ?= $(shell go env GOOS)
BUILD_ARCH ?= $(shell go env GOARCH)
DOCKER_PLATFORM = $(BUILD_OS)/$(BUILD_ARCH)

.PHONY: test
test:
	@go test -count=1 ./...

.PHONY: build
build:
	@mkdir -p ./build
	make build-platform BUILD_OS="linux" BUILD_ARCH="amd64"
	make build-platform BUILD_OS="linux" BUILD_ARCH="arm64"
	make build-platform BUILD_OS="darwin" BUILD_ARCH="amd64"
	make build-platform BUILD_OS="darwin" BUILD_ARCH="arm64"

build-platform:
	@mkdir -p ./build
	@GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) CGO_ENABLED=0 go build -o ./build/envoy-extauthz-$(BUILD_OS)-$(BUILD_ARCH) ./cmd/extauthz/main.go

.PHONY: docker
docker: build
	@docker build --platform=$(DOCKER_PLATFORM) -t gcr.io/openfga/openfga-extauthz:$(PACKAGE_VERSION) -f Dockerfile .

.PHONY: e2e
e2e: e2e-tools
	@./e2e/run.sh

e2e-tools:
	@which yq > /dev/null || go install github.com/mikefarah/yq/v4@latest
	@which fga > /dev/null || go install github.com/openfga/cli/cmd/fga@latest
