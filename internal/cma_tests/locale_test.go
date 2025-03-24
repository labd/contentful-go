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

	"github.com/stretchr/testify/assert"
)

func TestLocaleService_List(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/locale/list.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/locales", r.URL.Path)
	})

	defer ts.Close()

	locales, err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Locales().List(context.Background()).Next()
	assertions.Nil(err)
	assertions.Len(locales.Items, 1)
	assertions.Equal("34N35DoyUQAtaKwWTgZs34", locales.Items[0].Sys.ID)
}

func TestLocaleService_Get(t *testing.T) {
	assertions := assert.New(t)

	var tests = []struct {
		resultValidation func(assertions *assert.Assertions, locale *model.Locale, err error)
		path             string
		name             string
		statusCode       int
	}{
		{
			resultValidation: func(assertions *assert.Assertions, locale *model.Locale, err error) {
				assertions.Nil(err)
				assertions.Equal("U.S. English", locale.Name)
				assertions.Equal("en-US", locale.Code)
			},
			path:       "/locale/get.json",
			name:       "found",
			statusCode: 200,
		},
		{
			resultValidation: func(assertions *assert.Assertions, locale *model.Locale, err error) {
				assertions.NotNil(err)
				var contentfulError common.NotFoundError
				assertions.True(errors.As(err, &contentfulError))
			},
			path:       "/locale/not_found.json",
			statusCode: 404,
			name:       "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: tt.statusCode, Path: tt.path}, nil, func(r *http.Request) {
				assertions.Equal("GET", r.Method)
				assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/locales/4aGeQYgByqQFJtToAOh2JJ", r.URL.Path)
			})

			defer ts.Close()

			locale, err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Locales().Get(context.Background(), "4aGeQYgByqQFJtToAOh2JJ")
			tt.resultValidation(assertions, locale, err)
		})
	}
}

func TestLocaleService_Upsert_Create(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 201, Path: "/locale/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("POST", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/locales", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("German (Austria)", payload["name"])
		assertions.Equal("de-AT", payload["code"])
	})

	defer ts.Close()

	locale := &model.Locale{
		Name: "German (Austria)",
		Code: "de-AT",
	}

	err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Locales().Upsert(context.Background(), locale)
	assertions.Nil(err)
	assertions.Equal("4aGeQYgByqQFJtToAOh2JJ", locale.Sys.ID)
	assertions.Equal("U.S. English", locale.Name)
}

func TestLocaleService_Upsert_Update(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/locale/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/locales/4aGeQYgByqQFJtToAOh2JJ", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("modified-name", payload["name"])
		assertions.Equal("modified-code", payload["code"])

	})

	defer ts.Close()

	var locale *model.Locale
	err := testutil.ModelFromTestData("/locale/get.json", &locale)
	assertions.Nil(err)

	locale.Name = "modified-name"
	locale.Code = "modified-code"

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Locales().Upsert(context.Background(), locale)

	assertions.Nil(err)
}

func TestLocaleService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: ""}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/locales/4aGeQYgByqQFJtToAOh2JJ", r.URL.Path)
	})

	defer ts.Close()

	var locale *model.Locale
	err = testutil.ModelFromTestData("/locale/get.json", &locale)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Locales().Delete(context.Background(), locale)
	assertions.Nil(err)
}
