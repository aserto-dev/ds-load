package app

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/sdk/plugin"
)

type ExportTransportCmd struct {
}

func (t *ExportTransportCmd) Run(context *kong.Context) error {
	templateContent, err := Assets().ReadFile("assets/transform_template.tmpl")
	if err != nil {
		return err
	}
	return plugin.NewDSPlugin(plugin.WithTemplate(templateContent)).ExportTransform()
}
