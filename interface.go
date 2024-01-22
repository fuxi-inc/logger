package logger

import (
	"context"

	"go.uber.org/zap"
)

type Logger interface {
	Debug(ctx context.Context, message string, fields ...Field)
	Info(ctx context.Context, message string, args ...Field)
	Warn(ctx context.Context, message string, args ...Field)
	Error(ctx context.Context, message string, args ...Field)
	Panic(ctx context.Context, message string, args ...Field)
	ZapLogger(ctx context.Context) *zap.Logger // 这个是为了兼容dis-service加上的，后续有机会会下掉，别再加新的使用方式了。。
}
