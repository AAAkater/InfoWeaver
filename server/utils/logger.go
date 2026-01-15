package utils

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func InitLogger(isDev bool) *zap.SugaredLogger {
	var config zap.Config
	if isDev {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}
	l, err := config.Build()
	if err != nil {
		panic(fmt.Errorf("fatal error init logger: %w", err))
	}
	return l.Sugar()
}
