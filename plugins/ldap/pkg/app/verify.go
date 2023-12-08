package app

import (
	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/ldapclient"
)

type VerifyCmd struct {
	FetchCmd
}

func (v *VerifyCmd) Run(ctx *cc.CommonCtx) error {
	_, err := ldapclient.NewLDAPClient(v.User, v.Password, v.Host, v.BaseDn, v.UserFilter, v.GroupFilter)
	if err != nil {
		return err
	}

	return nil
}
