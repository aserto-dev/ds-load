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
	github.com/okta/okta-sdk-golang v1.1.0
	github.com/okta/okta-sdk-golang/v2 v2.18.0
	github.com/pkg/errors v0.9.1
	google.golang.org/grpc v1.54.0
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/aserto-dev/go-directory v0.21.0 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/go-yaml/yaml v2.1.0+incompatible // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/patrickmn/go-cache v0.0.0-20180815053127-5633e0862627 // indirect
	github.com/square/go-jose v2.4.1+incompatible // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
