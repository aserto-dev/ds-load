module github.com/aserto-dev/ds-load/plugins/auth0

go 1.19

replace github.com/aserto-dev/ds-load/sdk => ../../sdk

require (
	github.com/alecthomas/kong v0.7.1
	github.com/aserto-dev/ds-load/sdk v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	google.golang.org/protobuf v1.30.0
	gopkg.in/auth0.v5 v5.21.1
)

require (
	github.com/PuerkitoBio/rehttp v1.0.0 // indirect
	github.com/aserto-dev/go-directory v0.21.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/kr/text v0.2.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/oauth2 v0.7.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
