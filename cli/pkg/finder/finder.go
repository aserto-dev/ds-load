package finder

import "github.com/aserto-dev/ds-load/cli/pkg/plugin"

type Finder interface {
	Find() ([]*plugin.Plugin, error)
}
