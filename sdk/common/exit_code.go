package common

import "sync/atomic"

var (
	exitCode int32
)

func GetExitCode() int {
	return int(atomic.LoadInt32(&exitCode))
}

func SetExitCode(code int) {
	atomic.StoreInt32(&exitCode, int32(code)) //nolint:gosec
}
