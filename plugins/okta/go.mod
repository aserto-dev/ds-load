module github.com/aserto-dev/ds-load/plugins/okta

go 1.19

replace github.com/aserto-dev/ds-load/sdk => ../../sdk

require (
	github.com/alecthomas/kong v0.7.1
	github.com/aserto-dev/ds-load/sdk v0.0.0-00010101000000-000000000000
	github.com/okta/okta-sdk-golang/v2 v2.18.0
	google.golang.org/grpc v1.56.1
)

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/aserto-dev/go-directory v0.21.0 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/patrickmn/go-cache v0.0.0-20180815053127-5633e0862627 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.10.0 // indirect
	golang.org/x/net v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230526203410-71b5a4ffd15e // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
