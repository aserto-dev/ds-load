package app

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/sdk/common/version"
)

type CLI struct {
	Config          kong.ConfigFlag    `help:"Configuration file path" short:"c"`
	Version         VersionCmd         `cmd:"" help:"version information"`
	Fetch           FetchCmd           `cmd:"" help:"fetch hubspot data"`
	Transform       TransformCmd       `cmd:"" help:"transform hubspot data"`
	ExportTransform ExportTransportCmd `cmd:"" help:"export default transform template"`
	Exec            ExecCmd            `cmd:"" help:"fetch and transform hubspot data" default:"withargs"`
	GetRefreshToken GetTokenCmd        `cmd:"" help:"obtain a refresh token from hubspot"`
}

type VersionCmd struct {
}

func (cmd *VersionCmd) Run() error {
	fmt.Printf("%s - %s\n",
		AppName,
		version.GetInfo().String(),
	)
	return nil
}
