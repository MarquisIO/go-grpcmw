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

type TemplateData struct {
	DefinePackageLevel bool
	Package            string
	Services           []Service
}

func getTemplateData(src *descriptor.FileDescriptorProto, definePackageLevel bool) *TemplateData {
	pkg := &TemplateData{
		Package:            src.GetPackage(),
		DefinePackageLevel: definePackageLevel,
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
