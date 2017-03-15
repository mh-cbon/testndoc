# {{.Name}}

WIP.

{{template "badge/travis" .}}{{template "badge/appveyor" .}}{{template "badge/goreport" .}}{{template "badge/godoc" .}}

{{pkgdoc}}

Check a demo for Markdown [here][demo.md]

Credits goes to [test2doc](https://github.com/adams-sarah/test2doc)
for initial insipiration.

# TODOS

- add documentation generators
- add tests
- refactor code copied from test2doc
- implement noop url params extractor
- add doc to explain routers integration based on the existing one.

# Install

#### Go
{{template "go/install" .}}

# Supported routers

- [gorilla mux](github.com/gorilla/mux)

# Test suite integration

To integrate the documentation generator to your tests,

Add a `TestMain` function, [see the doc](https://golang.org/pkg/testing/#hdr-Main),
to export the documentation once the tests finished.

```go
func TestMain(m *testing.M) {
  exitCode := m.Run()                            // run the tests
  err := md.Export(testndoc.Recorder, "path/")   // export the documentation
  if err != nil {
    panic(err)
  }
  os.Exit(exitCode)
}
```

In your API test handlers, make use of a RouterProxy,

```go

import (
  "github.com/gorilla/mux"
  "github.com/mh-cbon/testndoc/recorder"
)

// declare a router instance proxy, give it the proxied router, and an api recorder.
var router = recorder.NewMuxRouter(mux.NewRouter(), testndoc.Recorder)

// setup a regular test function
func TestMyEpApi(t *testing.T) {

  	// Somewhere, setup a regular router.
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

    // execute the test assertions
    // ....
}
```

# Usage

### All in one example

{{file "main_test.go"}}

# Recipes

#### Release the project

```sh
gump patch -d # check
gump patch # bump
```
