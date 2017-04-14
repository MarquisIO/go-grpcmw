package descriptor

import (
	annotations "github.com/MarquisIO/BKND-gRPCMiddleware/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type Service struct {
	Package      string
	Service      string
	Methods      []*Method
	Interceptors *Interceptors
}

func GetService(pb *descriptor.ServiceDescriptorProto, pkg string) (s *Service, err error) {
	methods := pb.GetMethod()
	s = &Service{
		Package: pkg,
		Service: pb.GetName(),
		Methods: make([]*Method, len(methods)),
	}
	if pb.Options != nil {
		if s.Interceptors, err = GetInterceptors(pb.Options, annotations.E_ServiceInterceptors); err != nil {
			return nil, err
		}
	}
	for idx, method := range methods {
		if s.Methods[idx], err = GetMethod(method, s.Service, pkg); err != nil {
			return nil, err
		}
	}
	return
}