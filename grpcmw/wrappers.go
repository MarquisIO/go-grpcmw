package grpcmw

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// ServerStreamWrapper represents a wrapper for `grpc.ServerStream` that allows
// to modify the context.
type ServerStreamWrapper struct {
	grpc.ServerStream
	ctx context.Context
}

// WrapServerStream returns checks if `ss` is already a `*ServerStreamWrapper`.
// If it is, it returns the `ss`, otherwise it returns a new wrapper for
// `grpc.ServerStream`.
func WrapServerStream(ss grpc.ServerStream) *ServerStreamWrapper {
	if ret, ok := ss.(*ServerStreamWrapper); ok {
		return ret
	}
	return &ServerStreamWrapper{
		ServerStream: ss,
		ctx:          ss.Context(),
	}
}

// Context returns the context of the wrapper0
func (w ServerStreamWrapper) Context() context.Context {
	return w.ctx
}

// SetContext set the context of the wrapper to `ctx`.
func (w *ServerStreamWrapper) SetContext(ctx context.Context) {
	w.ctx = ctx
}
