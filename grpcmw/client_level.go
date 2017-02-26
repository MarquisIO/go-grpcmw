package grpcmw

import "google.golang.org/grpc"

// ClientInterceptor represent a client interceptor that uses both
// `UnaryClientInterceptor` and `StreamClientInterceptor` and that can be
// indexed.
type ClientInterceptor interface {
	AddGRPCUnaryInterceptor(i ...grpc.UnaryClientInterceptor) ClientInterceptor
	AddUnaryInterceptor(i ...UnaryClientInterceptor) ClientInterceptor
	UnaryClientInterceptor() UnaryClientInterceptor
	AddGRPCStreamInterceptor(i ...grpc.StreamClientInterceptor) ClientInterceptor
	AddStreamInterceptor(i ...StreamClientInterceptor) ClientInterceptor
	StreamClientInterceptor() StreamClientInterceptor
	Index() string
}

// ClientInterceptorRegister represents a register of `ClientInterceptor`,
// indexing them by using their method `Index`.
// It also implements `ClientInterceptor`.
type ClientInterceptorRegister interface {
	ClientInterceptor
	Register(level ClientInterceptor)
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
}

// NewClientInterceptor initializes a new `ClientInterceptor` with `index`
// as its index. It initializes the underlying `UnaryClientInterceptor` and
// `StreamClientInterceptor`
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

// NewClientInterceptorRegister initializes a `ClientInterceptorRegister` with
// an empty register and `index` as index as its index.
func NewClientInterceptorRegister(index string) ClientInterceptorRegister {
	return &higherClientInterceptorLevel{
		ClientInterceptor: NewClientInterceptor(index),
		sublevels:         make(map[string]ClientInterceptor),
	}
}

// Get returns the `ClientInterceptor` registered at the index `key`. If nothing
// is found, it returns (nil, false).
func (l higherClientInterceptorLevel) Get(key string) (interceptor ClientInterceptor, exists bool) {
	sub, exists := l.sublevels[key]
	return sub, exists
}

// Register registers `level` indexing it by calling its method `Index`.
// It overwrites any interceptor that has already been registered at this index.
func (l *higherClientInterceptorLevel) Register(level ClientInterceptor) {
	l.sublevels[level.Index()] = level
}
