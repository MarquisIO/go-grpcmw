package grpcmw

// ServerInterceptor represent a server interceptor that exposes both
// `UnaryServerInterceptor` and `StreamServerInterceptor` and that can be
// indexed.
type ServerInterceptor interface {
	UnaryServerInterceptor
	StreamServerInterceptor
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
	UnaryServerInterceptor
	StreamServerInterceptor
	index string
}

type higherServerInterceptorLevel struct {
	ServerInterceptor
	sublevels map[string]ServerInterceptor
}

// NewServerInterceptor initializes a new `ServerInterceptor` with `index`
// as its index.
func NewServerInterceptor(index string) ServerInterceptor {
	return &lowerServerInterceptor{
		UnaryServerInterceptor:  NewUnaryServerInterceptor(),
		StreamServerInterceptor: NewStreamServerInterceptor(),
		index: index,
	}
}

// Index returns the index of the `ServerInterceptor`.
func (l lowerServerInterceptor) Index() string {
	return l.index
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
