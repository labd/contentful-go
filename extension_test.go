package contentful

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labd/contentful-go/pkgs/common"

	"github.com/stretchr/testify/assert"
)

func TestExtensionsService_List(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/extensions")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("extension.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	collection, err := cmaClient.Extensions.List(spaceID).Next()
	assertions.Nil(err)

	extensions := collection.ToExtension()
	assertions.Equal(1, len(extensions))
	assertions.Equal("My awesome extension", extensions[0].Extension.Name)
	assertions.Equal("https://example.com/my", extensions[0].Extension.SRC)
}

func TestExtensionsService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/extensions/0xvkPW9FdQ1kkWlWZ8ga4x")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("extension_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	extension, err := cmaClient.Extensions.Get(spaceID, "0xvkPW9FdQ1kkWlWZ8ga4x")
	assertions.Nil(err)
	assertions.Equal("0xvkPW9FdQ1kkWlWZ8ga4x", extension.Sys.ID)
}

func TestExtensionsService_Get_2(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/extensions/0xvkPW9FdQ1kkWlWZ8ga4x")

		checkHeaders(r, assertions)

		w.WriteHeader(400)
		_, _ = fmt.Fprintln(w, readTestData("extension_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	_, err = cmaClient.Extensions.Get(spaceID, "0xvkPW9FdQ1kkWlWZ8ga4x")
	assertions.NotNil(err)
	var contentfulError common.ErrorResponse
	assertions.True(errors.As(err, &contentfulError))
}

func TestExtensionsService_Upsert_Create(t *testing.T) {
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/master/extensions")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		extension := payload["extension"].(map[string]interface{})
		assertions.Equal("https://example.com/my", extension["src"])
		assertions.Equal("My awesome extension", extension["name"])

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("extension_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	extension := &Extension{
		Extension: ExtensionDetails{
			SRC:  "https://example.com/my",
			Name: "My awesome extension",
			FieldTypes: []FieldTypes{
				{
					Type: "Symbol",
				},
				{
					Type: "Text",
				},
			},
			Sidebar: false,
		},
	}

	err := cmaClient.Extensions.Upsert(spaceID, extension)
	assertions.Nil(err)
	assertions.Equal("https://example.com/my", extension.Extension.SRC)
	assertions.Equal("My awesome extension", extension.Extension.Name)
}

func TestExtensionsService_Upsert_Update(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/master/extensions/0xvkPW9FdQ1kkWlWZ8ga4x")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		extension := payload["extension"].(map[string]interface{})
		assertions.Equal("https://example.com/my", extension["src"])
		assertions.Equal("The updated extension", extension["name"])

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("extension_updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	extension, err := extensionFromTestFile("extension_1.json")
	assertions.Nil(err)

	extension.Extension.Name = "The updated extension"

	err = cmaClient.Extensions.Upsert(spaceID, extension)
	assertions.Nil(err)
	assertions.Equal("The updated extension", extension.Extension.Name)
}

func TestExtensionsService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/master/extensions/0xvkPW9FdQ1kkWlWZ8ga4x")
		checkHeaders(r, assertions)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	extension, err := extensionFromTestFile("extension_1.json")
	assertions.Nil(err)

	err = cmaClient.Extensions.Delete(spaceID, extension.Sys.ID)
	assertions.Nil(err)
}
