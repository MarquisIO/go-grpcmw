package grpcmw

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type UnaryServerInterceptor interface {
	UnaryInterceptor() grpc.UnaryServerInterceptor
	AddGRPCUnaryInterceptor(i ...grpc.UnaryServerInterceptor) UnaryServerInterceptor
	AddUnaryInterceptor(i ...UnaryServerInterceptor) UnaryServerInterceptor
}

type unaryServerInterceptor struct {
	interceptors []grpc.UnaryServerInterceptor
}

func chainUnary(current grpc.UnaryServerInterceptor, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return current(ctx, req, info, next)
	}
}

func NewUnaryServerInterceptor(arr ...grpc.UnaryServerInterceptor) UnaryServerInterceptor {
	return &unaryServerInterceptor{
		interceptors: arr,
	}
}

func (ui *unaryServerInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// TODO: Find a more efficient way
		interceptor := handler
		for idx := len(ui.interceptors) - 1; idx >= 0; idx-- {
			interceptor = chainUnary(ui.interceptors[idx], info, interceptor)
		}
		return interceptor(ctx, req)
	}
}

func (ui *unaryServerInterceptor) AddGRPCUnaryInterceptor(arr ...grpc.UnaryServerInterceptor) UnaryServerInterceptor {
	ui.interceptors = append(ui.interceptors, arr...)
	return ui
}

func (ui *unaryServerInterceptor) AddUnaryInterceptor(arr ...UnaryServerInterceptor) UnaryServerInterceptor {
	for _, i := range arr {
		ui.interceptors = append(ui.interceptors, i.UnaryInterceptor())
	}
	return ui
}
