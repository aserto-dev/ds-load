package app

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/common/msg"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/auth0.v5/management"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(context *kong.Context) error {
	mgmt, err := management.New(
		cmd.Domain,
		management.WithClientCredentials(
			cmd.ClientID,
			cmd.ClientSecret,
		))
	if err != nil {
		return err
	}

	results := make(chan map[string]interface{}, 1)
	errCh := make(chan error, 1)
	go func() {
		Fetch(mgmt, results, errCh)
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

	if cmd.MaxChunkSize == 0 {
		cmd.MaxChunkSize = 1 // By default do not write begin and end batches
	}

	for input := range results {
		output, err := Transform(input, string(templateContent))
		if err != nil {
			return errors.Wrap(err, "transform template execute failed")
		}
		var directoryObject msg.Transform

		err = protojson.Unmarshal([]byte(output), &directoryObject)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal transformed data into directory objects and relations")
		}

		objectChunks, relationChunks := cmd.prepareChunks(&directoryObject)

		err = cmd.writeResponse(os.Stdout, objectChunks, relationChunks)
		if err != nil {
			return errors.Wrap(err, "failed to write chunks to output")
		}
	}

	return nil
}
