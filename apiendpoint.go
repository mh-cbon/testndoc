package testndoc

import (
	"encoding/json"
	"strings"
)

type APIDoc struct {
	EndPoints []APIEndPoint
}
type APIEndPoint struct {
	Path              string
	ParameterizedPath string
	HandlerFn         APIEndPointFn
	ReqRes            []APIRequestReponse
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

func (a APIRequestReponse) HasUrlParams() bool {
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
