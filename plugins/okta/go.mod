module github.com/aseto-dev/ds-load/plugins/okta

go 1.19

replace github.com/aserto-dev/ds-load/common => ../../common

replace github.com/aserto-dev/ds-load/plugins/sdk => ../../plugins/sdk

replace github.com/aserto-dev/ds-load/plugins/okta => ../../plugins/okta

require (
	github.com/alecthomas/kong v0.7.1
	github.com/alecthomas/kong-yaml v0.1.1
	github.com/aserto-dev/ds-load/common v0.0.0-20230517082410-8615661d7c89
	github.com/aserto-dev/ds-load/plugins/okta v0.0.0-00010101000000-000000000000
	github.com/aserto-dev/ds-load/plugins/sdk v0.0.0-00010101000000-000000000000
	github.com/aserto-dev/go-grpc v0.8.56
	github.com/pkg/errors v0.9.1
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/aserto-dev/go-directory v0.21.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.54.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
