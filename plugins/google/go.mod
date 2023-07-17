module github.com/aserto-dev/ds-load/plugins/google

go 1.19

replace github.com/aserto-dev/ds-load/sdk => ../../sdk

require github.com/aserto-dev/ds-load/sdk v0.0.0-00010101000000-000000000000

require (
	github.com/alecthomas/kong v0.7.1
	golang.org/x/oauth2 v0.9.0
)

require (
	cloud.google.com/go/compute v1.19.3 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/aws/aws-sdk-go v1.44.294 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/s2a-go v0.1.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.5 // indirect
	github.com/googleapis/gax-go/v2 v2.11.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/grpc v1.56.1 // indirect
)

require (
	github.com/aserto-dev/go-directory v0.21.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	google.golang.org/api v0.130.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230629202037-9506855d4529 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
