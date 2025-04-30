// get json from stdin and return json with v2 objects and v2 relations
package app

import (
	"context"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/sdk/plugin"
	"github.com/aserto-dev/ds-load/sdk/transform"
)

type TransformCmd struct {
	Template string `name:"template" short:"t" env:"DS_TEMPLATE_FILE" help:"transformation template file path" type:"path" optional:""`
}

const (
	defaultTimeout = 1500 * time.Millisecond
)

func (t *TransformCmd) Run(kongContext *kong.Context) error {
	template, err := t.getTemplateContent()
	if err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	goTemplateTransformer := transform.NewGoTemplateTransform(template)

	return t.transform(timeoutCtx, goTemplateTransformer)
}

func (t *TransformCmd) transform(ctx context.Context, transformer plugin.Transformer) error {
	return transformer.Transform(ctx, os.Stdin, os.Stdout, os.Stderr)
}

func (t *TransformCmd) getTemplateContent() ([]byte, error) {
	var (
		templateContent []byte
		err             error
	)

	if t.Template == "" {
		templateContent, err = Assets().ReadFile("assets/transform_template.tmpl")
		if err != nil {
			return nil, err
		}
	} else {
		templateContent, err = os.ReadFile(t.Template)
		if err != nil {
			return nil, err
		}
	}

	return templateContent, nil
}
