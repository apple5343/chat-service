package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LoggerKey    string = "logger"
	RequestIDKey string = "request_id"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type ZapLogger struct {
	serviceName string
	logger      *zap.Logger
}

func NewZapLogger(serviceName string) Logger {
	logger, _ := zap.NewProduction()
	logger = logger.WithOptions(zap.AddStacktrace(zapcore.FatalLevel))
	return &ZapLogger{
		serviceName: serviceName,
		logger:      logger,
	}
}

func (l *ZapLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("service", l.serviceName))
	fields = append(fields, zap.String("time", time.Now().Format("2006-01-02 15:04:05")))
	if ctx.Value(RequestIDKey) != nil {
		fields = append(fields, zap.String("request_id", ctx.Value(RequestIDKey).(string)))
	}
	l.logger.Info(msg, fields...)
}

func (l *ZapLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("service", l.serviceName))
	fields = append(fields, zap.String("time", time.Now().Format("2006-01-02 15:04:05")))
	if ctx.Value(RequestIDKey) != nil {
		fields = append(fields, zap.String("request_id", ctx.Value(RequestIDKey).(string)))
	}
	l.logger.Error(msg, fields...)
}

func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

func GetLogger(ctx context.Context) Logger {
	return ctx.Value(LoggerKey).(*ZapLogger)
}
