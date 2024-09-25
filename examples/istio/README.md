Reference: <https://istio.io/latest/docs/reference/config/security/authorization-policy/>

Edit the mesh config with the following command:

$ kubectl edit configmap istio -n istio-system

In the editor, add the extension provider definitions shown below:

The following content defines two external providers openfga-ext-authz-grpc and openfga-ext-authz-http using the same service ext-authz.foo.svc.cluster.local. The service implements both the HTTP and gRPC check API as defined by the Envoy ext_authz filter. You will deploy the service in the following step.

data:
  mesh: |-
    # Add the following content to define the external authorizers.
    extensionProviders:
    - name: "openfga-ext-authz-grpc"
      envoyExtAuthzGrpc:
        service: "ext-authz.foo.svc.cluster.local"
        port: "9000"
    - name: "openfga-ext-authz-http"
      envoyExtAuthzHttp:
        service: "ext-authz.foo.svc.cluster.local"
        port: "8000"
        includeRequestHeadersInCheck: ["x-ext-authz"]
