package md

import (
	"os"
	"strings"
	"text/template"

	"github.com/mh-cbon/testndoc"
)

// Export the documentation to MD format.
func Export(recorded testndoc.APIDoc, dest string) error {
	os.Remove(dest)
	f, err2 := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0755)
	if err2 != nil {
		return err2
	}
	defer f.Close()
	t, err := template.New("").Funcs(map[string]interface{}{"lower": strings.ToLower}).Parse(index)
	if err != nil {
		return err
	}
	return t.Execute(f, recorded)
}

var index = `
# API DOC

## TOC
{{range .SortedEP}}
- [{{.ParameterizedPath}}](#{{.ParameterizedPath}}){{range .SortedRequests}}
  - [{{.Title}}](#{{.Title}}){{end}}

{{end}}

{{range .SortedEP}}
### {{.ParameterizedPath}}
{{.Doc}}
{{range .SortedRequests}}
#### {{.Title}}
{{.Doc}}

__[{{.Response.Code}}] {{.Request.Method}}__ {{.Request.URL}}

{{- if .HasURLParams}}

##### Url Parameters

| Key | Value |
| --- | --- |
{{range $k, $v := .Request.URLParams}}| {{$k}} | {{$v}} |{{end}}
{{end}}

{{- if .HasGetParams}}

##### GET Parameters

| Key | Value |
| --- | --- |
{{range $k, $v := .Request.GetParams}}| {{$k}} | {{$v}} |{{end}}
{{end}}

{{- if .HasPostParams}}

##### POST Parameters

| Key | Value |
| --- | --- |
{{range $k, $v := .Request.PostParams}}| {{$k}} | {{$v}} |{{end}}
{{end}}

{{- if .HasRequestHeaders}}

##### Request Headers

| Key | Value |
| --- | --- |
{{range $k, $v := .Request.Headers}}| {{$k}} | {{$v}} |{{end}}
{{end}}

{{- if .HasRequestBody}}

##### Request Body
` + "```" + `
{{.GetRequestBody}}
` + "```" + `
{{end}}

{{- if .HasResponseHeaders}}

##### Response Headers

| Key | Value |
| --- | --- |
{{range $k, $v := .Response.Headers}}| {{$k}} | {{$v}} |{{end}}
{{end}}

{{- if .HasResponseBody}}

##### Response Body
` + "```" + `
{{.GetResponseBody}}
` + "```" + `
{{end}}
[TOP](#api-doc)
___________________
{{end}}
{{end}}`
