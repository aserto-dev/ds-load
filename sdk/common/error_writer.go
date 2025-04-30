package common

import (
	"io"
)

type errorOptions struct {
	SetExitCode bool
}

type ErrorOption func(*errorOptions)

func WithExitCode(eo *errorOptions) {
	eo.SetExitCode = true
}

type ErrorWriter struct {
	io.Writer
}

func NewErrorWriter(f io.Writer) ErrorWriter {
	return ErrorWriter{f}
}

func (e ErrorWriter) Error(err error, opts ...ErrorOption) {
	if err == nil {
		return
	}

	_, _ = e.Write([]byte(err.Error()))

	options := errorOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.SetExitCode {
		SetExitCode(1)
	}
}
