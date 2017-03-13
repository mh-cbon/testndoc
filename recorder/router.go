package recorder

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mh-cbon/testndoc"
)

// MuxRouter ...
type MuxRouter struct {
	*mux.Router
	api *testndoc.APIRecorder
}

// NewMuxRouter ...
func NewMuxRouter(r *mux.Router, a *testndoc.APIRecorder) *MuxRouter {
	return &MuxRouter{Router: r, api: a}
}

// HandleFunc ...
func (r *MuxRouter) HandleFunc(path string, h func(http.ResponseWriter, *http.Request)) *mux.Route {
	return r.Router.HandleFunc(path, r.api.FuncHandler(h, mux.Vars))
}
