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
	infoMsg := &msg.Info{
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
				Name:        "domain",
				Type:        msg.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
				Description: "auth0 domain name",
				Usage:       "--domain=yourdomain.auth0.com",
				Optional:    false,
			},
			{
				Name:        "client-id",
				Type:        msg.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
				Description: "auth0 client id",
				Usage:       "--client-id=yourclientid",
				Optional:    false,
			},
			{
				Name:        "client-secret",
				Type:        msg.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
				Description: "auth0 client secret",
				Usage:       "--client-secret=yourclientsecret",
				Optional:    false,
			},
		},
	}

	message, err := protojson.Marshal(infoMsg)
	if err != nil {
		return err
	}
	os.Stdout.Write(message)
	return nil
}
