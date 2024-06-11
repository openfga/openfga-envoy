ISTIO_VERSION=1.22.0
ISTIOCTL=$(shell pwd)/istio-$(ISTIO_VERSION)/bin/istioctl
KIND=kind
KUBECONFIG=$(shell pwd)/.kube/config
KUBECTL=kubectl --kubeconfig $(KUBECONFIG)
OPENFGA_EXTAUTHZ_VERSION ?= dev

CLUSTER_NAME=openfga-istio-example

.PHONY: install
install-istioctl:
	curl -L https://istio.io/downloadIstio | ISTIO_VERSION=$(ISTIO_VERSION) sh -

.PHONY: create-cluster
create-cluster:
	@$(KIND) create cluster --config kind-config.yaml --name $(CLUSTER_NAME)
	@$(KIND) export kubeconfig --name $(CLUSTER_NAME) --kubeconfig $(KUBECONFIG)

.PHONY: patch-mesh
patch-mesh: .istio-configmap-data.yaml

.istio-configmap-data.yaml:
	@$(KUBECTL) get configmap istio -n istio-system -o json > .istio-configmap.json.backup
	@cat .istio-configmap.json.backup | jq -r '.data.mesh' > .istio-configmap-data.yaml
	@yq -o=j -I=0 '.extensionProviders += {"name": "openfga-ext-authz-grpc", "envoyExtAuthzGrpc": {"service": "ext-authz.foo.svc.cluster.local", "port": "9000"}}' .istio-configmap-data.yaml > .istio-configmap-data.auth.json
	@sed 's@"@\\"@g' .istio-configmap-data.auth.json > .istio-configmap-data.auth.jsonc
	@$(MAKE) .patch-istio-configmap-data.yaml || ( rm .istio-configmap-data.auth.json*; exit 1)
	@rm .istio-configmap-data.auth.json*
	@$(KUBECTL) get configmap istio -n istio-system -o yaml | grep openfga-ext-authz-grpc
	
.patch-istio-configmap-data.yaml:
	@test -f .istio-configmap-data.auth.jsonc
	@$(KUBECTL) patch configmap istio -n istio-system --type json \
	  -p '[{"op": "replace", "path": "/data/mesh", "value":"$(shell cat .istio-configmap-data.auth.jsonc)"}]'

.PHONY: deploy
deploy:
	@$(ISTIOCTL) install --set profile=demo -c $(KUBECONFIG) -y
	@$(KUBECTL) label namespace default istio-injection=enabled
	@docker image -rm openfga-istio-example/openfga-extauthz:1.0
	@docker tag gcr.io/openfga/openfga-extauthz:$(OPENFGA_EXTAUTHZ_VERSION) openfga-istio-example/openfga-extauthz:0.1
	@$(KIND) load docker-image openfga-istio-example/openfga-extauthz:0.1
	@$(KUBECTL) apply -f ../../istio/deployment.yaml
	@$(MAKE) patch-mesh
	@$(KUBECTL) apply -f ./manifests/service-a.yaml

destroy-cluster:
	@$(KIND) delete cluster --name $(CLUSTER_NAME)

kubeconfig:
	@echo "Paste from clipboard in your terminal and hit 'return/enter'"
	@echo "export KUBECONFIG=\"$(KUBECONFIG)\"" | pbcopy
	