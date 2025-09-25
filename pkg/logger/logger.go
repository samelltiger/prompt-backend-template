package logger

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"llmapisrv/config"
)

const traceIDKey = "trace_id"

var log *zap.Logger

// Setup 初始化日志
func Setup(cfg config.Logger) {
	// 设置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 设置日志输出
	hook := lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 同时输出到控制台和文件
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(&hook),
			level,
		),
		zapcore.NewCore(
			// zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		),
	)

	// 添加调用者信息
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// Debug 调试级别日志
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Info 信息级别日志
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn 警告级别日志
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Error 错误级别日志
func Error(msg string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	log.Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	log.Fatal(msg, fields...)
}

// Field creates a zap.Field.  This is a helper to improve code readability.
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func getTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(traceIDKey).(string)
	if !ok || traceID == "" { // Handle missing or empty trace ID
		return "no-trace-id" // Or generate a new one if preferred
	}
	return traceID
}

// DebugWithCtx adds a trace_id field from the context to the log message.
func DebugWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, Field(traceIDKey, getTraceID(ctx)))
	log.Debug(msg, fields...)
}

// InfoWithCtx adds a trace_id field from the context to the log message.
func InfoWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, Field(traceIDKey, getTraceID(ctx)))
	log.Info(msg, fields...)
}

// WarnWithCtx adds a trace_id field from the context to the log message.
func WarnWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, Field(traceIDKey, getTraceID(ctx)))
	log.Warn(msg, fields...)
}

// ErrorWithCtx adds a trace_id field from the context to the log message.
func ErrorWithCtx(ctx context.Context, msg string, err error, fields ...zap.Field) {
	fields = append(fields, Field(traceIDKey, getTraceID(ctx)))
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	log.Error(msg, fields...)
}

// FatalWithCtx adds a trace_id field from the context to the log message.
func FatalWithCtx(ctx context.Context, msg string, err error, fields ...zap.Field) {
	fields = append(fields, Field(traceIDKey, getTraceID(ctx)))
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	log.Fatal(msg, fields...)
}

// InfofWithCtx logs a formatted info message with context trace ID
func InfofWithCtx(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fields := []zap.Field{Field(traceIDKey, getTraceID(ctx))}
	log.Info(msg, fields...)
}

// ErrorfWithCtx logs a formatted error message with context trace ID
func ErrorfWithCtx(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fields := []zap.Field{Field(traceIDKey, getTraceID(ctx))}
	log.Error(msg, fields...)
}

// Infof logs a formatted info message with context trace ID
func Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Info(msg)
}

// Errorf logs a formatted error message with context trace ID
func Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Error(msg)
}

func GetLogger() *zap.Logger {
	return log
}
