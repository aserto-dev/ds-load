package common

import (
	"os"
)

type ErrorWriter struct {
	os.File
}

type ErrorOptions struct {
	SetExitCode bool
}

type ErrorOption func(*ErrorOptions)

func NewErrorWriter(f *os.File) ErrorWriter {
	if f == nil {
		return ErrorWriter{*os.Stderr}
	}

	return ErrorWriter{*f}
}

func (e *ErrorWriter) Error(err error, opts ...ErrorOption) {
	if err == nil {
		return
	}

	_, _ = e.Write([]byte(err.Error()))

	options := ErrorOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.SetExitCode {
		SetExitCode(1)
	}
}

func WithExitCode(eo *ErrorOptions) {
	eo.SetExitCode = true
}
