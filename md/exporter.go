package md

import (
	"os"
	"text/template"

	"github.com/mh-cbon/testndoc"
)

// Export the documentation to MD format.
func Export(recorded testndoc.APIDoc, dest string) error {
	t, err := template.New("").Parse(index)
	if err != nil {
		return err
	}
	return t.Execute(os.Stdout, recorded)
}

var index = `
# API DOC

{{range .EndPoints}}
## {{.ParameterizedPath}}
	{{.Doc}}
	{{range .ReqRes}}
		#### {{.ParameterizedPath}}
		{{.Doc}}
	{{end}}
{{end}}
`
