package contentful

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppDefinitionsService_List(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/organizations/organization_id/app_definitions")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("app_definition.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.AppDefinitions.List("organization_id").Next()
	assertions.Nil(err)

	definitions := collection.ToAppDefinition()
	assertions.Equal(1, len(definitions))
	assertions.Equal("app_definition_id", definitions[0].Sys.ID)
	assertions.Equal("https://example.com/app.html", *definitions[0].SRC)
}

func TestAppDefinitionsService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/organizations/organization_id/app_definitions/app_definition_id")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("app_definition_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	definition, err := cma.AppDefinitions.Get("organization_id", "app_definition_id")
	assertions.Nil(err)
	assertions.Equal("app_definition_id", definition.Sys.ID)
	assertions.Equal("Hello world!", definition.Name)
}

func TestAppDefinitionsService_Get_2(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/organizations/organization_id/app_definitions/app_definition_id")

		checkHeaders(r, assertions)

		w.WriteHeader(400)
		_, _ = fmt.Fprintln(w, readTestData("app_definition_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.AppDefinitions.Get("organization_id", "app_definition_id")
	assertions.NotNil(err)
	var contentfulError ErrorResponse
	assertions.True(errors.As(err, &contentfulError))
}

func TestAppDefinitionsService_Get_3(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/organizations/organization_id/app_definitions/app_definition_id")

		checkHeaders(r, assertions)

		w.WriteHeader(404)
		_, _ = fmt.Fprintln(w, readTestData("app_definition_3.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.AppDefinitions.Get("organization_id", "app_definition_id")
	assertions.NotNil(err)
	var contentfulError NotFoundError
	assertions.True(errors.As(err, &contentfulError))
}

func TestAppDefinitionsService_Upsert_Create(t *testing.T) {
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/organizations/organization_id/app_definitions")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("Hello world!", payload["name"])
		assertions.Equal("https://example.com/app.html", payload["src"])

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("app_definition_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	src := "https://example.com/app.html"

	definition := &AppDefinition{
		Name: "Hello world!",
		SRC:  &src,
		Locations: []Locations{
			{
				Location: "entry-sidebar",
			},
		},
	}

	err := cma.AppDefinitions.Upsert("organization_id", definition)
	assertions.Nil(err)
	assertions.Equal("app_definition_id", definition.Sys.ID)
	assertions.Equal("Hello world!", definition.Name)
}

func TestAppDefinitionsService_Upsert_Create_Error(t *testing.T) {
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/organizations/organization_id/app_definitions")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("Hello world!", payload["name"])
		assertions.Equal("https://example.com/app.html", payload["src"])

		w.WriteHeader(422)
		_, _ = fmt.Fprintln(w, readTestData("error_validation_failed.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	src := "https://example.com/app.html"

	definition := &AppDefinition{
		Name: "Hello world!",
		SRC:  &src,
		Locations: []Locations{
			{
				Location: "entry-sidebar",
			},
		},
	}

	err := cma.AppDefinitions.Upsert("organization_id", definition)
	assertions.NotNil(err)
	var contentfulError ValidationFailedError
	assertions.True(errors.As(err, &contentfulError))
	assertions.Equal("ValidationFailed", contentfulError.err.Sys.ID)
	assertions.Equal("Validation error", contentfulError.err.Message)
	assertions.Equal("regexp", contentfulError.err.Details.Errors[0].Name)
	assertions.Equal("Does not match /(^https://)|(^http://localhost(:[0-9]+)?(/|$))/", contentfulError.err.Details.Errors[0].Details)
	assertions.Equal("localhost", contentfulError.err.Details.Errors[0].Value)
	assertions.Equal("src", contentfulError.err.Details.Errors[0].Path.([]any)[0])
}

func TestAppDefinitionsService_Upsert_Update(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/organizations/organization_id/app_definitions/app_definition_id")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("Hello Pluto", payload["name"])
		assertions.Equal("https://example.com/hellopluto.html", payload["src"])

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("app_definition_updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	definition, err := appDefinitionFromTestFile("app_definition_1.json")
	assertions.Nil(err)

	src := "https://example.com/hellopluto.html"

	definition.Name = "Hello Pluto"
	definition.SRC = &src

	err = cma.AppDefinitions.Upsert("organization_id", definition)
	assertions.Nil(err)
	assertions.Equal("Hello Pluto", definition.Name)
	assertions.Equal("https://example.com/hellopluto.html", *definition.SRC)
}

func TestAppDefinitionsService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.RequestURI, "/organizations/organization_id/app_definitions/app_definition_id")
		checkHeaders(r, assertions)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	definition, err := appDefinitionFromTestFile("app_definition_1.json")
	assertions.Nil(err)

	err = cma.AppDefinitions.Delete("organization_id", definition.Sys.ID)
	assertions.Nil(err)
}
