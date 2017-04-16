package grpcmw

import (
	"sync"

	"google.golang.org/grpc"
)

// ClientInterceptor represents a client interceptor that uses both
// `UnaryClientInterceptor` and `StreamClientInterceptor` and that can be
// indexed.
type ClientInterceptor interface {
	// AddGRPCUnaryInterceptor adds given unary interceptors to the chain.
	AddGRPCUnaryInterceptor(i ...grpc.UnaryClientInterceptor) ClientInterceptor
	// AddUnaryInterceptor is a convenient way for adding `UnaryClientInterceptor`
	// to the chain of unary interceptors.
	AddUnaryInterceptor(i ...UnaryClientInterceptor) ClientInterceptor
	// UnaryClientInterceptor returns the chain of unary interceptors.
	UnaryClientInterceptor() UnaryClientInterceptor
	// AddGRPCStreamInterceptor adds given stream interceptors to the chain.
	AddGRPCStreamInterceptor(i ...grpc.StreamClientInterceptor) ClientInterceptor
	// AddStreamInterceptor is a convenient way for adding
	// `StreamClientInterceptor` to the chain of stream interceptors.
	AddStreamInterceptor(i ...StreamClientInterceptor) ClientInterceptor
	// StreamClientInterceptor returns the chain of stream interceptors.
	StreamClientInterceptor() StreamClientInterceptor
	// Merge merges the given interceptors with the current interceptor.
	Merge(i ...ClientInterceptor) ClientInterceptor
	// Index returns the index of the `ClientInterceptor`.
	Index() string
}

// ClientInterceptorRegister represents a register of `ClientInterceptor`,
// indexing them by using their method `Index`.
// It also implements `ClientInterceptor`.
type ClientInterceptorRegister interface {
	ClientInterceptor
	// Register registers `level` at the index returned by its method `Index`.
	Register(level ClientInterceptor)
	// Get returns the `ClientInterceptor` registered at the index `key`. If
	// nothing is found, it returns (nil, false).
	Get(key string) (ClientInterceptor, bool)
}

type lowerClientInterceptor struct {
	unaries UnaryClientInterceptor
	streams StreamClientInterceptor
	index   string
}

type higherClientInterceptorLevel struct {
	ClientInterceptor
	sublevels map[string]ClientInterceptor
	lock      *sync.RWMutex
}

// NewClientInterceptor initializes a new `ClientInterceptor` with `index`
// as its index. It initializes the underlying `UnaryClientInterceptor` and
// `StreamClientInterceptor`.
// This implementation is thread-safe.
func NewClientInterceptor(index string) ClientInterceptor {
	return &lowerClientInterceptor{
		unaries: NewUnaryClientInterceptor(),
		streams: NewStreamClientInterceptor(),
		index:   index,
	}
}

// Index returns the index of the `ClientInterceptor`.
func (l lowerClientInterceptor) Index() string {
	return l.index
}

// AddGRPCUnaryInterceptor calls `AddGRPCInterceptor` of the underlying
// `UnaryClientInterceptor`. It returns the current instance of
// `ClientInterceptor` to allow chaining.
func (l *lowerClientInterceptor) AddGRPCUnaryInterceptor(arr ...grpc.UnaryClientInterceptor) ClientInterceptor {
	l.unaries.AddGRPCInterceptor(arr...)
	return l
}

// AddUnaryInterceptor calls `AddInterceptor` of the underlying
// `UnaryClientInterceptor`. It returns the current instance of
// `ClientInterceptor` to allow chaining.
func (l *lowerClientInterceptor) AddUnaryInterceptor(arr ...UnaryClientInterceptor) ClientInterceptor {
	l.unaries.AddInterceptor(arr...)
	return l
}

// UnaryClientInterceptor returns the underlying instance of
// `UnaryClientInterceptor`.
func (l *lowerClientInterceptor) UnaryClientInterceptor() UnaryClientInterceptor {
	return l.unaries
}

// AddGRPCStreamInterceptor calls `AddGRPCInterceptor` of the underlying
// `StreamClientInterceptor`. It returns the current instance of
// `ClientInterceptor` to allow chaining.
func (l *lowerClientInterceptor) AddGRPCStreamInterceptor(arr ...grpc.StreamClientInterceptor) ClientInterceptor {
	l.streams.AddGRPCInterceptor(arr...)
	return l
}

// AddStreamInterceptor calls `AddGRPCInterceptor` of the underlying
// `StreamClientInterceptor`. It returns the current instance of
// `ClientInterceptor` to allow chaining.
func (l *lowerClientInterceptor) AddStreamInterceptor(arr ...StreamClientInterceptor) ClientInterceptor {
	l.streams.AddInterceptor(arr...)
	return l
}

// StreamClientInterceptor returns the underlying instance of
// `StreamClientInterceptor`.
func (l *lowerClientInterceptor) StreamClientInterceptor() StreamClientInterceptor {
	return l.streams
}

// Merge merges the given interceptors with the current interceptor.
func (l *lowerClientInterceptor) Merge(interceptors ...ClientInterceptor) ClientInterceptor {
	for _, interceptor := range interceptors {
		l.AddUnaryInterceptor(interceptor.UnaryClientInterceptor()).
			AddStreamInterceptor(interceptor.StreamClientInterceptor())
	}
	return l
}

// NewClientInterceptorRegister initializes a `ClientInterceptorRegister` with
// an empty register and `index` as index as its index.
// This implementation is thread-safe.
func NewClientInterceptorRegister(index string) ClientInterceptorRegister {
	return &higherClientInterceptorLevel{
		ClientInterceptor: NewClientInterceptor(index),
		sublevels:         make(map[string]ClientInterceptor),
		lock:              &sync.RWMutex{},
	}
}

// Get returns the `ClientInterceptor` registered at the index `key`. If nothing
// is found, it returns (nil, false).
func (l higherClientInterceptorLevel) Get(key string) (interceptor ClientInterceptor, exists bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	interceptor, exists = l.sublevels[key]
	return
}

// Register registers `level` at the index returned by its method `Index`.
// It overwrites any interceptor that has already been registered at this index.
func (l *higherClientInterceptorLevel) Register(level ClientInterceptor) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.sublevels[level.Index()] = level
}
