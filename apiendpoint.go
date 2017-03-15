package testndoc

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
