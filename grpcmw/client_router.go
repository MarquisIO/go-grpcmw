package grpcmw

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ClientRouter represents a server router which allows to resolve routes from a
// middleware register and use the appropriate chain of interceptors.
type ClientRouter interface {
	GetRegister() ClientInterceptorRegister
	UnaryResolver() grpc.UnaryClientInterceptor
	StreamResolver() grpc.StreamClientInterceptor
}

type clientRouter struct {
	interceptors ClientInterceptorRegister
}

// NewClientRouter initializes a `ClientRouter`.
// This implementation is based on the official route format used by gRPC,
// which is the following:
// /package.service/method
//
// Based on this format, this implementation splits the middlewares into four
// levels:
// - the global level: these are the middlewares called at each request.
// - the package level: these are the middlewares called at each request to a
//   service from the corresponding package.
// - the service level: these are the middlewares called at each request to a
//   method from the corresponding service.
// - the method level: these are the middlewares called at each request to the
//   specific method.
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
		return nil, fmt.Errorf("Level %s do not implement grpcmw.ClientInterceptorRegister", lvl.Index())
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
	tokens := matchs[1:4]
	if len(matchs[4]) > 0 {
		tokens[1] = matchs[4]
	} else if len(matchs[5]) > 0 {
		tokens[0] = matchs[5]
	}
	return resolveClientInterceptorRec(tokens, lvl, cb, force)
}

// UnaryResolver returns a `grpc.UnaryClientInterceptor` that resolves the route
// of the request through the four levels of middlewares and imbricates them.
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

// StreamResolver returns a `grpc.StreamClientInterceptor` that resolves the
// route of the request through the four levels of middlewares and imbricates
// them.
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
// global level in the middleware chain.
func (r *clientRouter) GetRegister() ClientInterceptorRegister {
	return r.interceptors
}
