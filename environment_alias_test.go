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
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/environment_aliases")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("environment-alias.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.EnvironmentAliases.List(spaceID).Next()
	assert.Nil(err)
}

func TestEnvironmentAliasesServicesGet(t *testing.T) {
	var err error
	assert := assert.New(t)

	// Only tests master environment, as this is the only environment that always exists.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/environment_aliases/master")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("environment-alias_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.EnvironmentAliases.Get(spaceID, "master")
	assert.Nil(err)
}

func TestEnvironmentAliasesServiceUpdate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/environment_aliases/master")

		checkHeaders(r, assert)

		var payload EnvironmentAlias
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal("staging", payload.Alias.Sys.ID)

		w.WriteHeader(200)
		fmt.Fprintln(w, string(readTestData("environment-alias_1.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	environmentAlias, err := environmentAliasFromTestData("environment-alias_1.json")
	assert.Nil(err)

	environmentAlias.Alias.Sys.ID = "staging"

	err = cma.EnvironmentAliases.Update(spaceID, environmentAlias)
	assert.Nil(err)
}
