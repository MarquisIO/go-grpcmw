package grpcmw

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ClientRouter represents a client router which allows to add interceptors to
// routes, resolve these latter and use the appropriate chain of interceptors.
type ClientRouter interface {
	ClientInterceptorRegister
	UnaryResolver() grpc.UnaryClientInterceptor
	StreamResolver() grpc.StreamClientInterceptor
	AddUnaryClientInterceptor(route string, interceptor ...grpc.UnaryClientInterceptor) error
	AddStreamClientInterceptor(route string, interceptor ...grpc.StreamClientInterceptor) error
}

type clientRouter struct {
	ClientInterceptorRegister
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
		ClientInterceptorRegister: NewClientInterceptorRegister("global"),
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
		_, err := resolveClientInterceptor(method, r.ClientInterceptorRegister, func(lvl ClientInterceptor) {
			interceptor.AddUnaryInterceptor(lvl)
		}, false)
		if err != nil {
			return grpc.Errorf(codes.Internal, err.Error())
		}
		return interceptor.UnaryInterceptor()(ctx, method, req, reply, cc, invoker, opts...)
	}
}

// StreamResolver returns a `grpc.StreamClientInterceptor` that resolves the
// route of the request through the four levels of middlewares and imbricates
// them.
func (r *clientRouter) StreamResolver() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// TODO: Find a more efficient way to chain the interceptors
		interceptor := NewStreamClientInterceptor()
		_, err := resolveClientInterceptor(method, r.ClientInterceptorRegister, func(lvl ClientInterceptor) {
			interceptor.AddStreamInterceptor(lvl)
		}, false)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, err.Error())
		}
		return interceptor.StreamInterceptor()(ctx, desc, cc, method, streamer, opts...)
	}
}

// AddUnaryClientInterceptor parses the given route to get the right level and
// adds the `interceptors` to it. The route should be of the following format:
// - for the package level: /package
// - for the service level: /package.service
// - for the method level: /package.service/method
func (r *clientRouter) AddUnaryClientInterceptor(route string, interceptors ...grpc.UnaryClientInterceptor) error {
	lvl, err := resolveClientInterceptor(route, r.ClientInterceptorRegister, nil, true)
	if err != nil {
		return err
	}
	lvl.AddGRPCUnaryInterceptor(interceptors...)
	return nil
}

// AddStreamClientInterceptor parses the given route to get the right level and
// adds the `interceptors` to it. The route should be of the following format:
// - for the package level: /package
// - for the service level: /package.service
// - for the method level: /package.service/method
func (r *clientRouter) AddStreamClientInterceptor(route string, interceptors ...grpc.StreamClientInterceptor) error {
	lvl, err := resolveClientInterceptor(route, r.ClientInterceptorRegister, nil, true)
	if err != nil {
		return err
	}
	lvl.AddGRPCStreamInterceptor(interceptors...)
	return nil
}
