module github.com/aserto-dev/ds-load/plugins/auth0

go 1.19

replace github.com/aserto-dev/ds-load/common => ../../common

replace github.com/aserto-dev/ds-load/plugins/sdk => ../../plugins/sdk

require (
	github.com/alecthomas/kong v0.7.1
	github.com/aserto-dev/ds-load/common v0.0.0-20230517082410-8615661d7c89
	github.com/aserto-dev/ds-load/plugins/sdk v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/protobuf v1.30.0
	gopkg.in/auth0.v5 v5.21.1
)

require (
	github.com/PuerkitoBio/rehttp v1.0.0 // indirect
	github.com/aserto-dev/go-directory v0.21.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/oauth2 v0.4.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
