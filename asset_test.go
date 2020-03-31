package contentful

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssetsService_List(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("asset.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Assets.List(spaceID).Next()
	assertions.Nil(err)
}

func TestAssetsService_ListPublished(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/public/assets")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("asset.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Assets.ListPublished(spaceID).Next()
	assertions.Nil(err)
}

func TestAssetsService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("asset_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Assets.Get(spaceID, "1x0xpXu4pSGS4OukSyWGUK")
	assertions.Nil(err)
}

func TestAssetsService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK")
		checkHeaders(r, assertions)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test locale
	asset, err := assetFromTestData("asset_1.json")
	assertions.Nil(err)

	// delete locale
	err = cma.Assets.Delete(spaceID, asset)
	assertions.Nil(err)
}

func TestAssetsService_Process(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK/files//process")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("asset_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	asset, err := assetFromTestData("asset_1.json")
	assertions.Nil(err)

	err = cma.Assets.Process(spaceID, asset)
	assertions.Nil(err)
}

func TestAssetsService_Publish(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK/published")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("asset_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	asset, err := assetFromTestData("asset_1.json")
	assertions.Nil(err)

	err = cma.Assets.Publish(spaceID, asset)
	assertions.Nil(err)
}

func TestContentTypesService_Unpublish(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK/published")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("asset_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	asset, err := assetFromTestData("asset_1.json")
	assertions.Nil(err)

	err = cma.Assets.Unpublish(spaceID, asset)
	assertions.Nil(err)
}

func TestAssetsService_Archive(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK/archived")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("asset_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	asset, err := assetFromTestData("asset_1.json")
	assertions.Nil(err)

	err = cma.Assets.Archive(spaceID, asset)
	assertions.Nil(err)
}

func TestContentTypesService_Unarchive(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK/archived")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("asset_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	asset, err := assetFromTestData("asset_1.json")
	assertions.Nil(err)

	err = cma.Assets.Unarchive(spaceID, asset)
	assertions.Nil(err)
}
