package app

import (
	"context"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/common/msg"
	"github.com/aserto-dev/ds-load/plugins/okta/pkg/oktaclient"
	"github.com/aserto-dev/ds-load/plugins/sdk/transform"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(ctx *kong.Context) error {
	oktaClient, err := oktaclient.NewOktaClient(context.Background(), cmd.Domain, cmd.APIToken)
	if err != nil {
		return err
	}

	results := make(chan map[string]interface{}, 1)
	errCh := make(chan error, 1)
	go func() {
		cmd.Fetch(oktaClient, results, errCh)
		close(results)
		close(errCh)
	}()

	go func() {
		for err := range errCh {
			os.Stderr.WriteString(err.Error())
			os.Stderr.WriteString("\n")
		}
	}()

	fileReader := os.ReadFile
	template := cmd.TemplateFile
	if template == "" {
		fileReader = Assets().ReadFile
		template = "assets/transform_template.tmpl"
	}
	templateContent, err := fileReader(template)
	if err != nil {
		return err
	}

	tranformer := transform.NewTransformer(cmd.MaxChunkSize)
	for input := range results {
		output, err := tranformer.TransformToTemplate(input, string(templateContent))
		if err != nil {
			return errors.Wrap(err, "transform template execute failed")
		}
		var directoryObject msg.Transform

		err = protojson.Unmarshal([]byte(output), &directoryObject)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal transformed data into directory objects and relations")
		}

		objectChunks, relationChunks := tranformer.PrepareChunks(&directoryObject)

		err = tranformer.WriteChunks(os.Stdout, objectChunks, relationChunks)
		if err != nil {
			return errors.Wrap(err, "failed to write chunks to output")
		}
	}

	return nil
}
