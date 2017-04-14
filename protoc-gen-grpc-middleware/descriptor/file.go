package descriptor

import (
	annotations "github.com/MarquisIO/BKND-gRPCMiddleware/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type File struct {
	Package      string
	Name         string
	Services     []*Service
	Interceptors *Interceptors
}

func GetFile(pb *descriptor.FileDescriptorProto) (f *File, err error) {
	services := pb.GetService()
	f = &File{
		Name:     pb.GetName(),
		Package:  pb.GetPackage(),
		Services: make([]*Service, len(services)),
	}
	if pb.Options != nil {
		if f.Interceptors, err = GetInterceptors(pb.Options, annotations.E_PackageInterceptors); err != nil {
			return nil, err
		}
	}
	for idx, service := range services {
		if f.Services[idx], err = GetService(service, f.Package); err != nil {
			return nil, err
		}
	}
	if (f.Interceptors == nil || len(f.Interceptors.Symbols) == 0) && len(f.Services) == 0 {
		return nil, nil
	}
	return
}
