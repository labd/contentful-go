package contentful

import (
	"encoding/json"
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

func TestAssetsService_Upsert_Create(t *testing.T) {
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/assets")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		fields := payload["fields"].(map[string]interface{})
		title := fields["title"].(map[string]interface{})
		assertions.Equal("Doge", title["en-US"])

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("asset_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	asset := &Asset{
		locale: "en-US",
		Fields: &FileFields{
			Title:       "Doge",
			Description: "nice picture",
			File: &File{
				Name:        "doge.jpg",
				ContentType: "image/jpeg",
				URL:         "//images.contentful.com/cfexampleapi/1x0xpXu4pSGS4OukSyWGUK/cc1239c6385428ef26f4180190532818/doge.jpg",
				UploadURL:   "",
				Detail: &FileDetail{
					Size: 522943,
					Image: &FileImage{
						Width:  5800,
						Height: 4350,
					},
				},
			},
		},
	}

	err := cma.Assets.Upsert(spaceID, asset)
	assertions.Nil(err)
	assertions.Equal("Doge", asset.Fields.Title)
	assertions.Equal("doge.jpg", asset.Fields.File.Name)
}

func TestAssetsService_Upsert_Update(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		fields := payload["fields"].(map[string]interface{})
		title := fields["title"].(map[string]interface{})
		description := fields["description"].(map[string]interface{})
		assertions.Equal("Updated", title["en-US"])
		assertions.Equal("Lorum Ipsum", description["en-US"])

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("asset_updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	asset, err := assetFromTestData("asset_1.json")
	assertions.Nil(err)

	asset.Fields.Title = "Updated"
	asset.Fields.Description = "Lorum Ipsum"
	asset.locale = "en-US"

	err = cma.Assets.Upsert(spaceID, asset)
	assertions.Nil(err)
	assertions.Equal("Updated", asset.Fields.Title)
	assertions.Equal("Lorum Ipsum", asset.Fields.Description)
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

	// test asset
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

	// test asset
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

	// test asset
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

	// test asset
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

	// test asset
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

	// test asset
	asset, err := assetFromTestData("asset_1.json")
	assertions.Nil(err)

	err = cma.Assets.Unarchive(spaceID, asset)
	assertions.Nil(err)
}
