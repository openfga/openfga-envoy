.PHONY: run
run-extauthz:
	@go run extauthz/cmd/extauthz/main.go

.PHONY: test
test:
	@$(MAKE) -C extauthz test

.PHONY: build
build:
	@$(MAKE) -C extauthz build

.PHONY: e2e
e2e:
	@$(MAKE) -C extauthz e2e
