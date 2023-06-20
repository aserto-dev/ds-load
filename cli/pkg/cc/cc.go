package cc

import (
	"context"
	"log"
	"os"

	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/ds-load/cli/pkg/cc/iostream"
	logger "github.com/aserto-dev/logger"
	"github.com/rs/zerolog"
)

type CommonCtx struct {
	Context    context.Context
	UI         *clui.UI
	LogLevel   zerolog.Level
	Log        *zerolog.Logger
	ConfigPath string
}

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
		UI:         iostream.NewUI(iostream.DefaultIO()),
		LogLevel:   logLevelParsed,
		Log:        newLogger,
		ConfigPath: config,
	}
}

func GetLogLevel(intLevel int) zerolog.Level {
	logLevel := zerolog.FatalLevel

	switch intLevel {
	case 1:
		logLevel = zerolog.ErrorLevel
	case 2:
		logLevel = zerolog.InfoLevel
	case 3:
		logLevel = zerolog.DebugLevel
	case 4:
		logLevel = zerolog.TraceLevel
	}

	return logLevel
}
