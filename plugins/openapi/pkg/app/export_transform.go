package app

import (
	"os"

	"github.com/aserto-dev/ds-load/cli/pkg/cc"
	"github.com/aserto-dev/ds-load/sdk/transform"
)

type ExportTransformCmd struct{}

func (t *ExportTransformCmd) Run(ctx *cc.CommonCtx) error {
	templateContent, err := Assets().ReadFile("assets/transform_template.tmpl")
	if err != nil {
		return err
	}
	transformer := transform.NewGoTemplateTransform(templateContent)

	return transformer.ExportTransform(os.Stdout)
}
