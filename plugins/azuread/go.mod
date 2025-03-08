module github.com/aserto-dev/ds-load/plugins/azuread

go 1.23.7

toolchain go1.24.1

replace github.com/aserto-dev/ds-load/sdk => ../../sdk

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.16.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.8.0
	github.com/alecthomas/kong v1.8.1
	github.com/aserto-dev/ds-load/sdk v0.0.0-20241206112725-5d200d771446
	github.com/aserto-dev/msgraph-sdk-go v0.0.3
	github.com/microsoft/kiota-abstractions-go v1.8.1
	github.com/microsoft/kiota-authentication-azure-go v1.1.0
	github.com/microsoft/kiota-http-go v1.4.7
	github.com/microsoft/kiota-serialization-json-go v1.0.9
	github.com/microsoftgraph/msgraph-sdk-go-core v1.2.1
	github.com/pkg/errors v0.9.1
	google.golang.org/grpc v1.70.0
)

require (
	dario.cat/mergo v1.0.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.3.2 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.3.1 // indirect
	github.com/Masterminds/sprig/v3 v3.3.0 // indirect
	github.com/aserto-dev/go-directory v0.33.7 // indirect
	github.com/aserto-dev/logger v0.0.7 // indirect
	github.com/cjlapao/common-go v0.0.41 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dongri/phonenumber v0.1.12 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/microsoft/kiota-serialization-form-go v1.0.0 // indirect
	github.com/microsoft/kiota-serialization-multipart-go v1.0.0 // indirect
	github.com/microsoft/kiota-serialization-text-go v1.0.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/std-uritemplate/std-uritemplate/go/v2 v2.0.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.opentelemetry.io/otel v1.32.0 // indirect
	go.opentelemetry.io/otel/metric v1.32.0 // indirect
	go.opentelemetry.io/otel/trace v1.32.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250303144028-a0af3efb3deb // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250303144028-a0af3efb3deb // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
