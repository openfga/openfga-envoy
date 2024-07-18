#!/bin/bash

ISTIO_VERSION=1.22.0
ISTIOCTL=$(shell pwd)/istio-$(ISTIO_VERSION)/bin/istioctl
KIND=kind
KUBECONFIG=$(shell pwd)/.kube/config
KUBECTL=kubectl --kubeconfig $(KUBECONFIG)
HELM=KUBECONFIG=$(KUBECONFIG) helm
OPENFGA_EXTAUTHZ_VERSION ?= dev
BUILD_ARCH?=arm64

CLUSTER_NAME=openfga-istio-example

all: install-tools install-istioctl create-cluster deploy

install-tools:
	@which kubectl > /dev/null || brew install kubectl
	@which kind > /dev/null || brew install kind
	@which jq > /dev/null || brew install jq
	@which helm > /dev/null || brew install helm

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

.PHONY: patch-envoy-extauthz
patch-envoy-extauthz: 
	@$(KUBECTL) get configmap istio -n istio-system -o jsonpath='{.data.mesh}' \
		| yq '.extensionProviders | map(select(.name == "openfga-ext-authz-grpc")) | length' \
		| awk '$$0 == "1" { print "Already patched"; exit 1 }'
	@$(KUBECTL) get configmap istio -n istio-system -o json > .istio-configmap.json.backup
	@cat .istio-configmap.json.backup | jq -r '.data.mesh' | yq -o=j > .istio-configmap-data.json
	@jq -r '.extensionProviders += [{"name": "openfga-extauthz-grpc", "envoyExtAuthzGrpc": {"service": "extauthz.openfga.svc.cluster.local", "port": "9000"}}]' .istio-configmap-data.json > .istio-configmap-data.auth.json
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
# follows https://medium.com/@s4l1h/how-to-install-kind-and-istio-ingress-controller-3b510834c762
	@$(KUBECTL) patch deployments.apps -o yaml -n istio-system istio-ingressgateway -p '{"spec":{"template":{"spec":{"containers":[{"name":"istio-proxy","ports":[{"containerPort":8080,"hostPort":80},{"containerPort":8443,"hostPort":443}]}]}}}}'

deploy: install-istio
	@$(MAKE) -C ../../extauthz docker BUILD_OS=linux BUILD_ARCH=$(BUILD_ARCH) DOCKER_VERSION=example-istio
	@docker tag ghcr.io/openfga/openfga-extauthz:example-istio ghcr.io/openfga/openfga-extauthz:0.0.1
	@$(KIND) load docker-image --name $(CLUSTER_NAME) ghcr.io/openfga/openfga-extauthz:0.0.1
	
	@docker rmi ghcr.io/jcchavezs/httpmole:0.0.1 > /dev/null || true
	@docker pull --platform=linux/$(BUILD_ARCH) ghcr.io/jcchavezs/httpmole:main
	@docker tag ghcr.io/jcchavezs/httpmole:main ghcr.io/jcchavezs/httpmole:0.0.1
	@$(KIND) load docker-image --name $(CLUSTER_NAME) ghcr.io/jcchavezs/httpmole:0.0.1

	@$(KUBECTL) apply -f ./manifests/apps

	@$(KUBECTL) apply -f ./manifests/extauthz
	@$(HELM) repo add openfga https://openfga.github.io/helm-charts > /dev/null
	@$(KUBECTL) get svc openfga -n openfga > /dev/null || $(HELM) install openfga openfga/openfga --namespace openfga

	@$(MAKE) patch-envoy-extauthz || true

destroy-cluster:
	@$(KIND) delete cluster --name $(CLUSTER_NAME)

kubeconfig:
	@echo "ACTION REQUIRED: Paste from clipboard in your terminal and hit 'return/enter'"
	@echo "export KUBECONFIG=\"$(KUBECONFIG)\"" | pbcopy
	