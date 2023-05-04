module github.com/aserto-dev/ds-load/plugins/auth0

go 1.19

replace github.com/aserto-dev/ds-load/common => ../../common

require (
	github.com/alecthomas/kong v0.7.1
	github.com/aserto-dev/ds-load/common v0.0.0-00010101000000-000000000000
	github.com/aserto-dev/go-grpc v0.8.56
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/aserto-dev/go-directory v0.20.8 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.54.0 // indirect
)
