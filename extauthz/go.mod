module github.com/openfga/openfga-envoy/extauthz

go 1.22.1

require (
	github.com/envoyproxy/go-control-plane v0.12.1-0.20240419124334-0cebb2f428b3
	github.com/openfga/go-sdk v0.3.5
	github.com/stretchr/testify v1.9.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240415141817-7cd4c1c1f9ec
	google.golang.org/grpc v1.63.2
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.5

require (
	github.com/cncf/xds/go v0.0.0-20240329184929-0c46c01016dc // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.4 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
