package logging

import (
	"os"
	"strings"

	"notification-service/services/template-svc/internal/platform/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(cfg *config.Config) (*zap.Logger, error) {
	level := parseLevel(cfg.LogLevel)

	var encoder zapcore.Encoder
	if cfg.Env == "local" || cfg.Env == "development" {
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	consoleWriter := zapcore.AddSync(os.Stdout)

	logFilePath := strings.TrimRight(cfg.LogRootPath, "/") + "/app.log"
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    cfg.LogMaxSizeMB,
		MaxBackups: cfg.LogMaxBackups,
		MaxAge:     cfg.LogMaxAgeDays,
		Compress:   true,
	})

	consoleCore := zapcore.NewCore(encoder, consoleWriter, level)
	fileCore := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), fileWriter, level)

	core := zapcore.NewTee(consoleCore, fileCore)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
