package grpcmw

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ServerRouter represents route resolver that allows to use the appropriate
// chain of interceptors for a given gRPC request with an interceptor register.
type ServerRouter interface {
	// GetRegister returns the interceptor register of the router.
	GetRegister() ServerInterceptorRegister
	// UnaryResolver returns a `grpc.UnaryServerInterceptor` that uses the
	// appropriate chain of interceptors with the given unary gRPC request.
	UnaryResolver() grpc.UnaryServerInterceptor
	// StreamResolver returns a `grpc.StreamServerInterceptor` that uses the
	// appropriate chain of interceptors with the given stream gRPC request.
	StreamResolver() grpc.StreamServerInterceptor
}

type serverRouter struct {
	interceptors ServerInterceptorRegister
}

// NewServerRouter initializes a `ServerRouter`.
// This implementation is based on the official route format used by gRPC as
// defined here :
// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#appendix-a---grpc-for-protobuf
//
// Based on this format, this implementation splits the interceptors into four
// levels:
// - the global level: these are the interceptors called at each request.
// - the package level: these are the interceptors called at each request to a
//   service from the corresponding package.
// - the service level: these are the interceptors called at each request to a
//   method from the corresponding service.
// - the method level: these are the interceptors called at each request to the
//   specific method.
func NewServerRouter() ServerRouter {
	return &serverRouter{
		interceptors: NewServerInterceptorRegister("global"),
	}
}

func resolveServerInterceptorRec(pathTokens []string, lvl ServerInterceptor, cb func(lvl ServerInterceptor), force bool) (ServerInterceptor, error) {
	if cb != nil {
		cb(lvl)
	}
	if len(pathTokens) == 0 || len(pathTokens[0]) == 0 {
		return lvl, nil
	}
	reg, ok := lvl.(ServerInterceptorRegister)
	if !ok {
		return nil, fmt.Errorf("Level %s does not implement grpcmw.ServerInterceptorRegister", lvl.Index())
	}
	sub, exists := reg.Get(pathTokens[0])
	if !exists {
		if force {
			if len(pathTokens) == 1 {
				sub = NewServerInterceptor(pathTokens[0])
			} else {
				sub = NewServerInterceptorRegister(pathTokens[0])
			}
			reg.Register(sub)
		} else {
			return nil, nil
		}
	}
	return resolveServerInterceptorRec(pathTokens[1:], sub, cb, force)
}

func resolveServerInterceptor(route string, lvl ServerInterceptor, cb func(lvl ServerInterceptor), force bool) (ServerInterceptor, error) {
	// TODO: Find a more efficient way to resolve the route
	matchs := routeRegexp.FindStringSubmatch(route)
	if len(matchs) == 0 {
		return nil, errors.New("Invalid route")
	}
	return resolveServerInterceptorRec(matchs[1:], lvl, cb, force)
}

// UnaryResolver returns a `grpc.UnaryServerInterceptor` that uses the
// appropriate chain of interceptors with the given gRPC request.
func (r *serverRouter) UnaryResolver() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// TODO: Find a more efficient way to chain the interceptors
		interceptor := NewUnaryServerInterceptor()
		_, err := resolveServerInterceptor(info.FullMethod, r.interceptors, func(lvl ServerInterceptor) {
			interceptor.AddInterceptor(lvl.UnaryServerInterceptor())
		}, false)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, err.Error())
		}
		return interceptor.Interceptor()(ctx, req, info, handler)
	}
}

// StreamResolver returns a `grpc.StreamServerInterceptor` that uses the
// appropriate chain of interceptors with the given stream gRPC request.
func (r *serverRouter) StreamResolver() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// TODO: Find a more efficient way to chain the interceptors
		interceptor := NewStreamServerInterceptor()
		_, err := resolveServerInterceptor(info.FullMethod, r.interceptors, func(lvl ServerInterceptor) {
			interceptor.AddInterceptor(lvl.StreamServerInterceptor())
		}, false)
		if err != nil {
			return grpc.Errorf(codes.Internal, err.Error())
		}
		return interceptor.Interceptor()(srv, ss, info, handler)
	}
}

// GetRegister returns the underlying `ServerInterceptorRegister` which is the
// global level in the interceptor chain.
func (r *serverRouter) GetRegister() ServerInterceptorRegister {
	return r.interceptors
}
