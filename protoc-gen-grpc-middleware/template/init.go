package template

import (
	"bytes"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/MarquisIO/BKND-gRPCMiddleware/protoc-gen-grpc-middleware/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// Code template keys
const (
	initKey = "init"
)

// Code templates
const (
	initCode = `package {{.Package}}

import (
	grpcmw   "github.com/MarquisIO/BKND-gRPCMiddleware/grpcmw"
	registry "github.com/MarquisIO/BKND-gRPCMiddleware/grpcmw/registry"
)

type server{{template "pkgType" .}} struct {
	grpcmw.ServerInterceptor
}

type client{{template "pkgType" .}} struct {
	grpcmw.ClientInterceptor
}

var pkgInterceptors []string
{{with .Interceptors}}{{template "pkgInterceptors" .}}{{end}}
func RegisterServerInterceptors(router grpcmw.ServerRouter) *server{{template "pkgType" .}} {
	register := router.GetRegister()
	lvl, ok := register.Get("{{.Package}}")
	if !ok {
		lvl = grpcmw.NewServerInterceptorRegister("{{.Package}}")
		register.Register(lvl)
		for _, interceptor := range pkgInterceptors {
			lvl.Merge(registry.GetServerInterceptor(interceptor))
		}
	}
	return &server{{template "pkgType" .}}{
		ServerInterceptor: lvl,
	}
}

func RegisterClientInterceptors(router grpcmw.ClientRouter) *client{{template "pkgType" .}} {
	register := router.GetRegister()
	lvl, ok := register.Get("{{.Package}}")
	if !ok {
		lvl = grpcmw.NewClientInterceptorRegister("{{.Package}}")
		register.Register(lvl)
		for _, interceptor := range pkgInterceptors {
			lvl.Merge(registry.GetClientInterceptor(interceptor))
		}
	}
	return &client{{template "pkgType" .}}{
		ClientInterceptor: lvl,
	}
}

{{range .Services}}{{template "service" .}}{{end}}
`
)

var initCodeTpl = template.Must(template.New(initKey).Parse(initCode))

// Apply applies the given package descriptors and generates the appropriate
// code using go templates.
func Apply(pkgs map[string][]*descriptor.File) (*plugin.CodeGeneratorResponse, error) {
	res := &plugin.CodeGeneratorResponse{}
	for _, files := range pkgs {
		for idx, file := range files {
			buf := new(bytes.Buffer)
			dest := &plugin.CodeGeneratorResponse_File{}
			destName := strings.TrimSuffix(file.Name, filepath.Ext(file.Name)) + ".pb.mw.go"
			dest.Name = &destName
			templateKey := pkgKey
			if idx == 0 {
				templateKey = initKey
			}
			if err := initCodeTpl.ExecuteTemplate(buf, templateKey, file); err != nil {
				return nil, err
			}
			ct := buf.String()
			dest.Content = &ct
			res.File = append(res.File, dest)
		}
	}
	return res, nil
}
