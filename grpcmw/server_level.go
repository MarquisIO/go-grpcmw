package grpcmw

type ServerInterceptorLevel interface {
	UnaryServerInterceptor
	StreamServerInterceptor
	Index() string
}

type ServerInterceptorLevelRegister interface {
	ServerInterceptorLevel
	AddSubLevel(level ServerInterceptorLevel)
	Get(key string) (ServerInterceptorLevel, bool)
}

type LowerServerInterceptorLevel struct {
	UnaryServerInterceptor
	StreamServerInterceptor
	index string
}

type HigherServerInterceptorLevel struct {
	ServerInterceptorLevel
	sublevels map[string]ServerInterceptorLevel
}

func NewServerInterceptorLevel(index string) ServerInterceptorLevel {
	return &LowerServerInterceptorLevel{
		UnaryServerInterceptor:  NewUnaryServerInterceptor(),
		StreamServerInterceptor: NewStreamServerInterceptor(),
		index: index,
	}
}

func (l LowerServerInterceptorLevel) Index() string {
	return l.index
}

func NewServerInterceptorLevelRegister(index string) ServerInterceptorLevelRegister {
	return &HigherServerInterceptorLevel{
		ServerInterceptorLevel: NewServerInterceptorLevel(index),
		sublevels:              make(map[string]ServerInterceptorLevel),
	}
}

func (l HigherServerInterceptorLevel) Get(key string) (ServerInterceptorLevel, bool) {
	sub, exists := l.sublevels[key]
	return sub, exists
}

func (l *HigherServerInterceptorLevel) AddSubLevel(level ServerInterceptorLevel) {
	l.sublevels[level.Index()] = level
}
