module github.com/aserto-dev/ds-load/cli

go 1.21

replace github.com/aserto-dev/ds-load/sdk => ../sdk

require (
	github.com/alecthomas/kong v0.9.0
	github.com/aserto-dev/clui v0.8.3
	github.com/aserto-dev/ds-load/sdk v0.0.0-00010101000000-000000000000
	github.com/aserto-dev/go-aserto v0.31.2
	github.com/aserto-dev/go-directory v0.31.3
	github.com/aserto-dev/logger v0.0.4
	github.com/bufbuild/protovalidate-go v0.6.1
	github.com/fullstorydev/grpcurl v1.9.1
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.32.0
	github.com/stretchr/testify v1.9.0
	golang.org/x/exp v0.0.0-20240409090435-93d18d7e34b8
	golang.org/x/sync v0.7.0
	google.golang.org/grpc v1.63.2
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.33.0-20240401165935-b983156c5e99.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/aserto-dev/header v0.0.6 // indirect
	github.com/bufbuild/protocompile v0.10.0 // indirect
	github.com/cncf/xds/go v0.0.0-20231128003011-0fa0005c9caa // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/envoyproxy/go-control-plane v0.12.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.4 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/go-http-utils/headers v0.0.0-20181008091004-fed159eddc2a // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/cel-go v0.20.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/jhump/protoreflect v1.16.0 // indirect
	github.com/kyokomi/emoji v2.2.4+incompatible // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/samber/lo v1.39.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/term v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto v0.0.0-20240415151819-79826c84ba32 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240415151819-79826c84ba32 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240415151819-79826c84ba32 // indirect
	google.golang.org/protobuf v1.33.1-0.20240408130810-98873a205002 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
