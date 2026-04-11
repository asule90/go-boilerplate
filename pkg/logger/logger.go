package logger

import (
	"time"

	"github.com/sule/go-boilerplate/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger creates a zap.Logger based on APP_ENV configuration.
func InitLogger(cfg *config.Provider) (*zap.Logger, error) {
	loc, err := time.LoadLocation(cfg.App.Timezone)
	if err != nil {
		loc = time.UTC
	}

	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.In(loc).Format(time.RFC3339))
	}

	var zapCfg zap.Config
	if cfg.App.Env == "production" {
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		zapCfg.DisableStacktrace = true
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	zapCfg.EncoderConfig.EncodeTime = timeEncoder

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
