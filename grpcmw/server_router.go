package grpcmw

import (
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type ServerRouter interface {
	ServerInterceptorLevelRegister
	UnaryRouteResolver() grpc.UnaryServerInterceptor
	StreamRouteResolver() grpc.StreamServerInterceptor
	AddUnaryServerInterceptor(route string, interceptor ...grpc.UnaryServerInterceptor) error
	AddStreamServerInterceptor(route string, interceptor ...grpc.StreamServerInterceptor) error
}

type GRPCServerRouter struct {
	ServerInterceptorLevelRegister
}

var (
	routeRegexp = regexp.MustCompile("\\/(?:(.+)\\.(?:(.+)\\/(.+)|(.+))|(.+))")
)

func NewServerRouter() ServerRouter {
	return &GRPCServerRouter{
		ServerInterceptorLevelRegister: NewServerInterceptorLevelRegister("global"),
	}
}

func resolveRec(pathTokens []string, lvl ServerInterceptorLevel, cb func(lvl ServerInterceptorLevel), force bool) (ServerInterceptorLevel, error) {
	if cb != nil {
		cb(lvl)
	}
	if len(pathTokens) == 0 || len(pathTokens[0]) == 0 {
		return lvl, nil
	}
	reg, ok := lvl.(ServerInterceptorLevelRegister)
	if !ok {
		return nil, fmt.Errorf("Level %s do not implement grpcmw.ServerInterceptorLevelRegister", lvl.Index())
	}
	sub, exists := reg.Get(pathTokens[0])
	if !exists {
		if force {
			if len(pathTokens) == 1 {
				sub = NewServerInterceptorLevel(pathTokens[0])
			} else {
				sub = NewServerInterceptorLevelRegister(pathTokens[0])
			}
			reg.AddSubLevel(sub)
		} else {
			return nil, nil
		}
	}
	return resolveRec(pathTokens[1:], sub, cb, force)
}

func resolve(route string, lvl ServerInterceptorLevel, cb func(lvl ServerInterceptorLevel), force bool) (ServerInterceptorLevel, error) {
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
	return resolveRec(tokens, lvl, cb, force)
}

func (r *GRPCServerRouter) UnaryRouteResolver() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		interceptor := NewUnaryServerInterceptor()
		_, err := resolve(info.FullMethod, r.ServerInterceptorLevelRegister, func(lvl ServerInterceptorLevel) {
			interceptor.AddUnaryInterceptor(lvl)
		}, false)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, err.Error())
		}
		return interceptor.UnaryInterceptor()(ctx, req, info, handler)
	}
}

func (r *GRPCServerRouter) StreamRouteResolver() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		interceptor := NewStreamServerInterceptor()
		_, err := resolve(info.FullMethod, r.ServerInterceptorLevelRegister, func(lvl ServerInterceptorLevel) {
			interceptor.AddStreamInterceptor(lvl)
		}, false)
		if err != nil {
			return grpc.Errorf(codes.Internal, err.Error())
		}
		return interceptor.StreamInterceptor()(srv, ss, info, handler)
	}
}

func (r *GRPCServerRouter) AddUnaryServerInterceptor(route string, interceptors ...grpc.UnaryServerInterceptor) error {
	lvl, err := resolve(route, r.ServerInterceptorLevelRegister, nil, true)
	if err != nil {
		return err
	}
	lvl.AddGRPCUnaryInterceptor(interceptors...)
	return nil
}

func (r *GRPCServerRouter) AddStreamServerInterceptor(route string, interceptors ...grpc.StreamServerInterceptor) error {
	lvl, err := resolve(route, r.ServerInterceptorLevelRegister, nil, true)
	if err != nil {
		return err
	}
	lvl.AddGRPCStreamInterceptor(interceptors...)
	return nil
}
