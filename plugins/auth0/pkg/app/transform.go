// get json from stdin and return json with v2 objects and v2 relations
package app

import (
	"context"
	"os"

	"github.com/aserto-dev/ds-load/sdk/common/cc"
	"github.com/aserto-dev/ds-load/sdk/plugin"
	"github.com/aserto-dev/ds-load/sdk/template"
	"github.com/aserto-dev/ds-load/sdk/transform"
)

type TransformCmd struct {
	Template string `name:"template" short:"t" env:"DS_TEMPLATE_FILE" help:"transformation template file path" type:"path" optional:""`
}

func (t *TransformCmd) Run(ctx *cc.CommonCtx) error {
	templateContent, err := t.getTemplateContent()
	if err != nil {
		return err
	}

	goTemplateTransformer := transform.NewGoTemplateTransform(templateContent)

	return t.transform(ctx.Context, goTemplateTransformer)
}

func (t *TransformCmd) transform(ctx context.Context, transformer plugin.Transformer) error {
	return transformer.Transform(ctx, os.Stdin, os.Stdout, os.Stderr)
}

func (t *TransformCmd) getTemplateContent() ([]byte, error) {
	templateContent, err := Assets().ReadFile("assets/transform_template.tmpl")
	if err != nil {
		return nil, err
	}

	templateLoader := template.NewTemplateLoader(templateContent)

	templateContent, err = templateLoader.Load(t.Template)
	if err != nil {
		return nil, err
	}

	return templateContent, nil
}
