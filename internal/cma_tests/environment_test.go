package cma_tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/labd/contentful-go/internal/testutil"
	"github.com/labd/contentful-go/pkgs/common"
	"github.com/labd/contentful-go/pkgs/model"
	"github.com/labd/contentful-go/pkgs/util"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentService_List(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/environment/list.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments", r.URL.Path)
	})

	defer ts.Close()

	key, err := cma.WithSpaceId(testutil.SpaceID).Environments().List(context.Background()).Next()
	assertions.Nil(err)
	assertions.Len(key.Items, 1)
	assertions.Equal("master", key.Items[0].Sys.ID)
}

func TestEnvironmentService_Get(t *testing.T) {
	assertions := assert.New(t)

	var tests = []struct {
		resultValidation func(assertions *assert.Assertions, environment *model.Environment, err error)
		path             string
		name             string
		statusCode       int
	}{
		{
			resultValidation: func(assertions *assert.Assertions, environment *model.Environment, err error) {
				assertions.Nil(err)
				assertions.Equal("staging", environment.Name)
			},
			path:       "/environment/get.json",
			name:       "found",
			statusCode: 200,
		},
		{
			resultValidation: func(assertions *assert.Assertions, environment *model.Environment, err error) {
				assertions.NotNil(err)
				var contentfulError common.NotFoundError
				assertions.True(errors.As(err, &contentfulError))
			},
			path:       "/environment/not_found.json",
			statusCode: 404,
			name:       "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: tt.statusCode, Path: tt.path}, nil, func(r *http.Request) {
				assertions.Equal("GET", r.Method)
				assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/exampleapikey", r.URL.Path)
			})

			defer ts.Close()

			environment, err := cma.WithSpaceId(testutil.SpaceID).Environments().Get(context.Background(), "exampleapikey")
			tt.resultValidation(assertions, environment, err)
		})
	}
}

func TestEnvironmentService_Upsert_Create(t *testing.T) {
	assertions := assert.New(t)

	tests := []struct {
		name      string
		header    string
		sourceEnv *string
	}{
		{
			name:      "new",
			header:    "",
			sourceEnv: nil,
		},
		{
			name:      "copy",
			header:    "other",
			sourceEnv: util.ToPointer("other"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 201, Path: "/environment/get.json"}, nil, func(r *http.Request) {
				assertions.Equal("PUT", r.Method)
				assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing", r.URL.Path)

				var payload map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&payload)
				assertions.Nil(err)
				assertions.Equal("testing", payload["name"])

				assertions.Equal(tt.header, r.Header.Get("X-Contentful-Source-Environment"))
			})

			defer ts.Close()

			environment := &model.Environment{
				Name: "testing",
			}

			err := cma.WithSpaceId(testutil.SpaceID).Environments().Upsert(context.Background(), environment, tt.sourceEnv)
			assertions.Nil(err)
			assertions.Equal("staging", environment.Sys.ID)
			assertions.Equal("staging", environment.Name)
		})
	}
}

func TestEnvironmentService_Upsert_Update(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/environment/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/staging", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("This name is updated", payload["name"])

	})

	defer ts.Close()

	var environment *model.Environment
	err := testutil.ModelFromTestData("/environment/get.json", &environment)
	assertions.Nil(err)

	environment.Name = "This name is updated"

	err = cma.WithSpaceId(testutil.SpaceID).Environments().Upsert(context.Background(), environment, nil)

	assertions.Nil(err)
}

func TestEnvironmentService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: ""}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/staging", r.URL.Path)
	})

	defer ts.Close()

	var environment *model.Environment
	err = testutil.ModelFromTestData("/environment/get.json", &environment)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).Environments().Delete(context.Background(), environment)
	assertions.Nil(err)
}
