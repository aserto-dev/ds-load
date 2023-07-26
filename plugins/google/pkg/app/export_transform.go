package app

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/sdk/transform"
)

type ExportTransportCmd struct {
}

func (t *ExportTransportCmd) Run(context *kong.Context) error {
	templateContent, err := Assets().ReadFile("assets/transform_template.tmpl")
	if err != nil {
		return err
	}
	transformer := transform.NewGoTemplateTransform(templateContent)

	return transformer.ExportTransform(os.Stdout)
}
