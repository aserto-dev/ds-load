package app

import (
	"github.com/alecthomas/kong"
)

type FetchCmd struct {
	Domain    string `cmd:"" env:"DS_OKTA_DOMAIN"`
	APIToken  string `cmd:"" env:"DS_OKTA_TOKEN"`
	UserPID   string `cmd:"" env:"DS_OKTA_USER_PID"`
	UserEmail string `cmd:"" env:"DS_OKTA_USER_EMAIL"`
}

func (cmd *FetchCmd) Run(context *kong.Context) error {
	return nil
}
