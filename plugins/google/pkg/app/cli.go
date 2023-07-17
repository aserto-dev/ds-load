package app

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/sdk/common/version"
)

type CLI struct {
	Config          kong.ConfigFlag    `help:"Configuration file path" short:"c"`
	Version         VersionCmd         `cmd:"" help:"version information"`
	Fetch           FetchCmd           `cmd:"" help:"fetch google workspace data"`
	Transform       TransformCmd       `cmd:"" help:"transform google workspace data"`
	ExportTransform ExportTransportCmd `cmd:"" help:"export default transform template"`
	Exec            ExecCmd            `cmd:"" help:"fetch and transform google workspace data" default:"withargs"`
	GetRefreshToken GetTokenCmd        `cmd:"" help:"obtain a refresh token from GCP"`
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
