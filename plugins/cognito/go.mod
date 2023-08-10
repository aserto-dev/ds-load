module github.com/aserto-dev/ds-load/plugins/cognito

go 1.19

replace github.com/aserto-dev/ds-load/sdk => ../../sdk

require github.com/aserto-dev/ds-load/sdk v0.0.0-00010101000000-000000000000

require (
	github.com/alecthomas/kong v0.7.1
	github.com/aserto-dev/ds-load/cli v0.0.0-20230808135510-07f5c2157f1f
	github.com/aws/aws-sdk-go v1.44.294
	github.com/pkg/errors v0.9.1
)

require github.com/jmespath/go-jmespath v0.4.0 // indirect

require (
	github.com/aserto-dev/clui v0.8.1 // indirect
	github.com/aserto-dev/go-directory v0.21.0 // indirect
	github.com/aserto-dev/logger v0.0.3 // indirect
	github.com/dongri/phonenumber v0.0.0-20230428094603-9e2d44886294 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/kyokomi/emoji v2.2.4+incompatible // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/rs/zerolog v1.28.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/term v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230526203410-71b5a4ffd15e // indirect
	google.golang.org/grpc v1.56.1 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
