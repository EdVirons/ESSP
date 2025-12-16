package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(parseLevel(level))
	l, _ := cfg.Build()
	return l
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func Err(err error) zap.Field       { return zap.Error(err) }
func Str(k, v string) zap.Field     { return zap.String(k, v) }
func Int(k string, v int) zap.Field { return zap.Int(k, v) }
