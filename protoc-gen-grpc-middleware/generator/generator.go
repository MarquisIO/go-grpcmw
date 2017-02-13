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

const pkgCode = `package {{.Package}}

import (
	grpcmw "github.com/MarquisIO/BKND-gRPCMiddleware/grpcmw"
	grpc "google.golang.org/grpc"
)
{{range .Services}}
type {{.Service}}ServerRouter struct {
	grpcmw.ServerRouter
}

func Get{{.Service}}ServerRouter(r grpcmw.ServerRouter) *{{.Service}}ServerRouter {
	return &{{.Service}}ServerRouter{r}
}

func (r *{{.Service}}ServerRouter) AddUnaryInterceptor(interceptor ...grpc.UnaryServerInterceptor) error {
	return r.AddUnaryServerInterceptor("/{{.Package}}.{{.Service}}", interceptor...)
}

func (r *{{.Service}}ServerRouter) AddStreamInterceptor(interceptor ...grpc.StreamServerInterceptor) error {
	return r.AddStreamServerInterceptor("/{{.Package}}.{{.Service}}", interceptor...)
}
{{range .Methods}}{{if .ServerStream}}
func (r *{{.Service}}ServerRouter) AddInterceptorTo{{.Method}}(interceptor ...grpc.StreamServerInterceptor) error {
	return r.AddStreamServerInterceptor("/{{.Package}}.{{.Service}}/{{.Method}}", interceptor...)
}
{{else}}
func (r *{{.Service}}ServerRouter) AddInterceptorTo{{.Method}}(interceptor ...grpc.UnaryServerInterceptor) error {
	return r.AddUnaryServerInterceptor("/{{.Package}}.{{.Service}}/{{.Method}}", interceptor...)
}
{{end}}{{end}}{{end}}
`

var pkgCodeTpl = template.Must(template.New("package").Parse(pkgCode))

func New() Generator {
	return &generator{}
}

func (g generator) getResponseFile(src *descriptor.FileDescriptorProto) (dest *plugin.CodeGeneratorResponse_File) {
	if services := src.GetService(); len(services) > 0 {
		buf := new(bytes.Buffer)
		dest = &plugin.CodeGeneratorResponse_File{}
		srcName := src.GetName()
		destName := strings.TrimSuffix(srcName, filepath.Ext(srcName)) + ".pb.mw.go"
		dest.Name = &destName
		pkgCodeTpl.Execute(buf, GetPackage(src))
		ct := buf.String()
		dest.Content = &ct
	}
	return dest
}

func (g generator) Generate(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	res := &plugin.CodeGeneratorResponse{}
	for _, src := range req.GetProtoFile() {
		if dest := g.getResponseFile(src); dest != nil {
			res.File = append(res.File, dest)
		}
	}
	return res, nil
}
