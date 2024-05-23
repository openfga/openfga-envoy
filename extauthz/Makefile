.PHONY: build
build:
	@GOOS=linux go build -o build/envoy-extauthz cmd/extauthz/main.go

.PHONY: docker
docker: build
	@docker build -t envoy-extauthz -f Dockerfile .

.PHONY: e2e
e2e: build
	@docker compose -f e2e/docker-compose.yaml up -d --build
	sleep 2
	@curl -i localhost:8080 || ( docker compose -f e2e/docker-compose.yaml down && exit 1 )
	@docker compose -f e2e/docker-compose.yaml logs ext-authz > logs/ext-authz.log
	@docker compose -f e2e/docker-compose.yaml down
