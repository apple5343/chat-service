package grpc

import (
	"chat-service/errorz"
	"chat-service/pkg/logger"
	"context"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RequestsLogger(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx = logger.ContextWithLogger(ctx, l)
		h, err := handler(ctx, req)
		if err != nil {
			h, err = logError(ctx, l, h, info, err)
			return h, err
		}
		return h, nil
	}
}

func logError(ctx context.Context, l logger.Logger, h any, info *grpc.UnaryServerInfo, err error) (any, error) {
	code := codes.Internal
	if customErr := errorz.Parse(err); customErr != nil {
		code = codes.Code(customErr.StatusCode())
		l.Error(ctx, err.Error())
		if !reflect.ValueOf(h).IsNil() {
			return h, nil
		}
		return h, status.Errorf(code, customErr.Message())
	}
	l.Error(ctx, err.Error())
	return h, status.Errorf(code, errorz.ErrUnknown.Error())
}
