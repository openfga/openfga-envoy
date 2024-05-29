.PHONY: test
test:
	@go test -count=1 ./...

.PHONY: build
build:
	@mkdir -p ./build
	@GOOS=linux go build -o ./build/envoy-extauthz ./cmd/extauthz/main.go

.PHONY: docker
docker: build
	@docker build -t envoy-extauthz -f Dockerfile .

.PHONY: e2e
e2e:
	@./e2e/run.sh
