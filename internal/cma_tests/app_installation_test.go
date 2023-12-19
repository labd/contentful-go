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

func TestAppInstallationService_List(t *testing.T) {
	//var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "app_installation.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/master/app_installations", r.URL.Path)
	})

	defer ts.Close()

	installations, err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("master").AppInstallations().List(context.Background()).Next()
	assertions.Nil(err)
	assertions.Len(installations.Items, 1)
	assertions.Equal("world", installations.Items[0].Parameters["hello"])
}

func TestAppInstallationService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "app_installation_1.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/master/app_installations/app_definition_id", r.URL.Path)
	})

	defer ts.Close()

	installation, err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("master").AppInstallations().Get(context.Background(), "app_definition_id")
	assertions.Nil(err)
	assertions.Equal("world", installation.Parameters["hello"])
}

func TestAppInstallationService_GetNotFound(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 404, Path: "/app_installation/not_found.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/master/app_installations/app_definition_id", r.URL.Path)
	})

	defer ts.Close()

	_, err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("master").AppInstallations().Get(context.Background(), "app_definition_id")
	assertions.NotNil(err)
	var contentfulError common.NotFoundError
	assertions.True(errors.As(err, &contentfulError))
}

func TestAppInstallationService_Upsert(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "app_installation_updated.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/master/app_installations/<app_definition_id>", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		parameters := payload["parameters"].(map[string]interface{})
		assertions.Equal("ipsum", parameters["lorum"])

	})

	defer ts.Close()

	var installation *model.AppInstallation
	err := testutil.ModelFromTestData("app_installation_1.json", &installation)
	assertions.Nil(err)

	installation.Parameters["lorum"] = "ipsum"

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("master").AppInstallations().Upsert(context.Background(), installation)
	assertions.Nil(err)
	assertions.Equal("<app_definition_id>", installation.Sys.AppDefinition.Sys.ID)
	assertions.Equal("ipsum", installation.Parameters["lorum"])
}

func TestAppInstallationService_Upsert_Forbidden(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 403, Path: "app_installation/forbidden.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/master/app_installations/<app_definition_id>", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		parameters := payload["parameters"].(map[string]interface{})
		assertions.Equal("ipsum", parameters["lorum"])

	})

	defer ts.Close()

	var installation *model.AppInstallation
	err := testutil.ModelFromTestData("app_installation_1.json", &installation)
	assertions.Nil(err)

	installation.Parameters["lorum"] = "ipsum"

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("master").AppInstallations().Upsert(context.Background(), installation)
	assertions.NotNil(err)
	var errorResponse common.ErrorResponse
	assertions.True(errors.As(err, &errorResponse))
	assertions.Equal("Expected X-Contentful-Marketplace to be \"i-accept-end-user-license-agreement,i-accept-marketplace-terms-of-service,i-accept-privacy-policy\". Visit https://app.contentful.com/deeplink?link=apps&id=66frtrAqmWSowDJzQNDiD to read Marketplace Terms of Service, EULA and Privacy Policy.", errorResponse.Details.Reasons)
	assertions.Equal("Forbidden", errorResponse.Message)
}

func TestAppInstallationService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: ""}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/master/app_installations/<app_definition_id>", r.URL.Path)
	})

	defer ts.Close()

	var installation *model.AppInstallation
	err = testutil.ModelFromTestData("app_installation_1.json", &installation)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("master").AppInstallations().Delete(context.Background(), installation)
	assertions.Nil(err)
}
