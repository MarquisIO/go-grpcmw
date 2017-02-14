package grpcmw

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// StreamClientInterceptor represents a client interceptor for gRPC methods that
// return a stream. It allows chaining of `grpc.StreamClientInterceptor`
// and other `StreamClientInterceptor`.
type StreamClientInterceptor interface {
	StreamInterceptor() grpc.StreamClientInterceptor
	AddGRPCStreamInterceptor(i ...grpc.StreamClientInterceptor) StreamClientInterceptor
	AddStreamInterceptor(i ...StreamClientInterceptor) StreamClientInterceptor
}

// UnaryClientInterceptor represents a client interceptor for gRPC methods that
// return a single value instead of a stream. It allows chaining of
// `grpc.UnaryClientInterceptor` and other `UnaryClientInterceptor`.
type UnaryClientInterceptor interface {
	UnaryInterceptor() grpc.UnaryClientInterceptor
	AddGRPCUnaryInterceptor(i ...grpc.UnaryClientInterceptor) UnaryClientInterceptor
	AddUnaryInterceptor(i ...UnaryClientInterceptor) UnaryClientInterceptor
}

type streamClientInterceptor struct {
	interceptors []grpc.StreamClientInterceptor
}

type unaryClientInterceptor struct {
	interceptors []grpc.UnaryClientInterceptor
}

// NewStreamClientInterceptor returns a new `StreamClientInterceptor`.
// It initializes its interceptor chain with `arr`.
func NewStreamClientInterceptor(arr ...grpc.StreamClientInterceptor) StreamClientInterceptor {
	return &streamClientInterceptor{
		interceptors: arr,
	}
}

func chainStreamClientInterceptor(current grpc.StreamClientInterceptor, next grpc.Streamer) grpc.Streamer {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return current(ctx, desc, cc, method, next, opts...)
	}
}

// StreamInterceptor chains all added interceptors into a single
// `grpc.StreamClientInterceptor`.
//
// The `streamer` passed to each interceptor is either the next interceptor or,
// for the last element of the chain, the target method.
func (si streamClientInterceptor) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// TODO: Find a more efficient way
		interceptor := streamer
		for idx := len(si.interceptors) - 1; idx >= 0; idx-- {
			interceptor = chainStreamClientInterceptor(si.interceptors[idx], interceptor)
		}
		return interceptor(ctx, desc, cc, method, opts...)
	}
}

// AddGRPCStreamInterceptor adds `arr` to the chain of interceptors.
func (si *streamClientInterceptor) AddGRPCStreamInterceptor(arr ...grpc.StreamClientInterceptor) StreamClientInterceptor {
	si.interceptors = append(si.interceptors, arr...)
	return si
}

// AddStreamInterceptor is a convenient way for adding `StreamClientInterceptor`
// to the chain of interceptors. It only calls the method `StreamInterceptor`
// for each of them and append the return value to the chain.
func (si *streamClientInterceptor) AddStreamInterceptor(arr ...StreamClientInterceptor) StreamClientInterceptor {
	for _, i := range arr {
		si.interceptors = append(si.interceptors, i.StreamInterceptor())
	}
	return si
}

// NewUnaryClientInterceptor returns a new `UnaryClientInterceptor`.
// It initializes its interceptor chain with `arr`.
func NewUnaryClientInterceptor(arr ...grpc.UnaryClientInterceptor) UnaryClientInterceptor {
	return &unaryClientInterceptor{
		interceptors: arr,
	}
}

func chainUnaryClientInterceptor(current grpc.UnaryClientInterceptor, next grpc.UnaryInvoker) grpc.UnaryInvoker {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return current(ctx, method, req, reply, cc, next, opts...)
	}
}

// UnaryInterceptor chains all added interceptors into a single
// `grpc.UnaryClientInterceptor`.
//
// The `streamer` passed to each interceptor is either the next interceptor or,
// for the last element of the chain, the target method.
func (ui *unaryClientInterceptor) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// TODO: Find a more efficient way
		interceptor := invoker
		for idx := len(ui.interceptors) - 1; idx >= 0; idx-- {
			interceptor = chainUnaryClientInterceptor(ui.interceptors[idx], interceptor)
		}
		return interceptor(ctx, method, req, reply, cc, opts...)
	}
}

// AddGRPCUnaryInterceptor adds `arr` to the chain of interceptors.
func (ui *unaryClientInterceptor) AddGRPCUnaryInterceptor(arr ...grpc.UnaryClientInterceptor) UnaryClientInterceptor {
	ui.interceptors = append(ui.interceptors, arr...)
	return ui
}

// AddUnaryInterceptor is a convenient way for adding `UnaryClientInterceptor`
// to the chain of interceptors. It only calls the method `UnaryInterceptor`
// for each of them and append the return value to the chain.
func (ui *unaryClientInterceptor) AddUnaryInterceptor(arr ...UnaryClientInterceptor) UnaryClientInterceptor {
	for _, i := range arr {
		ui.interceptors = append(ui.interceptors, i.UnaryInterceptor())
	}
	return ui
}
