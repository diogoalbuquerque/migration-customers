package logger_test

import (
	"github.com/diogoalbuquerque/migration-customers/pkg/logger"
	"github.com/rs/zerolog"
	"testing"
)

func Test_Logger_Print(t *testing.T) {

	level := zerolog.ErrorLevel.String()
	l := logger.New(level)
	l.Error("Mock Log %s", level)
	l.Error("Mock Log %s", level)

	level = zerolog.WarnLevel.String()
	l = logger.New(level)
	l.Warn("Mock Log %s", level)

	level = zerolog.InfoLevel.String()
	l = logger.New(level)
	l.Info("Mock Log %s", level)

	level = zerolog.DebugLevel.String()
	l = logger.New(level)
	l.Debug("Mock Log %s", level)

	l = logger.New("default")
	l.Info("Mock Log default")
}
