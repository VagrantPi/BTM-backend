package logger

import (
	"context"
	"os"

	"github.com/josestg/lazy"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _zapLazy = lazy.New(func() (*Logger, error) {
	AtomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	zapOpts := make([]zap.Option, 0)
	zapOpts = append(
		zapOpts, zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel),
		),
		zap.AddCaller(),
	)
	if isDebug() {
		AtomicLevel.SetLevel(zap.DebugLevel)
		zapOpts = append(zapOpts, zap.Development())
	}

	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder.EncodeLevel = zapcore.CapitalLevelEncoder

	coreInfo := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		AtomicLevel,
	)

	core := zapcore.NewTee(coreInfo)
	zapLogger := zap.New(core, zapOpts...)
	return &Logger{zapLogger}, nil
})

func Zap() *Logger {
	return _zapLazy.Value().clone()
}

type Logger struct {
	logger *zap.Logger
}

func WithContext(ctx context.Context) *Logger {
	return Zap().WithContext(ctx)
}

func (log *Logger) WithContext(ctx context.Context) *Logger {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return log
	}
	l := log.With(
		zap.String("trace_id", span.SpanContext().TraceID().String()),
		zap.String("span_id", span.SpanContext().SpanID().String()),
	)
	return l
}

func (log *Logger) WithClassFunction(className string, functionName string) *Logger {
	l := log.With(
		zap.String("class", className),
		zap.String("functionName", functionName),
	)
	return l
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (log *Logger) With(fields ...zap.Field) *Logger {
	if len(fields) == 0 {
		return log
	}
	l := log.clone()
	l.logger = l.logger.With(fields...)
	return l
}

func (log *Logger) clone() *Logger {
	l := *log
	logger := *l.logger
	l.logger = &logger
	return &l
}

func (log *Logger) WithFunction(functionName string) *Logger {
	fields := []zap.Field{zap.String("function", functionName)}
	log.logger = log.logger.With(
		fields...,
	)
	return log
}

func (log *Logger) Sync() error {
	return log.logger.Sync()
}

func (log *Logger) Error(msg string, fields ...zap.Field) {
	log.logger.Error(msg, fields...)
}

func (log *Logger) Warn(msg string, fields ...zap.Field) {
	log.logger.Warn(msg, fields...)
}

func (log *Logger) Info(msg string, fields ...zap.Field) {
	log.logger.Info(msg, fields...)
}

func (log *Logger) Debug(msg string, fields ...zap.Field) {
	log.logger.Debug(msg, fields...)
}

func (log *Logger) GetZapLogger() *zap.Logger {
	return log.logger
}

func isDebug() bool {
	return len(os.Getenv("DEBUG")) != 0 && os.Getenv("DEBUG") == "true"
}
