package contentful

import (
	"encoding/json"
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

func TestAppInstallationsService_Upsert_Create(t *testing.T) {
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/master/app_installations")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		parameters := payload["parameters"].(map[string]interface{})
		assertions.Equal("world", parameters["hello"])

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, string(readTestData("app_installation_1.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	installation := &AppInstallation{
		Parameters: map[string]string{
			"hello": "world",
		},
	}

	err := cma.AppInstallations.Upsert(spaceID, "", installation)
	assertions.Nil(err)
	assertions.Equal("world", installation.Parameters["hello"])
}

func TestAppInstallationsService_Upsert_Update(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/master/app_installations/app_definition_id")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		parameters := payload["parameters"].(map[string]interface{})
		assertions.Equal("ipsum", parameters["lorum"])

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, string(readTestData("app_installation_updated.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	installation, err := appInstallationFromTestFile("app_installation_1.json")
	assertions.Nil(err)

	installation.Parameters["lorum"] = "ipsum"

	err = cma.AppInstallations.Upsert(spaceID, "app_definition_id", installation)
	assertions.Nil(err)
	assertions.Equal("ipsum", installation.Parameters["lorum"])
}
