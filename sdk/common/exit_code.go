package common

import (
	"io"
	"sync/atomic"
)

var exitCode int32

func GetExitCode() int {
	return int(atomic.LoadInt32(&exitCode))
}

func SetExitCode(code int) {
	atomic.StoreInt32(&exitCode, int32(code)) //nolint:gosec
}

func WriteErrorWithExitCode(w io.Writer, err error, code int) {
	_, _ = w.Write([]byte(err.Error()))

	SetExitCode(code)
}
