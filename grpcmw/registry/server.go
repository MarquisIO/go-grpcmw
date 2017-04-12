package registry

import (
	"sync"

	"github.com/MarquisIO/BKND-gRPCMiddleware/grpcmw"
)

var (
	serverLock     sync.Mutex
	serverRegistry = make(map[string]grpcmw.ServerInterceptor)
)

// GetServerInterceptor returns the `grpcmw.ServerInterceptor` registered at
// `index`. If nothing is at this `index`, it registers a new one using
// `grpcmw.NewServerInterceptor` and returns it.
// This is thread-safe.
func GetServerInterceptor(index string) grpcmw.ServerInterceptor {
	serverLock.Lock()
	defer serverLock.Unlock()
	intcp, ok := serverRegistry[index]
	if !ok {
		intcp = grpcmw.NewServerInterceptor(index)
		serverRegistry[index] = intcp
	}
	return intcp
}

// SetServerInterceptor registers `interceptor` at `index`. It replaces any
// interceptor that has been previously registered at this `index`.
// This is thread-safe.
func SetServerInterceptor(index string, interceptor grpcmw.ServerInterceptor) {
	serverLock.Lock()
	defer serverLock.Unlock()
	serverRegistry[index] = interceptor
}

// DeleteServerInterceptor deletes any interceptor registered at `index`.
// This is thread-safe.
func DeleteServerInterceptor(index string) {
	serverLock.Lock()
	defer serverLock.Unlock()
	delete(serverRegistry, index)
}
