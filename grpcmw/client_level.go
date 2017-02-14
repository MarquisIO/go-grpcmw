package grpcmw

// ClientInterceptor represent a client interceptor that exposes both
// `UnaryClientInterceptor` and `StreamClientInterceptor` and that can be
// indexed.
type ClientInterceptor interface {
	UnaryClientInterceptor
	StreamClientInterceptor
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
	UnaryClientInterceptor
	StreamClientInterceptor
	index string
}

type higherClientInterceptorLevel struct {
	ClientInterceptor
	sublevels map[string]ClientInterceptor
}

// NewClientInterceptor initializes a new `ClientInterceptor` with `index`
// as its index.
func NewClientInterceptor(index string) ClientInterceptor {
	return &lowerClientInterceptor{
		UnaryClientInterceptor:  NewUnaryClientInterceptor(),
		StreamClientInterceptor: NewStreamClientInterceptor(),
		index: index,
	}
}

// Index returns the index of the `ClientInterceptor`.
func (l lowerClientInterceptor) Index() string {
	return l.index
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
