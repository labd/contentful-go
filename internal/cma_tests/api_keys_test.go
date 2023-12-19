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

func TestAPIKeyService_List(t *testing.T) {
	//var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "api_key.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/api_keys", r.URL.Path)
	})

	defer ts.Close()

	key, err := cma.WithSpaceId(testutil.SpaceID).ApiKeys().List(context.Background()).Next()
	assertions.Nil(err)
	assertions.Len(key.Items, 1)
	assertions.Equal("exampleapikey", key.Items[0].Sys.ID)
}

func TestAPIKeyService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "api_key_1.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/api_keys/exampleapikey", r.URL.Path)
	})

	defer ts.Close()

	key, err := cma.WithSpaceId(testutil.SpaceID).ApiKeys().Get(context.Background(), "exampleapikey")
	assertions.Nil(err)
	assertions.Equal("Example API Key", key.Name)
	assertions.Equal("b4c0n73n7fu1", key.AccessToken)
	assertions.Equal("1Mx3FqXX5XCJDtNpVW4BZI", key.PreviewAPIKey.Sys.ID)
}

func TestAPIKeyService_GetNotFound(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 404, Path: "/api_key/not_found.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/api_keys/exampleapikey", r.URL.Path)
	})

	defer ts.Close()

	_, err = cma.WithSpaceId(testutil.SpaceID).ApiKeys().Get(context.Background(), "exampleapikey")
	assertions.NotNil(err)
	var contentfulError common.NotFoundError
	assertions.True(errors.As(err, &contentfulError))
}

func TestAPIKeyService_Upsert_Create(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 201, Path: "api_key_1.json"}, nil, func(r *http.Request) {
		assertions.Equal("POST", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/api_keys", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("Example API Key", payload["name"])

	})

	defer ts.Close()

	key := &model.APIKey{
		BaseAPIKey: model.BaseAPIKey{
			Name: "Example API Key",
			Environments: []model.Environments{
				{
					Sys: model.Sys{
						ID:       "master",
						Type:     "Link",
						LinkType: "Environment",
					},
				},
			},
		},
		PreviewAPIKey: &model.BaseAPIKey{
			Sys: &model.Sys{
				ID:       "1Mx3FqXX5XCJDtNpVW4BZI",
				Type:     "Link",
				LinkType: "PreviewApiKey",
			},
		},
	}

	err := cma.WithSpaceId(testutil.SpaceID).ApiKeys().Upsert(context.Background(), key)
	assertions.Nil(err)
	assertions.Equal("exampleapikey", key.Sys.ID)
	assertions.Equal("Example API Key", key.Name)
}

func TestAPIKeyService_Upsert_Update(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "api_key_updated.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/api_keys/exampleapikey", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("This name is updated", payload["name"])

	})

	defer ts.Close()

	var key *model.APIKey
	err := testutil.ModelFromTestData("api_key_1.json", &key)
	assertions.Nil(err)

	key.Name = "This name is updated"

	err = cma.WithSpaceId(testutil.SpaceID).ApiKeys().Upsert(context.Background(), key)

	assertions.Nil(err)
	assertions.Equal("This name is updated", key.Name)
}

func TestAPIKeyService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: ""}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/api_keys/exampleapikey", r.URL.Path)
	})

	defer ts.Close()

	var key *model.APIKey
	err = testutil.ModelFromTestData("api_key_1.json", &key)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).ApiKeys().Delete(context.Background(), key)
	assertions.Nil(err)
}
