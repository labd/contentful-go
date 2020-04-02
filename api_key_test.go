package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIKeyService_List(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/api_keys")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("api_key.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	res, err := cma.APIKeys.List(spaceID).Next()
	assertions.Nil(err)
	keys := res.ToAPIKey()
	assertions.Equal(1, len(keys))
	assertions.Equal("exampleapikey", keys[0].Sys.ID)
}

func TestAPIKeyService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/api_keys/exampleapikey")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("api_key_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	key, err := cma.APIKeys.Get(spaceID, "exampleapikey")
	assertions.Nil(err)
	assertions.Equal("exampleapikey", key.Sys.ID)
}

func TestAPIKeyService_Get_2(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/api_keys/exampleapikey")

		checkHeaders(r, assertions)

		w.WriteHeader(400)
		_, _ = fmt.Fprintln(w, readTestData("api_key_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.APIKeys.Get(spaceID, "exampleapikey")
	assertions.NotNil(err)
}

func TestAPIKeyService_Upsert_Create(t *testing.T) {
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/api_keys")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("Example API Key", payload["name"])

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("api_key_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	key := &APIKey{
		Name:        "Example API Key",
		AccessToken: "b4c0n73n7fu1",
		Environments: []Environments{
			{
				Sys: Sys{
					ID:       "master",
					Type:     "Link",
					LinkType: "Environment",
				},
			},
		},
		PreviewAPIKey: PreviewAPIKey{
			Sys: Sys{
				ID:       "1Mx3FqXX5XCJDtNpVW4BZI",
				Type:     "Link",
				LinkType: "PreviewApiKey",
			},
		},
	}

	err := cma.APIKeys.Upsert(spaceID, key)
	assertions.Nil(err)
	assertions.Equal("exampleapikey", key.Sys.ID)
	assertions.Equal("Example API Key", key.Name)
}

func TestAPIKeyService_Upsert_Update(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/api_keys/exampleapikey")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("This name is updated", payload["name"])

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("api_key_updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	key, err := apiKeyFromTestData("api_key_1.json")
	assertions.Nil(err)

	key.Name = "This name is updated"

	err = cma.APIKeys.Upsert(spaceID, key)
	assertions.Nil(err)
	assertions.Equal("This name is updated", key.Name)
}

func TestAPIKeyService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/api_keys/exampleapikey")
		checkHeaders(r, assertions)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test locale
	key, err := apiKeyFromTestData("api_key_1.json")
	assertions.Nil(err)

	// delete locale
	err = cma.APIKeys.Delete(spaceID, key)
	assertions.Nil(err)
}
