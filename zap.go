package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Logger = new(zapL)

type zapL struct {
	serviceLogger    *zap.Logger
	serviceErrLogger *zap.Logger
}

func (ds *zapL) Debug(ctx context.Context, message string, f ...Field) {
	ds.getLogger(ctx).With(f...).Debug(message)
}

func (ds *zapL) Info(ctx context.Context, message string, f ...Field) {
	ds.getLogger(ctx).With(f...).Info(message)
}

func (ds *zapL) Warn(ctx context.Context, message string, f ...Field) {
	ds.getLogger(ctx).With(f...).Warn(message)
}

func (ds *zapL) Error(ctx context.Context, message string, f ...Field) {
	ds.getErrorLogger(ctx).With(f...).Error(message)
}

func (ds *zapL) Panic(ctx context.Context, message string, f ...Field) {
	ds.getErrorLogger(ctx).With(f...).Panic(message)
}

func (ds *zapL) ZapLogger(ctx context.Context) *zap.Logger {
	inFields := make([]zapcore.Field, 0)
	inFields = append(inFields, zap.String("statusCode", "OK"))
	return ds.serviceErrLogger.With(inFields...)
}

func (ds *zapL) getLogger(ctx context.Context) *zap.Logger {
	inFields := make([]zapcore.Field, 0)
	inFields = append(inFields, zap.String("statusCode", "OK"))
	return ds.serviceLogger.With(inFields...)
}

func (ds *zapL) getErrorLogger(ctx context.Context) *zap.Logger {
	inFields := make([]zapcore.Field, 0)
	inFields = append(inFields, zap.String("statusCode", "ERROR"))
	return ds.serviceErrLogger.With(inFields...)
}
