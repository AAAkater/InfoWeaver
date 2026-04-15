package utils

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger = InitLogger()

func InitLogger() *zap.SugaredLogger {
	var config zap.Config

	config = zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Customize time format to YYYY-MM-DD HH:MM:SS.mmm
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	l, err := config.Build()
	if err != nil {
		panic(fmt.Errorf("fatal error init logger: %w", err))
	}
	return l.Sugar()
}
