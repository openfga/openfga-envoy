module github.com/openfga/openfga-envoy/extauthz

go 1.22.6

require (
	github.com/envoyproxy/go-control-plane v0.12.1-0.20240621013728-1eb8caab5155
	github.com/openfga/go-sdk v0.3.5
	github.com/openfga/openfga v1.6.1-0.20240906222438-b8787d5f9d21
	github.com/stretchr/testify v1.9.0
	go.uber.org/zap v1.27.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240822170219-fc7c04adadcd
	google.golang.org/grpc v1.66.0
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.5

require (
	github.com/cncf/xds/go v0.0.0-20240423153145-555b57ec207b // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.4 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
