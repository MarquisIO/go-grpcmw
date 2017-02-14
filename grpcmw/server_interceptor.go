package grpcmw

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// StreamServerInterceptor represents a server interceptor for gRPC methods that
// return a stream. It allows chaining of `grpc.StreamServerInterceptor`
// and other `StreamServerInterceptor`.
type StreamServerInterceptor interface {
	StreamInterceptor() grpc.StreamServerInterceptor
	AddGRPCStreamInterceptor(i ...grpc.StreamServerInterceptor) StreamServerInterceptor
	AddStreamInterceptor(i ...StreamServerInterceptor) StreamServerInterceptor
}

// UnaryServerInterceptor represents a server interceptor for gRPC methods that
// return a single value instead of a stream. It allows chaining of
// `grpc.UnaryServerInterceptor` and other `UnaryServerInterceptor`.
type UnaryServerInterceptor interface {
	UnaryInterceptor() grpc.UnaryServerInterceptor
	AddGRPCUnaryInterceptor(i ...grpc.UnaryServerInterceptor) UnaryServerInterceptor
	AddUnaryInterceptor(i ...UnaryServerInterceptor) UnaryServerInterceptor
}

type streamServerInterceptor struct {
	interceptors []grpc.StreamServerInterceptor
}

type unaryServerInterceptor struct {
	interceptors []grpc.UnaryServerInterceptor
}

// NewStreamServerInterceptor returns a new `StreamServerInterceptor`.
// It initializes its interceptor chain with `arr`.
func NewStreamServerInterceptor(arr ...grpc.StreamServerInterceptor) StreamServerInterceptor {
	return &streamServerInterceptor{
		interceptors: arr,
	}
}

func chainStream(current grpc.StreamServerInterceptor, info *grpc.StreamServerInfo, next grpc.StreamHandler) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return current(srv, stream, info, next)
	}
}

// StreamInterceptor chains all added interceptors into a single
// `grpc.StreamServerInterceptor`.
//
// The `handler` passed to each interceptor is either the next interceptor or,
// for the last element of the chain, the target method.
func (si streamServerInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// TODO: Find a more efficient way
		interceptor := handler
		for idx := len(si.interceptors) - 1; idx >= 0; idx-- {
			interceptor = chainStream(si.interceptors[idx], info, interceptor)
		}
		return interceptor(srv, ss)
	}
}

// AddGRPCStreamInterceptor adds `arr` to the chain of interceptors.
func (si *streamServerInterceptor) AddGRPCStreamInterceptor(arr ...grpc.StreamServerInterceptor) StreamServerInterceptor {
	si.interceptors = append(si.interceptors, arr...)
	return si
}

// AddStreamInterceptor is a convenient way for adding `StreamServerInterceptor`
// to the chain of interceptors. It only calls the method `StreamInterceptor`
// for each of them and append the return value to the chain.
func (si *streamServerInterceptor) AddStreamInterceptor(arr ...StreamServerInterceptor) StreamServerInterceptor {
	for _, i := range arr {
		si.interceptors = append(si.interceptors, i.StreamInterceptor())
	}
	return si
}

// NewUnaryServerInterceptor returns a new `UnaryServerInterceptor`.
// It initializes its interceptor chain with `arr`.
func NewUnaryServerInterceptor(arr ...grpc.UnaryServerInterceptor) UnaryServerInterceptor {
	return &unaryServerInterceptor{
		interceptors: arr,
	}
}

func chainUnary(current grpc.UnaryServerInterceptor, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return current(ctx, req, info, next)
	}
}

// UnaryInterceptor chains all added interceptors into a single
// `grpc.UnaryServerInterceptor`.
//
// The `handler` passed to each interceptor is either the next interceptor or,
// for the last element of the chain, the target method.
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

// AddGRPCUnaryInterceptor adds `arr` to the chain of interceptors.
func (ui *unaryServerInterceptor) AddGRPCUnaryInterceptor(arr ...grpc.UnaryServerInterceptor) UnaryServerInterceptor {
	ui.interceptors = append(ui.interceptors, arr...)
	return ui
}

// AddUnaryInterceptor is a convenient way for adding `UnaryServerInterceptor`
// to the chain of interceptors. It only calls the method `UnaryInterceptor`
// for each of them and append the return value to the chain.
func (ui *unaryServerInterceptor) AddUnaryInterceptor(arr ...UnaryServerInterceptor) UnaryServerInterceptor {
	for _, i := range arr {
		ui.interceptors = append(ui.interceptors, i.UnaryInterceptor())
	}
	return ui
}
