package generator

import (
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type Method struct {
	Package      string
	Service      string
	Method       string
	ClientStream bool
	ServerStream bool
}

type Service struct {
	Package string
	Service string
	Methods []Method
}

type Package struct {
	Package  string
	Services []Service
}

func GetPackage(src *descriptor.FileDescriptorProto) *Package {
	pkg := &Package{
		Package: src.GetPackage(),
	}
	for _, service := range src.GetService() {
		s := Service{
			Package: src.GetPackage(),
			Service: service.GetName(),
		}
		for _, method := range service.GetMethod() {
			s.Methods = append(s.Methods, Method{
				Package:      src.GetPackage(),
				Service:      service.GetName(),
				Method:       method.GetName(),
				ClientStream: method.GetClientStreaming(),
				ServerStream: method.GetServerStreaming(),
			})
		}
		pkg.Services = append(pkg.Services, s)
	}
	return pkg
}
