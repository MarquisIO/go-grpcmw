package descriptor

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"

	annotations "github.com/MarquisIO/go-grpcmw/proto"
)

// Method represents a method from a grpc service.
type Method struct {
	Package      string
	Service      string
	Method       string
	ClientStream bool
	ServerStream bool
	Interceptors *Interceptors
}

// GetMethod parses `pb` and builds from it a `Method` object.
func GetMethod(pb *descriptor.MethodDescriptorProto, service, pkg string) (method *Method, err error) {
	method = &Method{
		Package:      pkg,
		Service:      service,
		Method:       pb.GetName(),
		ClientStream: pb.GetClientStreaming(),
		ServerStream: pb.GetServerStreaming(),
	}
	if pb.Options != nil {
		if method.Interceptors, err = GetInterceptors(pb.Options, annotations.E_MethodInterceptors); err != nil {
			return nil, err
		}
	}
	return
}
