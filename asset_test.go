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

	collection, err := cma.Assets.List(spaceID).Next()
	assertions.Nil(err)
	asset := collection.ToAsset()
	assertions.Equal(3, len(asset))
	assertions.Equal("hehehe", asset[0].Fields.Title["en-US"])
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

	collection, err := cma.Assets.ListPublished(spaceID).Next()
	assertions.Nil(err)
	asset := collection.ToAsset()
	assertions.Equal(3, len(asset))
	assertions.Equal("hehehe", asset[0].Fields.Title["en-US"])
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

	asset, err := cma.Assets.Get(spaceID, "1x0xpXu4pSGS4OukSyWGUK")
	assertions.Nil(err)
	assertions.Equal("hehehe", asset.Fields.Title["en-US"])
}

func TestAssetsService_Get_2(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/1x0xpXu4pSGS4OukSyWGUK")

		checkHeaders(r, assertions)

		w.WriteHeader(400)
		_, _ = fmt.Fprintln(w, readTestData("asset_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Assets.Get(spaceID, "1x0xpXu4pSGS4OukSyWGUK")
	assertions.NotNil(err)
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
		assertions.Equal("hehehe", title["en-US"])

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
		Locale: "en-US",
		Fields: &AssetFields{
			Title: map[string]string{
				"en-US": "hehehe",
				"de":    "hehehe-de",
			},
			Description: map[string]string{
				"en-US": "asdfasf",
				"de":    "asdfasf-de",
			},
			File: map[string]*File{
				"en-US": {
					FileName:    "doge.jpg",
					ContentType: "image/jpeg",
					URL:         "//images.contentful.com/cfexampleapi/1x0xpXu4pSGS4OukSyWGUK/cc1239c6385428ef26f4180190532818/doge.jpg",
					UploadURL:   "",
					Details: &FileDetails{
						Size: 522943,
						Image: &ImageFields{
							Width:  5800,
							Height: 4350,
						},
					},
				},
			},
		},
	}

	err := cma.Assets.Upsert(spaceID, asset)
	assertions.Nil(err)
	assertions.Equal("hehehe", asset.Fields.Title["en-US"])
	assertions.Equal("d3b8dad44e5066cfb805e2357469ee64.png", asset.Fields.File["en-US"].FileName)
}

func TestAssetsService_Upsert_Update(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/assets/3HNzx9gvJScKku4UmcekYw")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		fields := payload["fields"].(map[string]interface{})
		title := fields["title"].(map[string]interface{})
		description := fields["description"].(map[string]interface{})
		assertions.Equal("updated", title["en-US"])
		assertions.Equal("also updated", description["en-US"])

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

	asset.Fields.Title["en-US"] = "updated"
	asset.Fields.Description["en-US"] = "also updated"

	err = cma.Assets.Upsert(spaceID, asset)
	assertions.Nil(err)
	assertions.Equal("updated", asset.Fields.Title["en-US"])
	assertions.Equal("also updated", asset.Fields.Description["en-US"])
}

func TestAssetsService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/assets/3HNzx9gvJScKku4UmcekYw")
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
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/3HNzx9gvJScKku4UmcekYw/files//process")

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
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/3HNzx9gvJScKku4UmcekYw/published")

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
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/3HNzx9gvJScKku4UmcekYw/published")

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
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/3HNzx9gvJScKku4UmcekYw/archived")

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
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/assets/3HNzx9gvJScKku4UmcekYw/archived")

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
