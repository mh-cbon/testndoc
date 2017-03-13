package testndoc

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

// following is to be refactored later.
// its a duplicate of test2doc which can be largely simplified.

// URL ...
type URL struct {
	rawURL            *url.URL
	ParameterizedPath string
	Parameters        []Parameter
}

// NewURL ...
func NewURL(req *http.Request) *URL {
	u := &URL{
		rawURL: req.URL,
	}
	u.ParameterizedPath, u.Parameters = paramPath(req)
	return u
}

func paramPath(req *http.Request) (string, []Parameter) {
	uri, err := url.QueryUnescape(req.URL.Path)
	if err != nil {
		// fall back to unescaped uri
		uri = req.URL.Path
	}

	vars := mux.Vars(req)
	params := []Parameter{}

	for k, v := range vars {
		uri = strings.Replace(uri, "/"+v, "/{"+k+"}", 1)
		params = append(params, MakeParameter(k, v))
	}

	var queryKeys []string
	queryParams := req.URL.Query()

	for k, vs := range queryParams {
		queryKeys = append(queryKeys, k)

		// just take first value
		params = append(params, MakeParameter(k, vs[0]))
	}

	var queryKeysStr string
	if len(queryKeys) > 0 {
		queryKeysStr = "{?" + strings.Join(queryKeys, ",") + "}"
	}

	uri = uri + queryKeysStr

	return uri, params
}

// ParameterType ...
type ParameterType int

// ...
const (
	Number ParameterType = iota
	String
	Boolean
)

const (
	numberRe = `^[0-9\.]+$`
	boolRe   = `^(?:[tT][rR][uU][eE]|[fF][aA][lL][sS][eE])$`
)

var (
	parameterTmpl *template.Template
	parameterFmt  = `    + {{.Name}}: {{.Value.Quote}} ({{.Type.String}}){{with .Description}} - {{.}}{{end}}`
)

func init() {
	parameterTmpl = template.Must(template.New("parameter").Parse(parameterFmt))
}

// Parameter ...
type Parameter struct {
	Name        string
	Description string
	Value       ParameterValue
	Type        ParameterType
	IsRequired  bool

	// TODO:
	// DefaultValue
}

// MakeParameter ...
func MakeParameter(key, val string) Parameter {
	return Parameter{
		Name:       key,
		Value:      ParameterValue(val),
		Type:       paramType(val),
		IsRequired: true, // assume anything in route URL is required
		// query params are a different story
	}
}

// ParameterValue ...
type ParameterValue string

// Quote ...
func (val ParameterValue) Quote() (qval string) {
	if len(val) > 0 {
		qval = fmt.Sprintf("`%s`", string(val))
	}

	return
}

func paramType(val string) ParameterType {
	if isBool(val) {
		return Boolean
	} else if isNumber(val) {
		return Number
	} else {
		return String
	}
}

func isBool(str string) bool {
	re := regexp.MustCompile(boolRe)
	return re.MatchString(str)
}

func isNumber(str string) bool {
	re := regexp.MustCompile(numberRe)
	return re.MatchString(str)
}

func (pt ParameterType) String() string {
	switch pt {
	case Number:
		return "number"
	case Boolean:
		return "boolean"
	default:
		return "string"
	}
}
