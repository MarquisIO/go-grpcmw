package generator

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

type Generator interface {
	Generate(*plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error)
}

type generator struct{}

const typePackageCode = `Interceptor_{{.Package}}`
const typeServiceCode = `Interceptor_{{.Package}}{{.Service}}`
const typeMethodCode = `{{if .ServerStream}}Stream{{else}}Unary{{end}}`

const methodCode = `func (s *server{{template "serviceType" .}}) {{.Method}}() grpcmw.{{template "methodType" .}}ServerInterceptor {
	method, ok := s.ServerInterceptor.(grpcmw.ServerInterceptorRegister).Get("{{.Method}}")
	if !ok {
		method = grpcmw.NewServerInterceptorRegister("{{.Method}}")
		s.ServerInterceptor.(grpcmw.ServerInterceptorRegister).Register(method)
	}
	return method.{{template "methodType" .}}ServerInterceptor()
}

func (s *client{{template "serviceType" .}}) {{.Method}}() grpcmw.{{template "methodType" .}}ClientInterceptor {
	method, ok := s.ClientInterceptor.(grpcmw.ClientInterceptorRegister).Get("{{.Method}}")
	if !ok {
		method = grpcmw.NewClientInterceptorRegister("{{.Method}}")
		s.ClientInterceptor.(grpcmw.ClientInterceptorRegister).Register(method)
	}
	return method.{{template "methodType" .}}ClientInterceptor()
}
`

const serviceCode = `type server{{template "serviceType" .}} struct {
	grpcmw.ServerInterceptor
}

type client{{template "serviceType" .}} struct {
	grpcmw.ClientInterceptor
}

func (i *server{{template "pkgType" .}}) {{.Service}}() *server{{template "serviceType" .}} {
	service, ok := i.ServerInterceptor.(grpcmw.ServerInterceptorRegister).Get("{{.Service}}")
	if !ok {
		service = grpcmw.NewServerInterceptorRegister("{{.Service}}")
		i.ServerInterceptor.(grpcmw.ServerInterceptorRegister).Register(service)
	}
	return &server{{template "serviceType" .}}{
		ServerInterceptor: service,
	}
}

func (i *client{{template "pkgType" .}}) {{.Service}}() *client{{template "serviceType" .}} {
	service, ok := i.ClientInterceptor.(grpcmw.ClientInterceptorRegister).Get("{{.Service}}")
	if !ok {
		service = grpcmw.NewClientInterceptorRegister("{{.Service}}")
		i.ClientInterceptor.(grpcmw.ClientInterceptorRegister).Register(service)
	}
	return &client{{template "serviceType" .}}{
		ClientInterceptor: service,
	}
}

{{range .Methods}}{{template "method" .}}{{end}}
`

const pkgCode = `type server{{template "pkgType" .}} struct {
	grpcmw.ServerInterceptor
}

type client{{template "pkgType" .}} struct {
	grpcmw.ClientInterceptor
}

func GetPackageServerInterceptors(router grpcmw.ServerRouter) *server{{template "pkgType" .}} {
	register := router.GetRegister()
	lvl, ok := register.Get("{{.Package}}")
	if !ok {
		lvl = grpcmw.NewServerInterceptorRegister("{{.Package}}")
		register.Register(lvl)
	}
	return &server{{template "pkgType" .}}{
		ServerInterceptor: lvl,
	}
}

func GetPackageClientInterceptors(router grpcmw.ClientRouter) *client{{template "pkgType" .}} {
	register := router.GetRegister()
	lvl, ok := register.Get("{{.Package}}")
	if !ok {
		lvl = grpcmw.NewClientInterceptorRegister("{{.Package}}")
		register.Register(lvl)
	}
	return &client{{template "pkgType" .}}{
		ClientInterceptor: lvl,
	}
}
`

const rootCode = `package {{.Package}}

import (
	grpcmw "github.com/MarquisIO/BKND-gRPCMiddleware/grpcmw"
)

{{if .DefinePackageLevel}}{{template "pkg" .}}{{end}}
{{range .Services}}{{template "service" .}}{{end}}
`

var rootCodeTpl = template.Must(template.New("code").Parse(rootCode))

func init() {
	template.Must(rootCodeTpl.New("pkg").Parse(pkgCode))
	template.Must(rootCodeTpl.New("service").Parse(serviceCode))
	template.Must(rootCodeTpl.New("method").Parse(methodCode))
	template.Must(rootCodeTpl.New("pkgType").Parse(typePackageCode))
	template.Must(rootCodeTpl.New("serviceType").Parse(typeServiceCode))
	template.Must(rootCodeTpl.New("methodType").Parse(typeMethodCode))
}

func New() Generator {
	return &generator{}
}

func (g generator) getResponseFile(src *descriptor.FileDescriptorProto, definePackageLevel bool) (dest *plugin.CodeGeneratorResponse_File) {
	if services := src.GetService(); len(services) > 0 {
		buf := new(bytes.Buffer)
		dest = &plugin.CodeGeneratorResponse_File{}
		srcName := src.GetName()
		destName := strings.TrimSuffix(srcName, filepath.Ext(srcName)) + ".pb.mw.go"
		dest.Name = &destName
		rootCodeTpl.Execute(buf, getTemplateData(src, definePackageLevel))
		ct := buf.String()
		dest.Content = &ct
	}
	return dest
}

func (g generator) Generate(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	res := &plugin.CodeGeneratorResponse{}
	definePackageLevel := true
	for _, src := range req.GetProtoFile() {
		if dest := g.getResponseFile(src, definePackageLevel); dest != nil {
			res.File = append(res.File, dest)
			definePackageLevel = false
		}
	}
	return res, nil
}
