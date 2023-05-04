package app

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/common/msg"
	"github.com/aserto-dev/ds-load/common/version"
	"github.com/aserto-dev/go-grpc/aserto/common/info/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type InfoCmd struct{}

func (cmd *InfoCmd) Run(context *kong.Context) error {
	ver := version.GetInfo()
	info := &msg.Info{
		Build: &info.BuildInfo{
			Version: ver.Version,
			Commit:  ver.Commit,
			Date:    ver.Date,
			Os:      ver.OS,
			Arch:    ver.Arch,
		},
		Description: AppDescription,
		Configs: []*msg.ConfigElement{
			{
				Name:        "auth0-domain",
				Type:        msg.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
				Description: "auth0 domain name",
				Usage:       "--auth0-domain=yourdomain.auth0.com",
				Optional:    false,
			},
			{
				Name:        "auth0-client-id",
				Type:        msg.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
				Description: "auth0 client id",
				Usage:       "--auth0-client-id=yourclientid",
				Optional:    false,
			},
			{
				Name:        "auth0-client-secret",
				Type:        msg.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
				Description: "auth0 client secret",
				Usage:       "--auth0-client-secret=yourclientsecret",
				Optional:    false,
			},
			{
				Name:        "auth0-connection-name",
				Type:        msg.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
				Description: "auth0 connection name - defaults to Username-Password-Authentication",
				Usage:       "--auth0-connection-name=yourconnectionname",
				Optional:    true,
			},
		},
	}

	message, err := protojson.Marshal(info)
	if err != nil {
		return err
	}
	os.Stdout.Write(message)
	return nil
}
