package grpcmw

import (
	"sync"

	"google.golang.org/grpc"
)

// ServerInterceptor represent a server interceptor that uses both
// `UnaryServerInterceptor` and `StreamServerInterceptor` and that can be
// indexed.
type ServerInterceptor interface {
	// AddGRPCUnaryInterceptor adds given unary interceptors to the chain.
	AddGRPCUnaryInterceptor(i ...grpc.UnaryServerInterceptor) ServerInterceptor
	// AddUnaryInterceptor is a convenient way for adding `UnaryServerInterceptor`
	// to the chain of unary interceptors.
	AddUnaryInterceptor(i ...UnaryServerInterceptor) ServerInterceptor
	// UnaryServerInterceptor returns the chain of unary interceptors.
	UnaryServerInterceptor() UnaryServerInterceptor
	// AddGRPCStreamInterceptor adds given stream interceptors to the chain.
	AddGRPCStreamInterceptor(i ...grpc.StreamServerInterceptor) ServerInterceptor
	// AddStreamInterceptor is a convenient way for adding
	// `StreamServerInterceptor` to the chain of stream interceptors.
	AddStreamInterceptor(i ...StreamServerInterceptor) ServerInterceptor
	// StreamServerInterceptor returns the chain of stream interceptors.
	StreamServerInterceptor() StreamServerInterceptor
	// Merge merges the given interceptors with the current interceptor.
	Merge(interceptors ...ServerInterceptor) ServerInterceptor
	// Index returns the index of the `ServerInterceptor`.
	Index() string
}

// ServerInterceptorRegister represents a register of `ServerInterceptor`,
// indexing them by using their method `Index`.
// It also implements `ServerInterceptor`.
type ServerInterceptorRegister interface {
	ServerInterceptor
	// Register registers `level` at the index returned by its method `Index`.
	Register(level ServerInterceptor)
	// Get returns the `ServerInterceptor` registered at the index `key`. If
	// nothing is found, it returns (nil, false).
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
	lock      *sync.RWMutex
}

// NewServerInterceptor initializes a new `ServerInterceptor` with `index`
// as its index. It initializes the underlying `UnaryServerInterceptor` and
// `StreamServerInterceptor`.
// This implementation is thread-safe.
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

// Merge merges the given interceptors with the current interceptor.
func (l *lowerServerInterceptor) Merge(interceptors ...ServerInterceptor) ServerInterceptor {
	for _, interceptor := range interceptors {
		l.AddUnaryInterceptor(interceptor.UnaryServerInterceptor()).
			AddStreamInterceptor(interceptor.StreamServerInterceptor())
	}
	return l
}

// NewServerInterceptorRegister initializes a `ServerInterceptorRegister` with
// an empty register and `index` as index as its index.
// This implementation is thread-safe.
func NewServerInterceptorRegister(index string) ServerInterceptorRegister {
	return &higherServerInterceptorLevel{
		ServerInterceptor: NewServerInterceptor(index),
		sublevels:         make(map[string]ServerInterceptor),
		lock:              &sync.RWMutex{},
	}
}

// Get returns the `ServerInterceptor` registered at the index `key`. If nothing
// is found, it returns (nil, false).
func (l higherServerInterceptorLevel) Get(key string) (interceptor ServerInterceptor, exists bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	interceptor, exists = l.sublevels[key]
	return
}

// Register registers `level` at the index returned by its method `Index`.
// It overwrites any interceptor that has already been registered at this index.
func (l *higherServerInterceptorLevel) Register(level ServerInterceptor) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.sublevels[level.Index()] = level
}
