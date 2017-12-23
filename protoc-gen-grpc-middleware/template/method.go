package template

import "text/template"

// Code template keys
const (
	methodTypeKey = "methodType"
	methodKey     = "method"
)

// Code templates
const (
	methodTypeCode = `{{if .}}Stream{{else}}Unary{{end}}`

	methodCode = `
func (s *server{{template "serviceType" .}}) {{.Method}}() grpcmw.{{template "methodType" .Stream}}ServerInterceptor {
	method, ok := s.ServerInterceptor.(grpcmw.ServerInterceptorRegister).Get("{{.Method}}")
	if !ok {
		method = grpcmw.NewServerInterceptorRegister("{{.Method}}")
		s.ServerInterceptor.(grpcmw.ServerInterceptorRegister).Register(method)
	}
	return method.{{template "methodType" .Stream}}ServerInterceptor()
}

func (s *client{{template "serviceType" .}}) {{.Method}}() grpcmw.{{template "methodType" .Stream}}ClientInterceptor {
	method, ok := s.ClientInterceptor.(grpcmw.ClientInterceptorRegister).Get("{{.Method}}")
	if !ok {
		method = grpcmw.NewClientInterceptorRegister("{{.Method}}")
		s.ClientInterceptor.(grpcmw.ClientInterceptorRegister).Register(method)
	}
	return method.{{template "methodType" .Stream}}ClientInterceptor()
}`
)

func init() {
	template.Must(initCodeTpl.New(methodKey).Parse(methodCode))
	template.Must(initCodeTpl.New(methodTypeKey).Parse(methodTypeCode))
}
