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
	@GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) CGO_ENABLED=0 go build -o ./build/envoy-extauthz-$(BUILD_OS)-$(BUILD_ARCH) ./cmd/extauthz/main.go

.PHONY: docker
docker: build
	@docker build --platform=$(DOCKER_PLATFORM) -t gcr.io/openfga/openfga-extauthz:$(PACKAGE_VERSION) -f Dockerfile .

.PHONY: e2e
e2e:
	@./e2e/run.sh
