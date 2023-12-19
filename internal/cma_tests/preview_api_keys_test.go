package cma_tests

import (
	"context"
	"errors"
	"github.com/flaconi/contentful-go/internal/testutil"
	"github.com/flaconi/contentful-go/pkgs/common"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreviewAPIKeyService_List(t *testing.T) {
	//var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "preview_api_keys/list.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/preview_api_keys", r.URL.Path)
	})

	defer ts.Close()

	key, err := cma.WithSpaceId(testutil.SpaceID).PreviewApiKeys().List(context.Background()).Next()
	assertions.Nil(err)
	assertions.Len(key.Items, 1)
	assertions.Equal("5leTEFZhsglc4EOkrnMNQ2", key.Items[0].Sys.ID)
}

func TestPreviewAPIKeyService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "preview_api_keys/single.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/preview_api_keys/exampleapikey", r.URL.Path)
	})

	defer ts.Close()

	key, err := cma.WithSpaceId(testutil.SpaceID).PreviewApiKeys().Get(context.Background(), "exampleapikey")
	assertions.Nil(err)
	assertions.Equal("example", key.Name)
	assertions.Equal("-fXcUaVWIW4njBjoqG1pJbNz65lg7mqxbugGyZyESzQ", key.AccessToken)
}

func TestPreviewAPIKeyService_GetNotFound(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockClient(t, assertions, testutil.ResponseData{StatusCode: 404, Path: "/preview_api_keys/not_found.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/preview_api_keys/exampleapikey", r.URL.Path)
	})

	defer ts.Close()

	_, err = cma.WithSpaceId(testutil.SpaceID).PreviewApiKeys().Get(context.Background(), "exampleapikey")
	assertions.NotNil(err)
	var contentfulError common.NotFoundError
	assertions.True(errors.As(err, &contentfulError))
}
