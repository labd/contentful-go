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
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/api_key/list.json"}, nil, func(r *http.Request) {
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
	assertions := assert.New(t)

	var tests = []struct {
		resultValidation func(assertions *assert.Assertions, apiKey *model.APIKey, err error)
		path             string
		name             string
		statusCode       int
	}{
		{
			resultValidation: func(assertions *assert.Assertions, apiKey *model.APIKey, err error) {
				assertions.Nil(err)
				assertions.Equal("Example API Key", apiKey.Name)
				assertions.Equal("b4c0n73n7fu1", apiKey.AccessToken)
				assertions.Equal("1Mx3FqXX5XCJDtNpVW4BZI", apiKey.PreviewAPIKey.Sys.ID)
			},
			path:       "/api_key/get.json",
			name:       "found",
			statusCode: 200,
		},
		{
			resultValidation: func(assertions *assert.Assertions, apiKey *model.APIKey, err error) {
				assertions.NotNil(err)
				var contentfulError common.NotFoundError
				assertions.True(errors.As(err, &contentfulError))
			},
			path:       "/api_key/not_found.json",
			statusCode: 404,
			name:       "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: tt.statusCode, Path: tt.path}, nil, func(r *http.Request) {
				assertions.Equal("GET", r.Method)
				assertions.Equal("/spaces/"+testutil.SpaceID+"/api_keys/exampleapikey", r.URL.Path)
			})

			defer ts.Close()

			apiKey, err := cma.WithSpaceId(testutil.SpaceID).ApiKeys().Get(context.Background(), "exampleapikey")
			tt.resultValidation(assertions, apiKey, err)
		})
	}
}

func TestAPIKeyService_Upsert_Create(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 201, Path: "/api_key/get.json"}, nil, func(r *http.Request) {
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
					Sys: model.BaseSys{
						ID:       "master",
						Type:     "Link",
						LinkType: "Environment",
					},
				},
			},
		},
		PreviewAPIKey: &model.PreviewAPIKeySys{
			Sys: model.BaseSys{
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

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/api_key/update.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/api_keys/exampleapikey", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("This name is updated", payload["name"])

	})

	defer ts.Close()

	var key *model.APIKey
	err := testutil.ModelFromTestData("/api_key/get.json", &key)
	assertions.Nil(err)

	key.Name = "This name is updated"

	err = cma.WithSpaceId(testutil.SpaceID).ApiKeys().Upsert(context.Background(), key)

	assertions.Nil(err)
	assertions.Equal("This name is updated", key.Name)
}

func TestAPIKeyService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: ""}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/api_keys/exampleapikey", r.URL.Path)
	})

	defer ts.Close()

	var key *model.APIKey
	err = testutil.ModelFromTestData("/api_key/get.json", &key)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).ApiKeys().Delete(context.Background(), key)
	assertions.Nil(err)
}
