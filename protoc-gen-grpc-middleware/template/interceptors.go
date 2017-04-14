package template

import "text/template"

// Code template keys
const (
	pkgInterceptorsKey = "pkgInterceptors"
)

// Code templates
const (
	pkgInterceptorsCode = `
{{if .Symbols}}func init() {
	pkgInterceptors = append(
		pkgInterceptors,{{range .Symbols}}
		"{{.}}",{{end}}
	)
}{{end}}
`
)

func init() {
	template.Must(initCodeTpl.New(pkgInterceptorsKey).Parse(pkgInterceptorsCode))
}
