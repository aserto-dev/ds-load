package version

import (
	"fmt"
	"runtime"
	"time"
)

// values set by linker using ldflag -X.
//
//nolint:gochecknoglobals // set by linker
var (
	ver    string
	date   string
	commit string
)

// Info - version info.
type Info struct {
	Version string
	Date    string
	Commit  string
	OS      string
	Arch    string
}

// GetInfo gets version stamp information.
func GetInfo() *Info {
	if ver == "" {
		ver = "0.0.0"
	}

	if date == "" {
		date = time.Now().Format(time.RFC3339)
	}

	if commit == "" {
		commit = "undefined"
	}

	return &Info{
		Version: ver,
		Date:    date,
		Commit:  commit,
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
	}
}

// String() returns the version info string.
func (vi *Info) String() string {
	return fmt.Sprintf("%s g%s %s-%s [%s]",
		vi.Version,
		vi.Commit,
		runtime.GOOS,
		runtime.GOARCH,
		vi.Date,
	)
}

func (vi *Info) AsMap() map[string]any {
	return map[string]any{
		"version": vi.Version,
		"date":    vi.Date,
		"commit":  vi.Commit,
		"os":      vi.OS,
		"arch":    vi.Arch,
	}
}
