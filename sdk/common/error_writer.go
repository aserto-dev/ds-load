package common

import "os"

type ErrorWriter struct {
	Writer *os.File
}

func (e *ErrorWriter) ErrorNoExitCode(err error) {
	if err == nil {
		return
	}

	_, _ = e.Writer.Write([]byte(err.Error()))
}

func (e *ErrorWriter) Error(err error) {
	if err == nil {
		return
	}

	_, _ = e.Writer.Write([]byte(err.Error()))

	SetExitCode(1)
}
