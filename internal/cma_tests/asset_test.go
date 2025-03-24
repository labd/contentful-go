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

func TestAssetService_List(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/asset/list.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/assets", r.URL.Path)
	})

	defer ts.Close()

	assets, err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Assets().List(context.Background()).Next()
	assertions.Nil(err)
	assertions.Len(assets.Items, 1)
	assertions.Equal("hehehe", assets.Items[0].Fields.Title["en-US"])
}

func TestAssetService_ListPublished(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/asset/list.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/public/assets", r.URL.Path)
	})

	defer ts.Close()

	assets, err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Assets().ListPublished(context.Background()).Next()
	assertions.Nil(err)
	assertions.Len(assets.Items, 1)
	assertions.Equal("hehehe", assets.Items[0].Fields.Title["en-US"])
}

func TestAssetService_Get(t *testing.T) {
	assertions := assert.New(t)

	var tests = []struct {
		resultValidation func(assertions *assert.Assertions, asset *model.Asset, err error)
		path             string
		name             string
		statusCode       int
	}{
		{
			resultValidation: func(assertions *assert.Assertions, asset *model.Asset, err error) {
				assertions.Nil(err)
				assertions.Equal("hehehe", asset.Fields.Title["en-US"])
			},
			path:       "/asset/get.json",
			name:       "found",
			statusCode: 200,
		},
		{
			resultValidation: func(assertions *assert.Assertions, asset *model.Asset, err error) {
				assertions.NotNil(err)
				var contentfulError common.NotFoundError
				assertions.True(errors.As(err, &contentfulError))
			},
			path:       "/asset/not_found.json",
			statusCode: 404,
			name:       "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: tt.statusCode, Path: tt.path}, nil, func(r *http.Request) {
				assertions.Equal("GET", r.Method)
				assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/assets/1x0xpXu4pSGS4OukSyWGUK", r.URL.Path)
			})

			defer ts.Close()

			asset, err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Assets().Get(context.Background(), "1x0xpXu4pSGS4OukSyWGUK")
			tt.resultValidation(assertions, asset, err)
		})
	}
}

func TestAssetService_Upsert_Create(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 201, Path: "/asset/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("POST", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/assets", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		fields := payload["fields"].(map[string]interface{})
		title := fields["title"].(map[string]interface{})
		assertions.Equal("hehehe", title["en-US"])
	})

	defer ts.Close()

	asset := &model.Asset{
		Fields: &model.AssetFields{
			Title: map[string]string{
				"en-US": "hehehe",
				"de":    "hehehe-de",
			},
			Description: map[string]string{
				"en-US": "asdfasf",
				"de":    "asdfasf-de",
			},
			File: map[string]*model.File{
				"en-US": {
					FileName:    "doge.jpg",
					ContentType: "image/jpeg",
					URL:         "//images.contentful.com/cfexampleapi/1x0xpXu4pSGS4OukSyWGUK/cc1239c6385428ef26f4180190532818/doge.jpg",
					UploadURL:   "",
					Details: &model.FileDetails{
						Size: 522943,
						Image: &model.ImageFields{
							Width:  5800,
							Height: 4350,
						},
					},
				},
			},
		},
	}

	err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Assets().Upsert(context.Background(), asset)
	assertions.Nil(err)
	assertions.Equal("3HNzx9gvJScKku4UmcekYw", asset.Sys.ID)
	assertions.Equal("hehehe", asset.Fields.Title["en-US"])
	assertions.Equal("d3b8dad44e5066cfb805e2357469ee64.png", asset.Fields.File["en-US"].FileName)
}

func TestAssetService_Upsert_Update(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/asset/update.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/assets/3HNzx9gvJScKku4UmcekYw", r.URL.Path)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		fields := payload["fields"].(map[string]interface{})
		title := fields["title"].(map[string]interface{})
		description := fields["description"].(map[string]interface{})
		assertions.Equal("updated", title["en-US"])
		assertions.Equal("also updated", description["en-US"])

	})

	defer ts.Close()

	var asset *model.Asset
	err := testutil.ModelFromTestData("/asset/get.json", &asset)
	assertions.Nil(err)

	asset.Fields.Title["en-US"] = "updated"
	asset.Fields.Description["en-US"] = "also updated"

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Assets().Upsert(context.Background(), asset)
	assertions.Nil(err)
	assertions.Equal("updated", asset.Fields.Title["en-US"])
	assertions.Equal("also updated", asset.Fields.Description["en-US"])
}

func TestAssetService_Process(t *testing.T) {
	var err error
	assertions := assert.New(t)

	first := true

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 204, Path: ""}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		if first {
			assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/assets/3HNzx9gvJScKku4UmcekYw/files/en-US/process", r.URL.Path)
			first = false
			return
		}
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/assets/3HNzx9gvJScKku4UmcekYw/files/de/process", r.URL.Path)

	})

	defer ts.Close()

	var asset *model.Asset
	err = testutil.ModelFromTestData("/asset/get.json", &asset)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Assets().Process(context.Background(), asset)
	assertions.Nil(err)
}

func TestAssetService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 204, Path: ""}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/assets/3HNzx9gvJScKku4UmcekYw", r.URL.Path)
	})

	defer ts.Close()

	var asset *model.Asset
	err = testutil.ModelFromTestData("/asset/get.json", &asset)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Assets().Delete(context.Background(), asset)
	assertions.Nil(err)
}

func TestAssetService_Publish(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/asset/update.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/assets/3HNzx9gvJScKku4UmcekYw/published", r.URL.Path)
	})

	defer ts.Close()

	var asset *model.Asset
	err = testutil.ModelFromTestData("/asset/get.json", &asset)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Assets().Publish(context.Background(), asset)
	assertions.Nil(err)
}

func TestAssetService_Unpublish(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/asset/update.json"}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/assets/3HNzx9gvJScKku4UmcekYw/published", r.URL.Path)
	})

	defer ts.Close()

	var asset *model.Asset
	err = testutil.ModelFromTestData("/asset/get.json", &asset)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Assets().Unpublish(context.Background(), asset)
	assertions.Nil(err)
}

func TestAssetService_Archive(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/asset/update.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/assets/3HNzx9gvJScKku4UmcekYw/archived", r.URL.Path)
	})

	defer ts.Close()

	var asset *model.Asset
	err = testutil.ModelFromTestData("/asset/get.json", &asset)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Assets().Archive(context.Background(), asset)
	assertions.Nil(err)
}

func TestAssetService_Unarchive(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/asset/update.json"}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/testing/assets/3HNzx9gvJScKku4UmcekYw/archived", r.URL.Path)
	})

	defer ts.Close()

	var asset *model.Asset
	err = testutil.ModelFromTestData("/asset/get.json", &asset)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("testing").Assets().Unarchive(context.Background(), asset)
	assertions.Nil(err)
}
