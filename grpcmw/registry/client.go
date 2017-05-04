package registry

import (
	"sync"

	"github.com/MarquisIO/go-grpcmw/grpcmw"
)

var (
	clientLock     sync.Mutex
	clientRegistry = make(map[string]grpcmw.ClientInterceptor)
)

// GetClientInterceptor returns the `grpcmw.ClientInterceptor` registered at
// `index`. If nothing is at this `index`, it registers a new one using
// `grpcmw.NewClientInterceptor` and returns it.
// This is thread-safe.
func GetClientInterceptor(index string) grpcmw.ClientInterceptor {
	clientLock.Lock()
	defer clientLock.Unlock()
	intcp, ok := clientRegistry[index]
	if !ok {
		intcp = grpcmw.NewClientInterceptor(index)
		clientRegistry[index] = intcp
	}
	return intcp
}

// SetClientInterceptor registers `interceptor` at `index`. It replaces any
// interceptor that has been previously registered at this `index`.
// This is thread-safe.
func SetClientInterceptor(index string, interceptor grpcmw.ClientInterceptor) {
	clientLock.Lock()
	defer clientLock.Unlock()
	clientRegistry[index] = interceptor
}

// DeleteClientInterceptor deletes any interceptor registered at `index`.
// This is thread-safe.
func DeleteClientInterceptor(index string) {
	clientLock.Lock()
	defer clientLock.Unlock()
	delete(clientRegistry, index)
}
