package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey string

const (
	loggerRequestIDKey ctxKey = "x-request-id"
	loggerTraceIDKey   ctxKey = "x-trace-id"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Sync()
}

type L struct {
	z *zap.Logger
}

func NewLogger(env string) Logger {
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if env == "dev" {
		loggerCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, err := loggerCfg.Build()
	if err != nil {
		panic("failed to build logger: " + err.Error())
	}

	lo := L{z: logger}

	return &lo
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, loggerRequestIDKey, requestID)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, loggerTraceIDKey, traceID)
}

func (l *L) Info(ctx context.Context, msg string, fields ...zap.Field) {
	requestID, _ := ctx.Value(loggerRequestIDKey).(string)

	traceID, _ := ctx.Value(loggerTraceIDKey).(string)

	fields = append(
		fields,
		zap.String(string(loggerRequestIDKey), requestID),
		zap.String(string(loggerTraceIDKey), traceID),
	)

	l.z.Info(msg, fields...)
}

func (l *L) Error(ctx context.Context, msg string, fields ...zap.Field) {
	requestID, _ := ctx.Value(loggerRequestIDKey).(string)

	traceID, _ := ctx.Value(loggerTraceIDKey).(string)

	fields = append(
		fields,
		zap.String(string(loggerRequestIDKey), requestID),
		zap.String(string(loggerTraceIDKey), traceID),
	)

	l.z.Error(msg, fields...)
}

func (l *L) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	requestID, _ := ctx.Value(loggerRequestIDKey).(string)

	traceID, _ := ctx.Value(loggerTraceIDKey).(string)

	fields = append(
		fields,
		zap.String(string(loggerRequestIDKey), requestID),
		zap.String(string(loggerTraceIDKey), traceID),
	)

	l.z.Debug(msg, fields...)
}

func (l *L) Sync() {
	_ = l.z.Sync()
}
