package template

import (
	"text/template"
)

// Code template keys
const (
	pkgKey     = "pkg"
	pkgTypeKey = "pkgType"
)

// Code templates
const (
	pkgTypeCode = `Interceptor_{{.Package}}`

	pkgCode = `package {{.Package}}

import (
	grpcmw "github.com/MarquisIO/go-grpcmw/grpcmw"
	registry "github.com/MarquisIO/go-grpcmw/grpcmw/registry"
)

var (
	_ = registry.GetClientInterceptor
)

{{with .Interceptors}}{{template "pkgInterceptors" .}}{{end}}
{{range .Services}}{{template "service" .}}{{end}}
`
)

func init() {
	template.Must(initCodeTpl.New(pkgKey).Parse(pkgCode))
	template.Must(initCodeTpl.New(pkgTypeKey).Parse(pkgTypeCode))
}
