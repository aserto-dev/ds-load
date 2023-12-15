package app

import (
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/ldapclient"
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
		UuidField:   v.UuidField,
	}

	_, err := ldapclient.NewLDAPClient(credentials, conOptions)
	if err != nil {
		return err
	}

	return nil
}
