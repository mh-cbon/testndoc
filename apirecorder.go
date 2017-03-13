package testndoc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
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
		st1 := SliceFramesUntil(st, "github.com/mh-cbon/backup/api.(*ApiRequestsResponses).FuncHandler")
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

// SliceFramesUntil ...
func SliceFramesUntil(frames []StackFrame, until string) []StackFrame {
	ret := []StackFrame{}
	for _, v := range frames {
		if len(v.Fn) >= len(until) && v.Fn[0:len(until)] == until {
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
