package plugin

import "io"

type PluginOption func(*DSPlugin)

func WithTemplate(templateContent []byte) PluginOption {
	return func(d *DSPlugin) {
		d.template = templateContent
	}
}

func WithOutputWriter(writer io.Writer) PluginOption {
	return func(d *DSPlugin) {
		d.outWriter = writer
	}
}

func WithErrorWriter(writer io.Writer) PluginOption {
	return func(d *DSPlugin) {
		d.errWriter = writer
	}
}
