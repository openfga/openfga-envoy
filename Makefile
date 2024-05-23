.PHONY: run
run:
	@go run extauthz/cmd/extauthz/main.go

.PHONY: test
test:
	@go test -count=1 ./...

.PHONY: build
build:
	@$(MAKE) -C extauthz build

.PHONY: e2e
e2e:
	@$(MAKE) -C extauthz e2e
