module github.com/aserto-dev/ds-load/cli

go 1.19

replace github.com/aserto-dev/ds-load/sdk => ../sdk

require (
	github.com/alecthomas/kong v0.7.1
	github.com/aserto-dev/clui v0.8.1
	github.com/aserto-dev/ds-load/sdk v0.0.0-00010101000000-000000000000
	github.com/aserto-dev/go-aserto v0.20.3
	github.com/aserto-dev/go-directory v0.21.0
	github.com/aserto-dev/logger v0.0.3
	github.com/fullstorydev/grpcurl v1.8.7
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.28.0
	github.com/stretchr/testify v1.8.4
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1
	golang.org/x/sync v0.3.0
	google.golang.org/grpc v1.56.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.14.0 // indirect
	github.com/jhump/protoreflect v1.12.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/kyokomi/emoji v2.2.4+incompatible // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/net v0.11.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/term v0.9.0 // indirect
	golang.org/x/text v0.10.0 // indirect
	google.golang.org/genproto v0.0.0-20230626202813-9b080da550b3 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230526203410-71b5a4ffd15e // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230526203410-71b5a4ffd15e // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
