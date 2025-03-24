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

func TestEntryService_List(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/entry/list.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/entries", r.URL.Path)
	})

	defer ts.Close()

	entries, err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").Entries().List(context.Background()).Next()
	assertions.Nil(err)
	assertions.Len(entries.Items, 1)
	assertions.Equal("5KsDBWseXY6QegucYAoacS", entries.Items[0].Sys.ID)
	assertions.Equal("Hello, World!", entries.Items[0].Fields["title"].(map[string]interface{})["en-US"])
}

func TestEntryService_Get(t *testing.T) {
	assertions := assert.New(t)

	var tests = []struct {
		resultValidation func(assertions *assert.Assertions, ct *model.Entry, err error)
		path             string
		name             string
		statusCode       int
	}{
		{
			resultValidation: func(assertions *assert.Assertions, entry *model.Entry, err error) {
				assertions.Nil(err)
				assertions.Equal("Hello, World!", entry.Fields["title"].(map[string]interface{})["en-US"])
				assertions.Equal("5KsDBWseXY6QegucYAoacS", entry.Sys.ID)
			},
			path:       "/entry/get.json",
			name:       "found",
			statusCode: 200,
		},
		{
			resultValidation: func(assertions *assert.Assertions, ct *model.Entry, err error) {
				assertions.NotNil(err)
				var contentfulError common.NotFoundError
				assertions.True(errors.As(err, &contentfulError))
			},
			path:       "/entry/not_found.json",
			statusCode: 404,
			name:       "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: tt.statusCode, Path: tt.path}, nil, func(r *http.Request) {
				assertions.Equal("GET", r.Method)
				assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/entries/63Vgs0BFK0USe4i2mQUGK6", r.URL.Path)
			})

			defer ts.Close()

			entry, err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").Entries().Get(context.Background(), "63Vgs0BFK0USe4i2mQUGK6")
			tt.resultValidation(assertions, entry, err)
		})
	}
}

func TestEntryService_Upsert_Create(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 201, Path: "/entry/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("POST", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/entries", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		fields := payload["fields"].(map[string]interface{})
		title := fields["title"].(map[string]interface{})
		body := fields["body"].(map[string]interface{})
		assertions.Equal("Hello, World!", title["en-US"].(string))
		assertions.Equal("Bacon is healthy!", body["en-US"].(string))
	})

	defer ts.Close()

	entry := &model.Entry{
		Locale: "en-US",
		Fields: map[string]interface{}{
			"title": map[string]interface{}{
				"en-US": "Hello, World!",
			},
			"body": map[string]interface{}{
				"en-US": "Bacon is healthy!",
			},
		},
	}

	err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").Entries().Upsert(context.Background(), "testContentType", entry)
	assertions.Nil(err)
	assertions.Equal("5KsDBWseXY6QegucYAoacS", entry.Sys.ID)
	assertions.Equal("Hello, World!", entry.Fields["title"].(map[string]interface{})["en-US"])
}

func TestEntryService_Upsert_Update(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/entry/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/entries/5KsDBWseXY6QegucYAoacS", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		fields := payload["fields"].(map[string]interface{})
		title := fields["title"].(map[string]interface{})
		body := fields["body"].(map[string]interface{})
		assertions.Equal("Hello, World!", title["en-US"].(string))
		assertions.Equal("Bacon is healthy!", body["en-US"].(string))
	})

	defer ts.Close()

	var entry *model.Entry
	err := testutil.ModelFromTestData("/entry/get.json", &entry)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").Entries().Upsert(context.Background(), "testContentType", entry)
	assertions.Nil(err)
}

func TestEntryService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 204, Path: ""}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/entries/5KsDBWseXY6QegucYAoacS", r.URL.Path)
	})

	defer ts.Close()

	var entry *model.Entry
	err = testutil.ModelFromTestData("/entry/get.json", &entry)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").Entries().Delete(context.Background(), entry)
	assertions.Nil(err)
}

func TestEntryService_Publish(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/entry/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/entries/5KsDBWseXY6QegucYAoacS/published", r.URL.Path)
	})

	defer ts.Close()

	var entry *model.Entry
	err = testutil.ModelFromTestData("/entry/get.json", &entry)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").Entries().Publish(context.Background(), entry)
	assertions.Nil(err)
}

func TestEntryService_Unpublish(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/entry/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/entries/5KsDBWseXY6QegucYAoacS/published", r.URL.Path)
	})

	defer ts.Close()

	var entry *model.Entry
	err = testutil.ModelFromTestData("/entry/get.json", &entry)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").Entries().Unpublish(context.Background(), entry)
	assertions.Nil(err)
}

func TestEntryService_Archive(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/entry/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/entries/5KsDBWseXY6QegucYAoacS/archived", r.URL.Path)
	})

	defer ts.Close()

	var entry *model.Entry
	err = testutil.ModelFromTestData("/entry/get.json", &entry)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").Entries().Archive(context.Background(), entry)
	assertions.Nil(err)
}

func TestEntryService_Unarchive(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/entry/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/entries/5KsDBWseXY6QegucYAoacS/archived", r.URL.Path)
	})

	defer ts.Close()

	var entry *model.Entry
	err = testutil.ModelFromTestData("/entry/get.json", &entry)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").Entries().Unarchive(context.Background(), entry)
	assertions.Nil(err)
}
