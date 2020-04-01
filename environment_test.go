package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentsServiceList(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("environments.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.Environments.List(spaceID).Next()
	assertions.Nil(err)
	environment := collection.ToEnvironment()
	assertions.Equal(1, len(environment))
	assertions.Equal("master", environment[0].Name)
}

func TestEnvironmentsServiceGet(t *testing.T) {
	var err error
	assertions := assert.New(t)

	// Only tests master environment, as this is the only environment that always exists.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("environment.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Environments.Get(spaceID, "master")
	assertions.Nil(err)
}

func TestEnvironmentsServiceUpsertCreate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments")

		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("staging", payload["name"])

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("environment_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	environment := &Environment{
		Name: "staging",
	}

	err = cma.Environments.Upsert(spaceID, environment)
	assertions.Nil(err)
}

func TestEnvironmentsServiceUpsertUpdate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/staging")

		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("modified-name", payload["name"])

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, string(readTestData("environment_1.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	environment, err := environmentFromTestData("environment_1.json")
	assertions.Nil(err)

	environment.Name = "modified-name"

	err = cma.Environments.Upsert(spaceID, environment)
	assertions.Nil(err)
}

func TestEnvironmentsServiceDelete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/staging")
		checkHeaders(r, assertions)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test environment
	environment, err := environmentFromTestData("environment_1.json")
	assertions.Nil(err)

	// delete environment
	err = cma.Environments.Delete(spaceID, environment)
	assertions.Nil(err)
}
