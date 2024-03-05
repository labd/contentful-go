package cma_tests

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/flaconi/contentful-go/internal/testutil"
	"github.com/flaconi/contentful-go/pkgs/common"
	"github.com/flaconi/contentful-go/pkgs/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppDefinitionService_List(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/app_definition/list.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/organizations/"+testutil.OrganizationId+"/app_definitions", r.URL.Path)
	})

	defer ts.Close()

	list, err := cma.WithOrganizationId(testutil.OrganizationId).AppDefinitions().List(context.Background()).Next()
	assertions.Nil(err)
	assertions.Len(list.Items, 1)
	assertions.Equal("app_definition_id", list.Items[0].Sys.ID)
	assertions.Equal("https://example.com/app.html", *list.Items[0].SRC)
}

func TestAppDefinitionService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/app_definition/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/organizations/"+testutil.OrganizationId+"/app_definitions/app_definition_id", r.URL.Path)
	})

	defer ts.Close()

	appDefinition, err := cma.WithOrganizationId(testutil.OrganizationId).AppDefinitions().Get(context.Background(), "app_definition_id")
	assertions.Nil(err)
	assertions.Equal("app_definition_id", appDefinition.Sys.ID)
	assertions.Equal("Hello world!", appDefinition.Name)
	assertions.Equal("https://example.com/app.html", *appDefinition.SRC)
	assertions.Equal("entry-sidebar", appDefinition.Locations[0].Location)
}

func TestAppDefinitionService_GetNotFound(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 404, Path: "/app_definition/not_found.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/organizations/"+testutil.OrganizationId+"/app_definitions/app_definition_id", r.URL.Path)
	})

	defer ts.Close()

	_, err = cma.WithOrganizationId(testutil.OrganizationId).AppDefinitions().Get(context.Background(), "app_definition_id")
	assertions.NotNil(err)
	var contentfulError common.NotFoundError
	assertions.True(errors.As(err, &contentfulError))
}

func TestAppDefinitionService_Upsert_Create(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 201, Path: "/app_definition/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("POST", r.Method)
		assertions.Equal("/organizations/"+testutil.OrganizationId+"/app_definitions", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("Hello world!", payload["name"])
		assertions.Equal("https://example.com/app.html", payload["src"])

	})

	defer ts.Close()

	src := "https://example.com/app.html"

	definition := &model.AppDefinition{
		Name: "Hello world!",
		SRC:  &src,
		Locations: []model.Locations{
			{
				Location: "entry-sidebar",
			},
		},
	}

	err := cma.WithOrganizationId(testutil.OrganizationId).AppDefinitions().Upsert(context.Background(), definition)
	assertions.Nil(err)
	assertions.Equal("app_definition_id", definition.Sys.ID)
	assertions.Equal("Hello world!", definition.Name)
}

func TestAppDefinitionService_Upsert_Create_ValidationError(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 422, Path: "/app_definition/error_validation_failed.json"}, nil, func(r *http.Request) {
		assertions.Equal("POST", r.Method)
		assertions.Equal("/organizations/"+testutil.OrganizationId+"/app_definitions", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("Hello world!", payload["name"])
		assertions.Equal("https://example.com/app.html", payload["src"])

	})

	defer ts.Close()

	src := "https://example.com/app.html"

	definition := &model.AppDefinition{
		Name: "Hello world!",
		SRC:  &src,
		Locations: []model.Locations{
			{
				Location: "entry-sidebar",
			},
		},
	}

	err := cma.WithOrganizationId(testutil.OrganizationId).AppDefinitions().Upsert(context.Background(), definition)
	assertions.NotNil(err)
	var contentfulError common.ValidationFailedError
	assertions.True(errors.As(err, &contentfulError))
	assertions.Equal("Value \"localhost\" in path \"src\" with details: \"Does not match /(^https://)|(^http://localhost(:[0-9]+)?(/|$))/\"\n", contentfulError.Error())
}

func TestAppDefinitionService_Upsert_Update(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/app_definition/update.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/organizations/"+testutil.OrganizationId+"/app_definitions/app_definition_id", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("Hello Pluto", payload["name"])
		assertions.Equal("https://example.com/hellopluto.html", payload["src"])
	})

	defer ts.Close()

	var definition *model.AppDefinition
	err := testutil.ModelFromTestData("/app_definition/get.json", &definition)
	assertions.Nil(err)

	src := "https://example.com/hellopluto.html"

	definition.Name = "Hello Pluto"
	definition.SRC = &src

	err = cma.WithOrganizationId(testutil.OrganizationId).AppDefinitions().Upsert(context.Background(), definition)

	assertions.Nil(err)
	assertions.Equal("Hello Pluto", definition.Name)
	assertions.Equal("https://example.com/hellopluto.html", *definition.SRC)
}

func TestAppDefinitionService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: ""}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/organizations/"+testutil.OrganizationId+"/app_definitions/app_definition_id", r.URL.Path)
	})

	defer ts.Close()

	var definition *model.AppDefinition
	err = testutil.ModelFromTestData("/app_definition/get.json", &definition)
	assertions.Nil(err)

	err = cma.WithOrganizationId(testutil.OrganizationId).AppDefinitions().Delete(context.Background(), definition)
	assertions.Nil(err)
}
