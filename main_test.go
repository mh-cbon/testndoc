package testndoc_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/mh-cbon/testndoc"
	"github.com/mh-cbon/testndoc/md"
	"github.com/mh-cbon/testndoc/recorder"
)

// ExmapleGorillaMux ...
func ExmapleGorillaMux() {

	// Create a proxy router to listen the requests / responses.
	router := recorder.NewMuxRouter(mux.NewRouter(), testndoc.Recorder)

	// Setup regular route handlers.
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("response"))
	})

	// Create a regular request.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		panic(err)
	}

	// Create a regular ResponseRecorder.
	rr := httptest.NewRecorder()

	// Execute a regular request on the router.
	router.ServeHTTP(rr, req)

	// Later, generate the doc
	md.Export(testndoc.Recorder, "path/")
}
