package contentful

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentsServiceList(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("environments.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Environments.List(spaceID).Next()
	assert.Nil(err)
}

func TestEnvironmentsServiceGet(t *testing.T) {
	var err error
	assert := assert.New(t)

	// Only tests master environment, as this is the only environment that always exists.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("environment.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Environments.Get(spaceID, "master")
	assert.Nil(err)
}
