package testndoc

import (
	"encoding/json"
	"sort"
	"strings"
)

type APIDoc struct {
	EndPoints []APIEndPoint
}

func (a APIDoc) SortedEP() []APIEndPoint {
	ByAPIEndPoint(ByPath).Sort(a.EndPoints)
	return a.EndPoints
}

type APIEndPoint struct {
	Path              string
	ParameterizedPath string
	HandlerFn         APIEndPointFn
	ReqRes            []APIRequestReponse
}

func (a APIEndPoint) SortedRequests() []APIRequestReponse {
	ByAPIRequestReponse(ByTitle).Sort(a.ReqRes)
	return a.ReqRes
}

func (a APIEndPoint) Doc() string {
	return a.HandlerFn.Doc
}

func (a APIEndPoint) GetMethods(j string) string {
	k := []string{}
	for _, d := range a.ReqRes {
		k = append(k, d.Request.Method)
	}
	return strings.Join(k, j)
}

type APIEndPointFn struct {
	StackFrame
	Doc string
}

type APIRequestReponse struct {
	RequestResponse
	TestEP APIEndPointFn
}

func (a APIRequestReponse) Doc() string {
	return a.TestEP.Doc
}

func (a APIRequestReponse) Title() string {
	return a.TestEP.Fn
}

func (a APIRequestReponse) HasURLParams() bool {
	return len(a.Request.URLParams) > 0
}

func (a APIRequestReponse) HasGetParams() bool {
	return len(a.Request.GetParams) > 0
}

func (a APIRequestReponse) HasPostParams() bool {
	return len(a.Request.PostParams) > 0
}

func (a APIRequestReponse) HasRequestHeaders() bool {
	return len(a.Request.Headers) > 0
}

func (a APIRequestReponse) HasResponseHeaders() bool {
	return len(a.Response.Headers) > 0
}

func (a APIRequestReponse) HasRequestBody() bool {
	return len(a.Request.Body) > 0
}

func (a APIRequestReponse) HasResponseBody() bool {
	return len(a.Response.Body) > 0
}

func (a APIRequestReponse) GetRequestBody() string {
	if a.HasRequestBody() {
		var test map[string]interface{}
		err := json.Unmarshal([]byte(a.Request.Body), &test)
		if err == nil {
			c, err2 := json.MarshalIndent(test, "", "    ")
			if err2 == nil {
				return string(c)
			}
		}
	}
	return a.Request.Body
}

func (a APIRequestReponse) GetResponseBody() string {
	if a.HasResponseBody() {
		var test map[string]interface{}
		err := json.Unmarshal([]byte(a.Response.Body), &test)
		if err == nil {
			c, err2 := json.MarshalIndent(test, "", "    ")
			if err2 == nil {
				return string(c)
			}
		}
	}
	return a.Response.Body
}

// sorting

// ByPath ...
func ByPath(p1, p2 *APIEndPoint) bool {
	return p1.ParameterizedPath < p2.ParameterizedPath
}

// ByAPIEndPoint ...
type ByAPIEndPoint func(p1, p2 *APIEndPoint) bool

// Sort ...
func (by ByAPIEndPoint) Sort(list []APIEndPoint) {
	ps := &APIEndPointSorter{
		endpoints: list,
		by:        by,
	}
	sort.Sort(ps)
}

// APIEndPointSorter ...
type APIEndPointSorter struct {
	endpoints []APIEndPoint
	by        func(p1, p2 *APIEndPoint) bool // Closure used in the Less method.
}

// Len  ...
func (s *APIEndPointSorter) Len() int {
	return len(s.endpoints)
}

// Swap  ...
func (s *APIEndPointSorter) Swap(i, j int) {
	s.endpoints[i], s.endpoints[j] = s.endpoints[j], s.endpoints[i]
}

// Less  ...
func (s *APIEndPointSorter) Less(i, j int) bool {
	return s.by(&s.endpoints[i], &s.endpoints[j])
}

// ByTitle ...
func ByTitle(p1, p2 *APIRequestReponse) bool {
	return p1.Title() < p2.Title()
}

// ByAPIRequestReponse ...
type ByAPIRequestReponse func(p1, p2 *APIRequestReponse) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by ByAPIRequestReponse) Sort(list []APIRequestReponse) {
	ps := &APIRequestReponseSorter{
		reqres: list,
		by:     by,
	}
	sort.Sort(ps)
}

// APIRequestReponseSorter ...
type APIRequestReponseSorter struct {
	reqres []APIRequestReponse
	by     func(p1, p2 *APIRequestReponse) bool // Closure used in the Less method.
}

// Len  ...
func (s *APIRequestReponseSorter) Len() int {
	return len(s.reqres)
}

// Swap  ...
func (s *APIRequestReponseSorter) Swap(i, j int) {
	s.reqres[i], s.reqres[j] = s.reqres[j], s.reqres[i]
}

// Less  ...
func (s *APIRequestReponseSorter) Less(i, j int) bool {
	return s.by(&s.reqres[i], &s.reqres[j])
}
