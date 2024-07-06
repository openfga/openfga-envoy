.PHONY: test
test:
	@go test -count=1 ./...

.PHONY: build
build:
	@mkdir -p ./build
	@GOOS=linux CGO_ENABLED=0 go build -o ./build/envoy-extauthz-linux ./cmd/extauthz/main.go

.PHONY: docker
docker: build
	@docker build -t envoy-extauthz:dev -f Dockerfile .

.PHONY: e2e
e2e:
	@./e2e/run.sh
