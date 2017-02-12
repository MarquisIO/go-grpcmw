package grpcmw

import (
	"google.golang.org/grpc"
)

type StreamServerInterceptor interface {
	StreamInterceptor() grpc.StreamServerInterceptor
	AddGRPCStreamInterceptor(i ...grpc.StreamServerInterceptor) StreamServerInterceptor
	AddStreamInterceptor(i ...StreamServerInterceptor) StreamServerInterceptor
}

type streamServerInterceptor struct {
	interceptors []grpc.StreamServerInterceptor
}

func chainStream(current grpc.StreamServerInterceptor, info *grpc.StreamServerInfo, next grpc.StreamHandler) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return current(srv, stream, info, next)
	}
}

func NewStreamServerInterceptor(arr ...grpc.StreamServerInterceptor) StreamServerInterceptor {
	return &streamServerInterceptor{
		interceptors: arr,
	}
}

func (si *streamServerInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// TODO: Find a more efficient way
		interceptor := handler
		for idx := len(si.interceptors) - 1; idx >= 0; idx-- {
			interceptor = chainStream(si.interceptors[idx], info, interceptor)
		}
		return interceptor(srv, ss)
	}
}

func (si *streamServerInterceptor) AddGRPCStreamInterceptor(arr ...grpc.StreamServerInterceptor) StreamServerInterceptor {
	si.interceptors = append(si.interceptors, arr...)
	return si
}

func (si *streamServerInterceptor) AddStreamInterceptor(arr ...StreamServerInterceptor) StreamServerInterceptor {
	for _, i := range arr {
		si.interceptors = append(si.interceptors, i.StreamInterceptor())
	}
	return si
}
