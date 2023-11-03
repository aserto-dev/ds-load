module github.com/aserto-dev/ds-load/plugins/azuread

go 1.19

replace github.com/aserto-dev/ds-load/sdk => ../../sdk

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.6.1
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.3.0
	github.com/alecthomas/kong v0.7.1
	github.com/aserto-dev/ds-load/cli v0.0.0-20231020133027-e931efc4fad6
	github.com/aserto-dev/ds-load/sdk v0.0.0-00010101000000-000000000000
	github.com/aserto-dev/msgraph-sdk-go v0.0.1
	github.com/microsoft/kiota-abstractions-go v1.0.0
	github.com/microsoft/kiota-authentication-azure-go v1.0.0
	github.com/microsoft/kiota-http-go v1.0.0
	github.com/microsoft/kiota-serialization-json-go v1.0.2
	github.com/microsoftgraph/msgraph-sdk-go-core v1.0.0
	github.com/pkg/errors v0.9.1
	google.golang.org/grpc v1.58.3
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.31.0-20231017183020-0de7443d03cf.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.3.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.0.0 // indirect
	github.com/aserto-dev/clui v0.8.1 // indirect
	github.com/aserto-dev/go-directory v0.30.0 // indirect
	github.com/aserto-dev/logger v0.0.3 // indirect
	github.com/cjlapao/common-go v0.0.39 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dongri/phonenumber v0.0.0-20230428094603-9e2d44886294 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/kyokomi/emoji v2.2.4+incompatible // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/microsoft/kiota-serialization-form-go v1.0.0 // indirect
	github.com/microsoft/kiota-serialization-text-go v1.0.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/rs/zerolog v1.31.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	go.opentelemetry.io/otel v1.15.1 // indirect
	go.opentelemetry.io/otel/trace v1.15.1 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/term v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto v0.0.0-20231012201019-e917dd12ba7a // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231016165738-49dd2c1f3d0b // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231012201019-e917dd12ba7a // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
