package app

import (
	"os"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/fetch"
	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/ldapclient"
)

type FetchCmd struct {
	User     string `short:"u" help:"LDAP user" env:"LDAP_USER" required:""`
	Password string `short:"p" help:"LDAP password" env:"LDAP_PASSWORD" required:""`
	Host     string `short:"s" help:"LDAP host" env:"LDAP_HOST" required:""`
}

func (cmd *FetchCmd) Run(ctx *cc.CommonCtx) error {
	ldapClient, err := ldapclient.NewLDAPClient(cmd.User, cmd.Password, cmd.Host)
	defer ldapClient.Close()

	fetcher, err := fetch.New(ldapClient)
	if err != nil {
		return err
	}

	return fetcher.Fetch(ctx.Context, os.Stdout, os.Stderr)
}
