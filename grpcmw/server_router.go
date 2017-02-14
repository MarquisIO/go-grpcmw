package grpcmw

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ServerRouter represents a server router which allows to add interceptors to
// routes, resolve these latter and use the appropriate chain of interceptors.
type ServerRouter interface {
	ServerInterceptorRegister
	UnaryResolver() grpc.UnaryServerInterceptor
	StreamResolver() grpc.StreamServerInterceptor
	AddUnaryServerInterceptor(route string, interceptor ...grpc.UnaryServerInterceptor) error
	AddStreamServerInterceptor(route string, interceptor ...grpc.StreamServerInterceptor) error
}

type serverRouter struct {
	ServerInterceptorRegister
}

// NewServerRouter initializes a `ServerRouter`.
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
func NewServerRouter() ServerRouter {
	return &serverRouter{
		ServerInterceptorRegister: NewServerInterceptorRegister("global"),
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
		return nil, fmt.Errorf("Level %s do not implement grpcmw.ServerInterceptorRegister", lvl.Index())
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
	return resolveServerInterceptorRec(tokens, lvl, cb, force)
}

// UnaryResolver returns a `grpc.UnaryServerInterceptor` that resolves the route
// of the request through the four levels of middlewares and imbricates them.
func (r *serverRouter) UnaryResolver() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		interceptor := NewUnaryServerInterceptor()
		_, err := resolveServerInterceptor(info.FullMethod, r.ServerInterceptorRegister, func(lvl ServerInterceptor) {
			interceptor.AddUnaryInterceptor(lvl)
		}, false)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, err.Error())
		}
		return interceptor.UnaryInterceptor()(ctx, req, info, handler)
	}
}

// StreamResolver returns a `grpc.StreamServerInterceptor` that resolves the
// route of the request through the four levels of middlewares and imbricates
// them.
func (r *serverRouter) StreamResolver() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		interceptor := NewStreamServerInterceptor()
		_, err := resolveServerInterceptor(info.FullMethod, r.ServerInterceptorRegister, func(lvl ServerInterceptor) {
			interceptor.AddStreamInterceptor(lvl)
		}, false)
		if err != nil {
			return grpc.Errorf(codes.Internal, err.Error())
		}
		return interceptor.StreamInterceptor()(srv, ss, info, handler)
	}
}

// AddUnaryServerInterceptor parses the given route to get the right level and
// adds the `interceptors` to it. The route should be of the following format:
// - for the package level: /package
// - for the service level: /package.service
// - for the method level: /package.service/method
func (r *serverRouter) AddUnaryServerInterceptor(route string, interceptors ...grpc.UnaryServerInterceptor) error {
	lvl, err := resolveServerInterceptor(route, r.ServerInterceptorRegister, nil, true)
	if err != nil {
		return err
	}
	lvl.AddGRPCUnaryInterceptor(interceptors...)
	return nil
}

// AddStreamServerInterceptor parses the given route to get the right level and
// adds the `interceptors` to it. The route should be of the following format:
// - for the package level: /package
// - for the service level: /package.service
// - for the method level: /package.service/method
func (r *serverRouter) AddStreamServerInterceptor(route string, interceptors ...grpc.StreamServerInterceptor) error {
	lvl, err := resolveServerInterceptor(route, r.ServerInterceptorRegister, nil, true)
	if err != nil {
		return err
	}
	lvl.AddGRPCStreamInterceptor(interceptors...)
	return nil
}
