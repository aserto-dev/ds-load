package cc

import (
	"context"

	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/ds-load/cli/pkg/cc/iostream"
)

type CommonCtx struct {
	Context   context.Context
	UI        *clui.UI
	Verbosity int
}

func NewCommonContext(verbosity int) *CommonCtx {
	return &CommonCtx{
		Context:   context.Background(),
		UI:        iostream.NewUI(iostream.DefaultIO()),
		Verbosity: verbosity,
	}
}
