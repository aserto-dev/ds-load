// get json from stdin and return json with v2 objects and v2 relations
package app

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/sdk/plugin"
)

type TransformCmd struct {
	Template     string `name:"template" short:"t" env:"DS_TEMPLATE_FILE" help:"transformation template file path" type:"path" optional:""`
	MaxChunkSize int    `name:"max-chunk-size" env:"DS_MAX_CHUNK_SIZE" help:"maximum chunk size" default:"20" optional:""`
}

func (t *TransformCmd) Run(context *kong.Context) error {
	content, err := t.getTemplateContent()
	if err != nil {
		return err
	}
	return plugin.NewDSPlugin(plugin.WithTemplate(content), plugin.WithMaxChunkSize(t.MaxChunkSize)).Transform()
}

func (t *TransformCmd) getTemplateContent() ([]byte, error) {
	var templateContent []byte
	var err error
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
