package exec

import (
	"context"
	"io"
	"os"

	"github.com/aserto-dev/ds-load/sdk/plugin"
	"github.com/rs/zerolog"
)

func Execute(ctx context.Context, log *zerolog.Logger, transformer plugin.Transformer, fetcher plugin.Fetcher) error {
	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()

	go func() {
		err := fetcher.Fetch(ctx, pipeWriter, os.Stderr)
		if err != nil {
			log.Printf("Could not fetch data %s", err.Error())
		}
		pipeWriter.Close()
	}()

	return transformer.Transform(ctx, pipeReader, os.Stdout, os.Stderr)
}
