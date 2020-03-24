package contentful

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssetsServiceList(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("spaces-id1-assets.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Assets.List(spaceID).Next()
	assert.Nil(err)
}

func TestAssetsServiceListPublished(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/public/assets")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("spaces-id1-assets.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Assets.ListPublished(spaceID).Next()
	assert.Nil(err)
}

func TestAssetsServiceGet(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("spaces-id1-assets-1x0xpXu4pSGS4OukSyWGUK.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Assets.Get(spaceID, "1x0xpXu4pSGS4OukSyWGUK")
	assert.Nil(err)
}

func TestAssetsServiceDelete(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK")
		checkHeaders(r, assert)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test locale
	asset, err := assetFromTestData("spaces-id1-assets-1x0xpXu4pSGS4OukSyWGUK.json")
	assert.Nil(err)

	// delete locale
	err = cma.Assets.Delete(spaceID, asset)
	assert.Nil(err)
}

func TestAssetsServiceProcess(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK/files//process")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("spaces-id1-assets-1x0xpXu4pSGS4OukSyWGUK.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	asset, err := assetFromTestData("spaces-id1-assets-1x0xpXu4pSGS4OukSyWGUK.json")
	assert.Nil(err)

	err = cma.Assets.Process(spaceID, asset)
	assert.Nil(err)
}

func TestAssetsServicePublish(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK/published")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("spaces-id1-assets-1x0xpXu4pSGS4OukSyWGUK.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	asset, err := assetFromTestData("spaces-id1-assets-1x0xpXu4pSGS4OukSyWGUK.json")
	assert.Nil(err)

	err = cma.Assets.Publish(spaceID, asset)
	assert.Nil(err)
}

func TestContentTypesServiceUnpublish(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK/published")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("spaces-id1-assets-1x0xpXu4pSGS4OukSyWGUK.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	asset, err := assetFromTestData("spaces-id1-assets-1x0xpXu4pSGS4OukSyWGUK.json")
	assert.Nil(err)

	err = cma.Assets.Unpublish(spaceID, asset)
	assert.Nil(err)
}

func TestAssetsServiceArchive(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK/archived")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("spaces-id1-assets-1x0xpXu4pSGS4OukSyWGUK.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	asset, err := assetFromTestData("spaces-id1-assets-1x0xpXu4pSGS4OukSyWGUK.json")
	assert.Nil(err)

	err = cma.Assets.Archive(spaceID, asset)
	assert.Nil(err)
}

func TestContentTypesServiceUnarchive(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK/archived")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("spaces-id1-assets-1x0xpXu4pSGS4OukSyWGUK.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	asset, err := assetFromTestData("spaces-id1-assets-1x0xpXu4pSGS4OukSyWGUK.json")
	assert.Nil(err)

	err = cma.Assets.Unarchive(spaceID, asset)
	assert.Nil(err)
}
