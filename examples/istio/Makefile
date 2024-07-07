#!/bin/bash

ISTIO_VERSION=1.22.0
ISTIOCTL=$(shell pwd)/istio-$(ISTIO_VERSION)/bin/istioctl
KIND=kind
KUBECONFIG=$(shell pwd)/.kube/config
KUBECTL=kubectl --kubeconfig $(KUBECONFIG)
OPENFGA_EXTAUTHZ_VERSION ?= dev

CLUSTER_NAME=openfga-istio-example

.PHONY: download-istioctl
download-istioctl:
	@curl -L https://istio.io/downloadIstio | ISTIO_VERSION=$(ISTIO_VERSION) sh -

.PHONY: install-istioctl
install-istioctl:
	@test -d istio-$(ISTIO_VERSION) || $(MAKE) download-istioctl

.PHONY: create-cluster
create-cluster:
	@$(KIND) create cluster --config kind-config.yaml --name $(CLUSTER_NAME)
	@$(KIND) export kubeconfig --name $(CLUSTER_NAME) --kubeconfig $(KUBECONFIG)

.PHONY: patch-mesh
patch-mesh: 
	@$(KUBECTL) get configmap istio -n istio-system -o jsonpath='{.data.mesh}' \
		| yq '.extensionProviders | map(select(.name == "openfga-ext-authz-grpc")) | length' \
		| awk '$$0 == "1" { print "Already patched"; exit 1 }'
	@$(KUBECTL) get configmap istio -n istio-system -o json > .istio-configmap.json.backup
	@cat .istio-configmap.json.backup | jq -r '.data.mesh' | yq -o=j > .istio-configmap-data.json
	@jq -r '.extensionProviders += [{"name": "openfga-ext-authz-grpc", "envoyExtAuthzGrpc": {"service": "ext-authz.foo.svc.cluster.local", "port": "9000"}}]' .istio-configmap-data.json > .istio-configmap-data.auth.json
	@cat .istio-configmap-data.auth.json | yq -P > .istio-configmap-data.auth.yaml
	@printf '"' > .istio-configmap-data.auth.oneline.yaml
	@sed 's/$$/\\n/' .istio-configmap-data.auth.yaml | sed 's/"//g' | tr -d '\n' | head -c -2 >> .istio-configmap-data.auth.oneline.yaml
	@printf '"' >> .istio-configmap-data.auth.oneline.yaml
	@$(MAKE) .patch-istio-configmap-data.yaml || ( rm .istio-configmap-data.auth.json*; exit 1)
	@rm .istio-configmap*
	
.patch-istio-configmap-data.yaml:
	@jq --argjson mesh "$$(cat .istio-configmap-data.auth.oneline.yaml)" '.data.mesh = $$mesh' .istio-configmap.json.backup > .istio-configmap.json.apply
	@$(KUBECTL) apply -f .istio-configmap.json.apply

install-istio:
	@$(ISTIOCTL) install --set profile=demo -c $(KUBECONFIG) -y

deploy: install-istio
	@$(KUBECTL) label namespace default istio-injection=enabled
	@docker image rm openfga-istio-example/openfga-extauthz:0.1 || true
	@docker tag gcr.io/openfga/openfga-extauthz:$(OPENFGA_EXTAUTHZ_VERSION) openfga-istio-example/openfga-extauthz:0.1
	@$(KIND) load docker-image --name $(CLUSTER_NAME) openfga-istio-example/openfga-extauthz:0.1
	@$(KUBECTL) apply -f ../../istio/deployment.yaml
	@$(MAKE) patch-mesh || true
	@$(KUBECTL) apply -f ./manifests/service-a.yaml

destroy-cluster:
	@$(KIND) delete cluster --name $(CLUSTER_NAME)

kubeconfig:
	@echo "Paste from clipboard in your terminal and hit 'return/enter'"
	@echo "export KUBECONFIG=\"$(KUBECONFIG)\"" | pbcopy
	