package grpcmw

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// StreamClientInterceptor represents a client interceptor for gRPC methods that
// return a stream. It allows chaining of `grpc.StreamClientInterceptor`
// and other `StreamClientInterceptor`.
type StreamClientInterceptor interface {
	Interceptor() grpc.StreamClientInterceptor
	AddGRPCInterceptor(i ...grpc.StreamClientInterceptor) StreamClientInterceptor
	AddInterceptor(i ...StreamClientInterceptor) StreamClientInterceptor
}

// UnaryClientInterceptor represents a client interceptor for gRPC methods that
// return a single value instead of a stream. It allows chaining of
// `grpc.UnaryClientInterceptor` and other `UnaryClientInterceptor`.
type UnaryClientInterceptor interface {
	Interceptor() grpc.UnaryClientInterceptor
	AddGRPCInterceptor(i ...grpc.UnaryClientInterceptor) UnaryClientInterceptor
	AddInterceptor(i ...UnaryClientInterceptor) UnaryClientInterceptor
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

// Interceptor chains all added interceptors into a single
// `grpc.StreamClientInterceptor`.
//
// The `streamer` passed to each interceptor is either the next interceptor or,
// for the last element of the chain, the target method.
func (si streamClientInterceptor) Interceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// TODO: Find a more efficient way
		interceptor := streamer
		for idx := len(si.interceptors) - 1; idx >= 0; idx-- {
			interceptor = chainStreamClientInterceptor(si.interceptors[idx], interceptor)
		}
		return interceptor(ctx, desc, cc, method, opts...)
	}
}

// AddGRPCInterceptor adds `arr` to the chain of interceptors.
func (si *streamClientInterceptor) AddGRPCInterceptor(arr ...grpc.StreamClientInterceptor) StreamClientInterceptor {
	si.interceptors = append(si.interceptors, arr...)
	return si
}

// AddInterceptor is a convenient way for adding `StreamClientInterceptor`
// to the chain of interceptors. It only calls the method `StreamInterceptor`
// for each of them and append the return value to the chain.
func (si *streamClientInterceptor) AddInterceptor(arr ...StreamClientInterceptor) StreamClientInterceptor {
	for _, i := range arr {
		si.interceptors = append(si.interceptors, i.Interceptor())
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

// Interceptor chains all added interceptors into a single
// `grpc.UnaryClientInterceptor`.
//
// The `streamer` passed to each interceptor is either the next interceptor or,
// for the last element of the chain, the target method.
func (ui *unaryClientInterceptor) Interceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// TODO: Find a more efficient way
		interceptor := invoker
		for idx := len(ui.interceptors) - 1; idx >= 0; idx-- {
			interceptor = chainUnaryClientInterceptor(ui.interceptors[idx], interceptor)
		}
		return interceptor(ctx, method, req, reply, cc, opts...)
	}
}

// AddGRPCInterceptor adds `arr` to the chain of interceptors.
func (ui *unaryClientInterceptor) AddGRPCInterceptor(arr ...grpc.UnaryClientInterceptor) UnaryClientInterceptor {
	ui.interceptors = append(ui.interceptors, arr...)
	return ui
}

// AddInterceptor is a convenient way for adding `UnaryClientInterceptor`
// to the chain of interceptors. It only calls the method `UnaryInterceptor`
// for each of them and append the return value to the chain.
func (ui *unaryClientInterceptor) AddInterceptor(arr ...UnaryClientInterceptor) UnaryClientInterceptor {
	for _, i := range arr {
		ui.interceptors = append(ui.interceptors, i.Interceptor())
	}
	return ui
}
