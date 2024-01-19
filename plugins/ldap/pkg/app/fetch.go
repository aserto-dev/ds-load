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
	BaseDn      string `short:"b" help:"LDAP base DN" env:"LDAP_BASE_DN" required:""`
	UserFilter  string `short:"f" help:"LDAP user filter" env:"LDAP_USER_FILTER" default:"(&(objectClass=organizationalPerson))"`
	GroupFilter string `short:"g" help:"LDAP group filter" env:"LDAP_GROUP_FILTER" default:"(&(objectClass=groupOfNames))"`
	Insecure    bool   `short:"i" help:"Allow insecure LDAP connection" env:"LDAP_INSECURE" default:"false"`
	UUIDField   string `short:"U" help:"LDAP field to use as UUID" env:"LDAP_UUID_FIELD" default:"objectGUID"`
}

func (cmd *FetchCmd) Run(ctx *cc.CommonCtx) error {
	credentials := &ldapclient.Credentials{
		User:     cmd.User,
		Password: cmd.Password,
	}

	conOptions := &ldapclient.ConnectionOptions{
		Host:        cmd.Host,
		BaseDN:      cmd.BaseDn,
		UserFilter:  cmd.UserFilter,
		GroupFilter: cmd.GroupFilter,
		Insecure:    cmd.Insecure,
		UUIDField:   cmd.UUIDField,
	}

	ldapClient, err := ldapclient.NewLDAPClient(credentials, conOptions, ctx.Log)
	if err != nil {
		return err
	}
	defer ldapClient.Close()

	fetcher, err := fetch.New(ldapClient)
	if err != nil {
		return err
	}

	return fetcher.Fetch(ctx.Context, os.Stdout, os.Stderr)
}
