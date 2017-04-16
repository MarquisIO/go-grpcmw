package template

import "text/template"

// Code template keys
const (
	pkgInterceptorsKey = "pkgInterceptors"
)

// Code templates
const (
	pkgInterceptorsCode = `
{{if .Indexes}}func init() {
	pkgInterceptors = append(
		pkgInterceptors,{{range .Indexes}}
		"{{.}}",{{end}}
	)
}{{end}}
`
)

func init() {
	template.Must(initCodeTpl.New(pkgInterceptorsKey).Parse(pkgInterceptorsCode))
}
