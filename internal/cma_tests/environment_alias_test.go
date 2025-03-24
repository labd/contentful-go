package cma_tests

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/labd/contentful-go/internal/testutil"
	"github.com/labd/contentful-go/pkgs/model"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentAliasService_List(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/environment_alias/list.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environment_aliases", r.URL.Path)
	})

	defer ts.Close()

	list, err := cma.WithSpaceId(testutil.SpaceID).EnvironmentAliases().List(context.Background()).Next()
	assertions.Nil(err)
	assertions.Len(list.Items, 1)
	assertions.Equal("master-18-3-2020", list.Items[0].Alias.Sys.ID)
}

func TestEnvironmentAliasService_Upsert_Update(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/environment_alias/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environment_aliases/master", r.URL.Path)

		var payload map[string]any
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("staging", payload["environment"].(map[string]any)["sys"].(map[string]any)["id"])

	})

	defer ts.Close()

	var environment *model.EnvironmentAlias
	err := testutil.ModelFromTestData("/environment_alias/get.json", &environment)
	assertions.Nil(err)

	environment.Alias.Sys.ID = "staging"

	err = cma.WithSpaceId(testutil.SpaceID).EnvironmentAliases().Upsert(context.Background(), environment)

	assertions.Nil(err)
}

func TestEnvironmentAliasService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: ""}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environment_aliases/master", r.URL.Path)
	})

	defer ts.Close()

	var environment *model.EnvironmentAlias
	err = testutil.ModelFromTestData("/environment_alias/get.json", &environment)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).EnvironmentAliases().Delete(context.Background(), environment)
	assertions.Nil(err)
}
