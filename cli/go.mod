module github.com/aserto-dev/ds-load/cli

go 1.23.7

toolchain go1.24.1

replace github.com/aserto-dev/ds-load/sdk => ../sdk

require (
	github.com/alecthomas/kong v1.8.1
	github.com/aserto-dev/ds-load/sdk v0.0.0-00010101000000-000000000000
	github.com/aserto-dev/go-aserto v0.33.5
	github.com/aserto-dev/go-directory v0.33.7
	github.com/bufbuild/protovalidate-go v0.7.3
	github.com/fullstorydev/grpcurl v1.9.2
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.33.0
	github.com/stretchr/testify v1.10.0
	golang.org/x/sync v0.12.0
	google.golang.org/grpc v1.70.0
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.35.2-20241127180247-a33202765966.1 // indirect
	cel.dev/expr v0.19.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/aserto-dev/header v0.0.8 // indirect
	github.com/aserto-dev/logger v0.0.7 // indirect
	github.com/bufbuild/protocompile v0.10.0 // indirect
	github.com/cncf/xds/go v0.0.0-20240905190251-b4127c9b8d78 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/envoyproxy/go-control-plane v0.13.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.1.0 // indirect
	github.com/go-http-utils/headers v0.0.0-20181008091004-fed159eddc2a // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/cel-go v0.22.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jhump/protoreflect v1.16.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/samber/lo v1.49.1 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	golang.org/x/exp v0.0.0-20241210194714-1829a127f884 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250303144028-a0af3efb3deb // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250303144028-a0af3efb3deb // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
