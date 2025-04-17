package cc

import (
	"context"
	"log"
	"os"

	logger "github.com/aserto-dev/logger"
	"github.com/rs/zerolog"
)

type CommonCtx struct {
	Context    context.Context
	LogLevel   zerolog.Level
	Log        *zerolog.Logger
	ConfigPath string
}

const (
	errLogLevel   = 1
	infoLogLevel  = 2
	debugLogLevel = 3
	traceLogLevel = 4
)

func NewCommonContext(verbosity int, config string) *CommonCtx {
	logLevelParsed := GetLogLevel(verbosity)
	logCfg := &logger.Config{
		Prod:           false,
		LogLevelParsed: logLevelParsed,
	}

	newLogger, err := logger.NewLogger(os.Stdout, os.Stderr, logCfg)
	if err != nil {
		log.Fatalf("failed to initialize logger: %s", err.Error())
	}

	return &CommonCtx{
		Context:    context.Background(),
		LogLevel:   logLevelParsed,
		Log:        newLogger,
		ConfigPath: config,
	}
}

func GetLogLevel(intLevel int) zerolog.Level {
	logLevel := zerolog.FatalLevel

	switch intLevel {
	case errLogLevel:
		logLevel = zerolog.ErrorLevel
	case infoLogLevel:
		logLevel = zerolog.InfoLevel
	case debugLogLevel:
		logLevel = zerolog.DebugLevel
	case traceLogLevel:
		logLevel = zerolog.TraceLevel
	}

	return logLevel
}
