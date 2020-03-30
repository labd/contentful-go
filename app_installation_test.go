package contentful

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppInstallationsService_List(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/app_installations")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("app_installation.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.AppInstallations.List(spaceID).Next()
	assertions.Nil(err)

	installation := collection.ToAppInstallation()
	assertions.Equal(1, len(installation))
	assertions.Equal("world", installation[0].Parameters["hello"])
}

func TestAppInstallationsService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/app_installations/app_definition_id")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("app_installation_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	installation, err := cma.AppInstallations.Get(spaceID, "app_definition_id")
	assertions.Nil(err)
	assertions.Equal("world", installation.Parameters["hello"])
}
