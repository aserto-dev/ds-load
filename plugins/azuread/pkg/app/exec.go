package app

import (
	"github.com/alecthomas/kong"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(context *kong.Context) error {
	return nil
}
