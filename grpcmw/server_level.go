package grpcmw

import "google.golang.org/grpc"

// ServerInterceptor represent a server interceptor that uses both
// `UnaryServerInterceptor` and `StreamServerInterceptor` and that can be
// indexed.
type ServerInterceptor interface {
	AddGRPCUnaryInterceptor(i ...grpc.UnaryServerInterceptor) ServerInterceptor
	AddUnaryInterceptor(i ...UnaryServerInterceptor) ServerInterceptor
	UnaryServerInterceptor() UnaryServerInterceptor
	AddGRPCStreamInterceptor(i ...grpc.StreamServerInterceptor) ServerInterceptor
	AddStreamInterceptor(i ...StreamServerInterceptor) ServerInterceptor
	StreamServerInterceptor() StreamServerInterceptor
	Index() string
}

// ServerInterceptorRegister represents a register of `ServerInterceptor`,
// indexing them by using their method `Index`.
// It also implements `ServerInterceptor`.
type ServerInterceptorRegister interface {
	ServerInterceptor
	Register(level ServerInterceptor)
	Get(key string) (ServerInterceptor, bool)
}

type lowerServerInterceptor struct {
	unaries UnaryServerInterceptor
	streams StreamServerInterceptor
	index   string
}

type higherServerInterceptorLevel struct {
	ServerInterceptor
	sublevels map[string]ServerInterceptor
}

// NewServerInterceptor initializes a new `ServerInterceptor` with `index`
// as its index. It initializes the underlying `UnaryServerInterceptor` and
// `StreamServerInterceptor`
func NewServerInterceptor(index string) ServerInterceptor {
	return &lowerServerInterceptor{
		unaries: NewUnaryServerInterceptor(),
		streams: NewStreamServerInterceptor(),
		index:   index,
	}
}

// Index returns the index of the `ServerInterceptor`.
func (l lowerServerInterceptor) Index() string {
	return l.index
}

// AddGRPCUnaryInterceptor calls `AddGRPCInterceptor` of the underlying
// `UnaryServerInterceptor`. It returns the current instance of
// `ServerInterceptor` to allow chaining.
func (l *lowerServerInterceptor) AddGRPCUnaryInterceptor(arr ...grpc.UnaryServerInterceptor) ServerInterceptor {
	l.unaries.AddGRPCInterceptor(arr...)
	return l
}

// AddUnaryInterceptor calls `AddInterceptor` of the underlying
// `UnaryServerInterceptor`. It returns the current instance of
// `ServerInterceptor` to allow chaining.
func (l *lowerServerInterceptor) AddUnaryInterceptor(arr ...UnaryServerInterceptor) ServerInterceptor {
	l.unaries.AddInterceptor(arr...)
	return l
}

// UnaryServerInterceptor returns the underlying instance of
// `UnaryServerInterceptor`.
func (l *lowerServerInterceptor) UnaryServerInterceptor() UnaryServerInterceptor {
	return l.unaries
}

// AddGRPCStreamInterceptor calls `AddGRPCInterceptor` of the underlying
// `StreamServerInterceptor`. It returns the current instance of
// `ServerInterceptor` to allow chaining.
func (l *lowerServerInterceptor) AddGRPCStreamInterceptor(arr ...grpc.StreamServerInterceptor) ServerInterceptor {
	l.streams.AddGRPCInterceptor(arr...)
	return l
}

// AddStreamInterceptor calls `AddGRPCInterceptor` of the underlying
// `StreamServerInterceptor`. It returns the current instance of
// `ServerInterceptor` to allow chaining.
func (l *lowerServerInterceptor) AddStreamInterceptor(arr ...StreamServerInterceptor) ServerInterceptor {
	l.streams.AddInterceptor(arr...)
	return l
}

// StreamServerInterceptor returns the underlying instance of
// `StreamServerInterceptor`.
func (l *lowerServerInterceptor) StreamServerInterceptor() StreamServerInterceptor {
	return l.streams
}

// NewServerInterceptorRegister initializes a `ServerInterceptorRegister` with
// an empty register and `index` as index as its index.
func NewServerInterceptorRegister(index string) ServerInterceptorRegister {
	return &higherServerInterceptorLevel{
		ServerInterceptor: NewServerInterceptor(index),
		sublevels:         make(map[string]ServerInterceptor),
	}
}

// Get returns the `ServerInterceptor` registered at the index `key`. If nothing
// is found, it returns (nil, false).
func (l higherServerInterceptorLevel) Get(key string) (interceptor ServerInterceptor, exists bool) {
	sub, exists := l.sublevels[key]
	return sub, exists
}

// Register registers `level` indexing it by calling its method `Index`.
// It overwrites any interceptor that has already been registered at this index.
func (l *higherServerInterceptorLevel) Register(level ServerInterceptor) {
	l.sublevels[level.Index()] = level
}
