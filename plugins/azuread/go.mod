module github.com/aserto-dev/ds-load/plugins/azuread

go 1.19

replace github.com/aserto-dev/ds-load/sdk => ../../sdk

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.6.1
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.3.0
	github.com/alecthomas/kong v0.7.1
	github.com/aserto-dev/ds-load/sdk v0.0.0-00010101000000-000000000000
	github.com/aserto-dev/msgraph-sdk-go v0.0.1
	github.com/microsoft/kiota-authentication-azure-go v1.0.0
	github.com/microsoft/kiota-http-go v1.0.0
	github.com/microsoft/kiota-serialization-json-go v1.0.2
	github.com/pkg/errors v0.9.1
	google.golang.org/grpc v1.56.1
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.3.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.0.0 // indirect
	github.com/aserto-dev/go-directory v0.21.0 // indirect
	github.com/cjlapao/common-go v0.0.39 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/microsoft/kiota-abstractions-go v1.0.0 // indirect
	github.com/microsoft/kiota-serialization-form-go v1.0.0 // indirect
	github.com/microsoft/kiota-serialization-text-go v1.0.0 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	go.opentelemetry.io/otel v1.15.1 // indirect
	go.opentelemetry.io/otel/trace v1.15.1 // indirect
	golang.org/x/crypto v0.10.0 // indirect
	golang.org/x/net v0.11.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/text v0.10.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230526203410-71b5a4ffd15e // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
