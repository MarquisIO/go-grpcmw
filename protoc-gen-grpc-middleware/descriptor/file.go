package descriptor

import (
	"github.com/MarquisIO/go-grpcmw/annotations"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// File represents a protobuf file.
type File struct {
	Package      string
	Name         string
	Services     []*Service
	Interceptors *Interceptors
}

// GetFile parses `pb` and builds a `File` object from it.
// If the file does not define any service nor any interceptor option, it does
// not return anything.
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
	if f.Interceptors == nil && len(f.Services) == 0 {
		return nil, nil
	}
	return
}
