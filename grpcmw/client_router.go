package grpcmw

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ClientRouter represents route resolver that allows to use the appropriate
// chain of interceptors for a given gRPC request with an interceptor register.
type ClientRouter interface {
	// GetRegister returns the interceptor register of the router.
	GetRegister() ClientInterceptorRegister
	// SetRegister sets the interceptor register of the router.
	SetRegister(reg ClientInterceptorRegister)
	// UnaryResolver returns a `grpc.UnaryClientInterceptor` that uses the
	// appropriate chain of interceptors with the given unary gRPC request.
	UnaryResolver() grpc.UnaryClientInterceptor
	// StreamResolver returns a `grpc.StreamClientInterceptor` that uses the
	// appropriate chain of interceptors with the given stream gRPC request.
	StreamResolver() grpc.StreamClientInterceptor
}

type clientRouter struct {
	interceptors ClientInterceptorRegister
}

// NewClientRouter initializes a `ClientRouter`.
// This implementation is based on the official route format used by gRPC as
// defined here :
// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#appendix-a---grpc-for-protobuf
//
// Based on this format, this implementation splits the interceptors into four
// levels:
//   - the global level: these are the interceptors called at each request.
//   - the package level: these are the interceptors called at each request to a
//     service from the corresponding package.
//   - the service level: these are the interceptors called at each request to a
//     method from the corresponding service.
//   - the method level: these are the interceptors called at each request to the
//     specific method.
func NewClientRouter() ClientRouter {
	return &clientRouter{
		interceptors: NewClientInterceptorRegister("global"),
	}
}

func resolveClientInterceptorRec(pathTokens []string, lvl ClientInterceptor, cb func(lvl ClientInterceptor), force bool) (ClientInterceptor, error) {
	if cb != nil {
		cb(lvl)
	}
	if len(pathTokens) == 0 || len(pathTokens[0]) == 0 {
		return lvl, nil
	}
	reg, ok := lvl.(ClientInterceptorRegister)
	if !ok {
		return nil, fmt.Errorf("Level %s does not implement grpcmw.ClientInterceptorRegister", lvl.Index())
	}
	sub, exists := reg.Get(pathTokens[0])
	if !exists {
		if force {
			if len(pathTokens) == 1 {
				sub = NewClientInterceptor(pathTokens[0])
			} else {
				sub = NewClientInterceptorRegister(pathTokens[0])
			}
			reg.Register(sub)
		} else {
			return nil, nil
		}
	}
	return resolveClientInterceptorRec(pathTokens[1:], sub, cb, force)
}

func resolveClientInterceptor(route string, lvl ClientInterceptor, cb func(lvl ClientInterceptor), force bool) (ClientInterceptor, error) {
	// TODO: Find a more efficient way to resolve the route
	matchs := routeRegexp.FindStringSubmatch(route)
	if len(matchs) == 0 {
		return nil, errors.New("Invalid route")
	}
	return resolveClientInterceptorRec(matchs[1:], lvl, cb, force)
}

// UnaryResolver returns a `grpc.UnaryClientInterceptor` that uses the
// appropriate chain of interceptors with the given gRPC request.
func (r *clientRouter) UnaryResolver() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// TODO: Find a more efficient way to chain the interceptors
		interceptor := NewUnaryClientInterceptor()
		_, err := resolveClientInterceptor(method, r.interceptors, func(lvl ClientInterceptor) {
			interceptor.AddInterceptor(lvl.UnaryClientInterceptor())
		}, false)
		if err != nil {
			return grpc.Errorf(codes.Internal, err.Error())
		}
		return interceptor.Interceptor()(ctx, method, req, reply, cc, invoker, opts...)
	}
}

// StreamResolver returns a `grpc.StreamClientInterceptor` that uses the
// appropriate chain of interceptors with the given stream gRPC request.
func (r *clientRouter) StreamResolver() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// TODO: Find a more efficient way to chain the interceptors
		interceptor := NewStreamClientInterceptor()
		_, err := resolveClientInterceptor(method, r.interceptors, func(lvl ClientInterceptor) {
			interceptor.AddInterceptor(lvl.StreamClientInterceptor())
		}, false)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, err.Error())
		}
		return interceptor.Interceptor()(ctx, desc, cc, method, streamer, opts...)
	}
}

// GetRegister returns the underlying `ClientInterceptorRegister` which is the
// global level in the interceptor chain.
func (r *clientRouter) GetRegister() ClientInterceptorRegister {
	return r.interceptors
}

// SetRegister sets the interceptor register of the router.
func (r *clientRouter) SetRegister(reg ClientInterceptorRegister) {
	r.interceptors = reg
}
