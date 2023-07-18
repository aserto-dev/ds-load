package plugin

import (
	"context"
	"io"
	"log"
	"os"
	"sync"

	"github.com/aserto-dev/ds-load/sdk/common/js"
)

type Fetcher interface {
	Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error
}

type Transformer interface {
	Transform(ctx context.Context, reader io.Reader, outputWriter, errorWriter io.Writer) error
}

type Publisher interface {
	Publish(ctx context.Context, reader io.Reader) error
}

type Plugin interface {
	Fetcher
	Transformer
}

type DSPlugin struct {
	template  []byte
	outWriter io.Writer
	errWriter io.Writer
}

func NewDSPlugin(options ...PluginOption) *DSPlugin {
	plugin := &DSPlugin{}
	for _, o := range options {
		o(plugin)
	}

	if plugin.outWriter == nil {
		plugin.outWriter = os.Stdout
	}

	if plugin.errWriter == nil {
		plugin.errWriter = os.Stderr
	}

	return plugin
}

// json encodes results and prints to plugin writer.
func (plugin *DSPlugin) WriteFetchOutput(results chan map[string]interface{}, errCh chan error, transformMessage bool) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for err := range errCh {
			_, wErr := plugin.errWriter.Write([]byte(err.Error() + "\n"))
			if wErr != nil {
				log.Fatalf("cannot write to output: %s", wErr.Error())
			}
		}
		wg.Done()
	}()

	writer, err := js.NewJSONArrayWriter(plugin.outWriter)
	if err != nil {
		return err
	}
	defer writer.Close()
	for result := range results {
		err := writer.Write(result)
		if err != nil {
			return err
		}
	}
	wg.Wait()
	return nil
}

func (plugin *DSPlugin) ExportTransform() error {
	_, err := plugin.outWriter.Write(plugin.template)
	if err != nil {
		log.Fatalf("cannot write to output: %s", err.Error())
	}

	return nil
}
