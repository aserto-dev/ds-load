package app

import (
	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/ldapclient"
	"github.com/aserto-dev/ds-load/sdk/common/cc"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	credentials := &ldapclient.Credentials{
		User:     v.User,
		Password: v.Password,
	}

	conOptions := &ldapclient.ConnectionOptions{
		Host:        v.Host,
		BaseDN:      v.BaseDn,
		UserFilter:  v.UserFilter,
		GroupFilter: v.GroupFilter,
		Insecure:    v.Insecure,
		IDField:     v.IDField,
	}

	client, err := ldapclient.NewLDAPClient(credentials, conOptions, ctx.Log)
	if err != nil {
		return err
	}
	defer client.Close()

	return nil
}
