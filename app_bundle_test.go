package contentful

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppBundleService_Create(t *testing.T) {
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/organizations/organization_id/app_definitions/definition_id/app_bundles")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		_, ok := payload["comment"]
		assertions.True(ok)

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("app_bundle/create.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	bundle, err := cma.AppBundle.Create("organization_id", "definition_id", "comment", "upload_id")
	assertions.Nil(err)
	assertions.Equal("app_bundle_id", bundle.Sys.ID)
	assertions.Equal("comment", bundle.Comment)
	assertions.Len(bundle.Files, 1)
}

func TestAppBundleService_Create_Error(t *testing.T) {
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/organizations/organization_id/app_definitions/definition_id/app_bundles")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		_, ok := payload["comment"]
		assertions.False(ok)

		w.WriteHeader(400)
		_, _ = fmt.Fprintln(w, readTestData("app_bundle/upload_error.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	bundle, err := cma.AppBundle.Create("organization_id", "definition_id", "", "upload_id")
	assertions.NotNil(err)
	assertions.Nil(bundle)
	assertions.Equal("Failed to find upload", err.Error())
}
