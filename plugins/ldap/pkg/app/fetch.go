package app

import (
	"os"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/fetch"
	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/ldapclient"
)

type FetchCmd struct {
	User        string `short:"u" help:"LDAP user" env:"LDAP_USER" required:""`
	Password    string `short:"p" help:"LDAP password" env:"LDAP_PASSWORD" required:""`
	Host        string `short:"s" help:"LDAP host" env:"LDAP_HOST" required:""`
	BaseDn      string `short:"b" help:"LDAP base DN" env:"LDAP_BASE_DN" default:"dc=example,dc=org"`
	UserFilter  string `short:"f" help:"LDAP user filter" env:"LDAP_USER_FILTER" default:"(&(objectClass=organizationalPerson))"`
	GroupFilter string `short:"g" help:"LDAP group filter" env:"LDAP_GROUP_FILTER" default:"(&(objectClass=groupOfNames))"`
}

func (cmd *FetchCmd) Run(ctx *cc.CommonCtx) error {
	ldapClient, err := ldapclient.NewLDAPClient(cmd.User, cmd.Password, cmd.Host, cmd.BaseDn, cmd.UserFilter, cmd.GroupFilter)
	defer ldapClient.Close()

	fetcher, err := fetch.New(ldapClient)
	if err != nil {
		return err
	}

	return fetcher.Fetch(ctx.Context, os.Stdout, os.Stderr)
}
