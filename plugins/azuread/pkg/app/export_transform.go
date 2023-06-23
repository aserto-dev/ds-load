package app

import (
	"os"

	"github.com/alecthomas/kong"
)

type ExportTransportCmd struct {
}

func (t *ExportTransportCmd) Run(context *kong.Context) error {
	content, err := Assets().ReadFile("assets/transform_template.tmpl")
	if err != nil {
		return err
	}
	os.Stdout.WriteString(string(content))
	return nil
}
