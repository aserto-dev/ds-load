module github.com/aserto-dev/ds-load/cli

go 1.23.7

toolchain go1.24.1

replace github.com/aserto-dev/ds-load/sdk => ../sdk

require (
	github.com/alecthomas/kong v1.10.0
	github.com/aserto-dev/ds-load/sdk v0.0.0-20250408143332-e8965667fcc0
	github.com/aserto-dev/go-aserto v0.33.8
	github.com/aserto-dev/go-directory v0.33.10
	github.com/fullstorydev/grpcurl v1.9.3
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.34.0
	github.com/stretchr/testify v1.10.0
	golang.org/x/sync v0.13.0
	google.golang.org/grpc v1.71.1
)

require (
	github.com/aserto-dev/errors v0.0.17 // indirect
	github.com/aserto-dev/header v0.0.11 // indirect
	github.com/aserto-dev/logger v0.0.9 // indirect
	github.com/bufbuild/protocompile v0.14.1 // indirect
	github.com/cncf/xds/go v0.0.0-20250326154945-ae57f3c0d45f // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/envoyproxy/go-control-plane/envoy v1.32.4 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/go-http-utils/headers v0.0.0-20181008091004-fed159eddc2a // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jhump/protoreflect v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/samber/lo v1.49.1 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250409194420-de1ac958c67a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250409194420-de1ac958c67a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
