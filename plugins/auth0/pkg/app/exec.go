package app

import (
	"net/http"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/common/js"
	"github.com/aserto-dev/ds-load/common/msg"
	"github.com/aserto-dev/ds-load/plugins/auth0/pkg/httpclient"
	"github.com/aserto-dev/ds-load/plugins/sdk/transform"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/auth0.v5/management"
)

type ExecCmd struct {
	FetchCmd
	TransformCmd
}

func (cmd *ExecCmd) Run(context *kong.Context) error {
	options := []management.ManagementOption{
		management.WithClientCredentials(
			cmd.ClientID,
			cmd.ClientSecret,
		),
	}
	if cmd.RateLimit {
		client := http.DefaultClient
		client.Transport = httpclient.NewTransport(http.DefaultTransport)
		options = append(options, management.WithClient(client))
	}

	mgmt, err := management.New(
		cmd.Domain,
		options...,
	)
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
	jsonWriter := js.NewJSONArrayWriter(os.Stdout)

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

		chunks := tranformer.PrepareChunks(&directoryObject)
		err = tranformer.WriteChunks(jsonWriter, chunks)
		if err != nil {
			return errors.Wrap(err, "failed to write chunks to output")
		}
	}
	jsonWriter.Close()

	return nil
}
