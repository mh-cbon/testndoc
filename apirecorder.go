package testndoc

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"runtime"
	"strings"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/types/typeutil"
)

// APIRecorder ...
type APIRecorder map[string][]RequestResponse

// Add ...
func (a *APIRecorder) Add(r RequestResponse) {
	t := *a
	if _, ok := t[r.ParameterizedPath]; ok {
		t[r.ParameterizedPath] = append(t[r.ParameterizedPath], r)
	} else {
		t[r.ParameterizedPath] = []RequestResponse{r}
	}
}

// Save ...
func (a *APIRecorder) Save(to string) error {
	out, err := os.Create(to)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	return enc.Encode(a)
}

// Load ...
func (a *APIRecorder) Load(from string) error {
	b, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, a)
}

// LoadMetadata load api metadata information
// by digging ans searching your source code,
// looking for anotations and comments.
// You shall have recorded(or loaded) requests/responses before this.
func (a *APIRecorder) LoadMetadata() (APIDoc, error) {
	ret := APIDoc{}
	for eppath, records := range *a {

		ep := APIEndPoint{}
		for i, record := range records {

			Pkg, Type, Func := record.HandlerFn.Split()
			prog, err := loadPkg(Pkg)
			if err != nil {
				return ret, err
			}
			hDoc := findTypeDoc(prog, Pkg, Type, Func)

			TestPkg, TestType, TestFunc := record.TestFn.Split()
			TestProg, err := loadPkg(TestPkg)
			if err != nil {
				return ret, err
			}
			tDoc := findTypeDoc(TestProg, TestPkg, TestType, TestFunc)

			if i == 0 {
				ep.Path = record.ParameterizedPath
				ep.ParameterizedPath = eppath
				ep.HandlerFn = APIEndPointFn{StackFrame: record.HandlerFn, Doc: hDoc}
				ep.ReqRes = []APIRequestReponse{}
			}
			ep.ReqRes = append(ep.ReqRes, APIRequestReponse{
				RequestResponse: record,
				TestEP:          APIEndPointFn{StackFrame: record.TestFn, Doc: tDoc},
			})
		}
		ret.EndPoints = append(ret.EndPoints, ep)
	}
	return ret, nil
}

var loadedProgs = map[string]*loader.Program{}

func loadPkg(path string) (*loader.Program, error) {
	if x, ok := loadedProgs[path]; ok {
		return x, nil
	}
	var conf loader.Config
	conf.ParserMode = parser.ParseComments
	conf.ImportWithTests(path)
	prog, err := conf.Load()
	if err != nil {
		return prog, err
	}
	loadedProgs[path] = prog
	return prog, nil
}

// todo: return err if not found or similar
func findTypeDoc(prog *loader.Program, path string, Type string, Func string) string {
	pkg := prog.Package(path).Pkg
	if pkg == nil {
		log.Printf("%v not found", path)
		return ""
	}

	// look for a type's method
	if Type != "" {
		obj := pkg.Scope().Lookup(Type)
		if obj == nil {
			log.Printf("%s.%s not found\n", pkg.Path(), Type)
			return ""
		}
		for _, sel := range typeutil.IntuitiveMethodSet(obj.Type(), nil) {
			if sel.Obj().Name() == Func {
				_, paths, _ := prog.PathEnclosingInterval(sel.Obj().Pos(), sel.Obj().Pos())
				for _, n := range paths {
					switch n := n.(type) {
					case *ast.GenDecl:
						return n.Doc.Text()
					case *ast.FuncDecl:
						return n.Doc.Text()
					}
				}
			}
		}
	} else {
		// search a package func
		obj := pkg.Scope().Lookup(Func)
		if obj == nil {
			log.Printf("%s.%s not found\n", pkg.Path(), Func)
			return ""
		}
		_, paths, _ := prog.PathEnclosingInterval(obj.Pos(), obj.Pos())
		for _, n := range paths {
			switch n := n.(type) {
			case *ast.GenDecl:
				return n.Doc.Text()
			case *ast.FuncDecl:
				return n.Doc.Text()
			}
		}
	}
	return ""
}

// FuncHandler ...
func (a *APIRecorder) FuncHandler(
	h func(http.ResponseWriter, *http.Request),
	extractor func(*http.Request) map[string]string,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		u := NewURL(req)

		buf, _ := ioutil.ReadAll(req.Body)

		req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		resp := NewResponseRecorder()
		h(resp, req)
		copyHeader(w.Header(), resp.Header())
		w.WriteHeader(resp.Code)
		resBody := resp.Body.Bytes()

		w.Write(resBody)

		rr := RequestResponse{
			Request: Request{
				URL:        req.URL.String(),
				Method:     req.Method,
				Body:       string(buf),
				URLParams:  map[string]string{},
				GetParams:  map[string][]string{},
				PostParams: map[string][]string{},
				Headers:    map[string][]string{},
			},
			Response: Response{
				Code:    resp.Code,
				Body:    string(resBody),
				Headers: map[string][]string{},
			},
			ParameterizedPath: u.ParameterizedPath,
			TestFn:            resp.TestFn,
			HandlerFn:         resp.HandlerFn,
		}
		for _, p := range u.Parameters {
			rr.Request.URLParams[p.Name] = string(p.Value)
		}
		for k, v := range req.URL.Query() {
			rr.Request.GetParams[k] = v
		}
		for k, v := range req.Form {
			rr.Request.PostParams[k] = v
		}
		for k, v := range req.Header {
			rr.Request.Headers[k] = v
		}
		for k, v := range resp.Header() {
			rr.Response.Headers[k] = v
		}

		a.Add(rr)
	}
}

// RequestResponse ...
type RequestResponse struct {
	Request           Request
	Response          Response
	ParameterizedPath string
	HandlerFn         StackFrame
	TestFn            StackFrame
}

// Request ...
type Request struct {
	URL        string
	Method     string
	URLParams  map[string]string
	Headers    map[string][]string
	GetParams  map[string][]string
	PostParams map[string][]string
	Body       string
}

// Response ...
type Response struct {
	Code    int
	Headers map[string][]string
	Body    string
}

// ResponseRecorder ...
type ResponseRecorder struct {
	*httptest.ResponseRecorder
	HandlerFn StackFrame
	TestFn    StackFrame
	recorded  bool
}

// NewResponseRecorder ...
func NewResponseRecorder() *ResponseRecorder {
	return &ResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
}

// Write always succeeds and writes to rw.Body, if not nil.
func (rw *ResponseRecorder) record() {
	if !rw.recorded {
		st := makeStackTrace()
		st1 := SliceFramesUntil(st, "github.com/mh-cbon/testndoc.(*APIRecorder).FuncHandler")
		rw.HandlerFn = st1[len(st1)-1]
		st2 := SliceFramesUntil(st, "testing.tRunner")
		rw.TestFn = st2[len(st2)-1]
	}
	rw.recorded = true
}

// Write always succeeds and writes to rw.Body, if not nil.
func (rw *ResponseRecorder) Write(buf []byte) (int, error) {
	rw.record()
	return rw.ResponseRecorder.Write(buf)
}

// WriteString always succeeds and writes to rw.Body, if not nil.
func (rw *ResponseRecorder) WriteString(str string) (int, error) {
	rw.record()
	return rw.ResponseRecorder.WriteString(str)
}

// WriteHeader sets rw.Code. After it is called, changing rw.Header
// will not affect rw.HeaderMap.
func (rw *ResponseRecorder) WriteHeader(code int) {
	rw.record()
	rw.ResponseRecorder.WriteHeader(code)
}

// StackFrame ...
type StackFrame struct {
	File string
	Fn   string
}

// matches public func ~~ (*Config).(the/path/pkg.RmTask)-fm
var methodR = regexp.MustCompile(`(?i)(\([^)]+\))\.(\([^)]+\))-fm$`)

// matches private func ~~ (*Config).RmTask-fm
var methodR2 = regexp.MustCompile(`(?i)(\([^)]+\))\.([^-]+)-fm$`)

// matches public func ~~ github.com/mh-cbon/backup/api.TestTaskStart
var funcPublic = regexp.MustCompile(`(?i)(.+)\.([^.]+)$`)

func (s StackFrame) Split() (string, string, string) {
	Pkg := ""
	Type := ""
	Func := ""

	methodRes := methodR.FindAllStringSubmatch(s.Fn, -1)
	if len(methodRes) > 0 {
		d := strings.Index(s.Fn, "(")
		Pkg = s.Fn[:d-1]

		Type = methodRes[0][1]       // something like (*DisksInfo)
		Type = Type[1 : len(Type)-1] // removes parenthesis
		if strings.HasPrefix(Type, "*") {
			Type = Type[1:]
		}
		Func = methodRes[0][2]       // ~ (github.com/mh-cbon/backup/api.read)
		Func = Func[1 : len(Func)-1] // removes parenthesis
		e := strings.LastIndex(Func, ".")
		Func = Func[e+1:]
	} else {
		methodRes = methodR2.FindAllStringSubmatch(s.Fn, -1)
		if len(methodRes) > 0 {
			d := strings.Index(s.Fn, "(")
			Pkg = s.Fn[:d-1]

			Type = methodRes[0][1]
			Type = Type[1 : len(Type)-1] // removes parenthesis
			if strings.HasPrefix(Type, "*") {
				Type = Type[1:]
			}
			Func = methodRes[0][2]
		} else {
			funcRes := funcPublic.FindAllStringSubmatch(s.Fn, -1)
			Pkg = funcRes[0][1]
			Func = funcRes[0][2]
		}
	}
	return Pkg, Type, Func
}

// SliceFramesUntil ...
func SliceFramesUntil(frames []StackFrame, until string) []StackFrame {
	ret := []StackFrame{}
	for _, v := range frames {
		if len(v.Fn) >= len(until) && strings.Index(v.Fn, until) > -1 {
			break
		}
		ret = append(ret, v)
	}
	return ret
}

func makeStackTrace() []StackFrame {
	ret := []StackFrame{}
	i := 1
	max := 150

	var pc uintptr
	var file, fnName string
	var ok bool

	for i < max {
		pc, file, _, ok = runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		fnName = fn.Name()

		ret = append(ret, StackFrame{
			File: file,
			Fn:   fnName,
		})
		i++
	}
	return ret
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
