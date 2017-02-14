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

const routerCode = `import grpcmw "github.com/MarquisIO/BKND-gRPCMiddleware/grpcmw"

type {{.Package}}ServerRouter struct {
	grpcmw.ServerRouter
}

type {{.Package}}ClientRouter struct {
	grpcmw.ClientRouter
}

func GetServerRouter(r grpcmw.ServerRouter) *{{.Package}}ServerRouter {
	return &{{.Package}}ServerRouter{r}
}

func (r *{{.Package}}ServerRouter) AddUnaryInterceptorToPackage(interceptor ...grpc.UnaryServerInterceptor) error {
	return r.AddUnaryServerInterceptor("/{{.Package}}", interceptor...)
}

func (r *{{.Package}}ServerRouter) AddStreamInterceptorToPackage(interceptor ...grpc.StreamServerInterceptor) error {
	return r.AddStreamServerInterceptor("/{{.Package}}", interceptor...)
}

func GetClientRouter(r grpcmw.ClientRouter) *{{.Package}}ClientRouter {
	return &{{.Package}}ClientRouter{r}
}

func (r *{{.Package}}ClientRouter) AddUnaryInterceptorToPackage(interceptor ...grpc.UnaryClientInterceptor) error {
	return r.AddUnaryClientInterceptor("/{{.Package}}", interceptor...)
}

func (r *{{.Package}}ClientRouter) AddStreamInterceptorToPackage(interceptor ...grpc.StreamClientInterceptor) error {
	return r.AddStreamClientInterceptor("/{{.Package}}", interceptor...)
}`

const pkgCode = `package {{.Package}}

import (
	grpc "google.golang.org/grpc"
)
{{if .DefineRouter}}{{template "router" .}}{{end}}
{{range .Services}}
func (r *{{.Package}}ServerRouter) AddUnaryInterceptorToService{{.Service}}(interceptor ...grpc.StreamServerInterceptor) error {
	return r.AddStreamServerInterceptor("/{{.Package}}.{{.Service}}", interceptor...)
}

func (r *{{.Package}}ServerRouter) AddStreamInterceptorToService{{.Service}}(interceptor ...grpc.UnaryServerInterceptor) error {
	return r.AddUnaryServerInterceptor("/{{.Package}}.{{.Service}}", interceptor...)
}

func (r *{{.Package}}ClientRouter) AddUnaryInterceptorToService{{.Service}}(interceptor ...grpc.StreamClientInterceptor) error {
	return r.AddStreamClientInterceptor("/{{.Package}}.{{.Service}}", interceptor...)
}

func (r *{{.Package}}ClientRouter) AddStreamInterceptorToService{{.Service}}(interceptor ...grpc.UnaryClientInterceptor) error {
	return r.AddUnaryClientInterceptor("/{{.Package}}.{{.Service}}", interceptor...)
}
{{range .Methods}}{{if .ServerStream}}
func (r *{{.Package}}ServerRouter) AddInterceptorToMethod{{.Method}}(interceptor ...grpc.StreamServerInterceptor) error {
	return r.AddStreamServerInterceptor("/{{.Package}}.{{.Service}}/{{.Method}}", interceptor...)
}
{{else}}
func (r *{{.Package}}ServerRouter) AddInterceptorToMethod{{.Method}}(interceptor ...grpc.UnaryServerInterceptor) error {
	return r.AddUnaryServerInterceptor("/{{.Package}}.{{.Service}}/{{.Method}}", interceptor...)
}
{{end}}{{if .ClientStream}}
func (r *{{.Package}}ClientRouter) AddInterceptorToMethod{{.Method}}(interceptor ...grpc.StreamClientInterceptor) error {
	return r.AddStreamClientInterceptor("/{{.Package}}.{{.Service}}/{{.Method}}", interceptor...)
}
{{else}}
func (r *{{.Package}}ClientRouter) AddInterceptorToMethod{{.Method}}(interceptor ...grpc.UnaryClientInterceptor) error {
	return r.AddUnaryClientInterceptor("/{{.Package}}.{{.Service}}/{{.Method}}", interceptor...)
}
{{end}}{{end}}{{end}}`

var (
	pkgCodeTpl    = template.Must(template.New("package").Parse(pkgCode))
	routerCodeTpl = template.Must(pkgCodeTpl.New("router").Parse(routerCode))
)

func New() Generator {
	return &generator{}
}

func (g generator) getResponseFile(src *descriptor.FileDescriptorProto, defineRouter bool) (dest *plugin.CodeGeneratorResponse_File) {
	if services := src.GetService(); len(services) > 0 {
		buf := new(bytes.Buffer)
		dest = &plugin.CodeGeneratorResponse_File{}
		srcName := src.GetName()
		destName := strings.TrimSuffix(srcName, filepath.Ext(srcName)) + ".pb.mw.go"
		dest.Name = &destName
		pkgCodeTpl.Execute(buf, getTemplateData(src, defineRouter))
		ct := buf.String()
		dest.Content = &ct
	}
	return dest
}

func (g generator) Generate(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	res := &plugin.CodeGeneratorResponse{}
	defineRouter := true
	for _, src := range req.GetProtoFile() {
		if dest := g.getResponseFile(src, defineRouter); dest != nil {
			res.File = append(res.File, dest)
			defineRouter = false
		}
	}
	return res, nil
}
