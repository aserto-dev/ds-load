module github.com/aserto-dev/ds-load/plugins/cognito

go 1.19

replace github.com/aserto-dev/ds-load/sdk => ../../sdk

require github.com/aserto-dev/ds-load/sdk v0.0.0-00010101000000-000000000000

require (
	github.com/alecthomas/kong v0.7.1
	github.com/aws/aws-sdk-go v1.44.294
)

require github.com/jmespath/go-jmespath v0.4.0 // indirect

require (
	github.com/aserto-dev/go-directory v0.21.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230526203410-71b5a4ffd15e // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
