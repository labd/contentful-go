package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentAliasesServicesList(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environment_aliases")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("environment-alias.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.EnvironmentAliases.List(spaceID).Next()
	assertions.Nil(err)
}

func TestEnvironmentAliasesServicesGet(t *testing.T) {
	var err error
	assertions := assert.New(t)

	// Only tests master environment, as this is the only environment that always exists.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environment_aliases/master")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("environment-alias_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.EnvironmentAliases.Get(spaceID, "master")
	assertions.Nil(err)
}

func TestEnvironmentAliasesServiceUpdate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environment_aliases/master")

		checkHeaders(r, assertions)

		var payload EnvironmentAlias
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("staging", payload.Alias.Sys.ID)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, string(readTestData("environment-alias_1.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	environmentAlias, err := environmentAliasFromTestData("environment-alias_1.json")
	assertions.Nil(err)

	environmentAlias.Alias.Sys.ID = "staging"

	err = cma.EnvironmentAliases.Update(spaceID, environmentAlias)
	assertions.Nil(err)
}
