package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalesServiceList(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/locales")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("locales.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Locales.List(spaceID).Next()
	assertions.Nil(err)
}

func TestLocalesServiceGet(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/locales/4aGeQYgByqQFJtToAOh2JJ")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("locale_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	locale, err := cma.Locales.Get(spaceID, "4aGeQYgByqQFJtToAOh2JJ")
	assertions.Nil(err)
	assertions.Equal("U.S. English", locale.Name)
	assertions.Equal("en-US", locale.Code)
}

func TestLocalesServiceUpsertCreate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/locales")

		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("German (Austria)", payload["name"])
		assertions.Equal("de-AT", payload["code"])

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("locale_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	locale := &Locale{
		Name: "German (Austria)",
		Code: "de-AT",
	}

	err = cma.Locales.Upsert(spaceID, locale)
	assertions.Nil(err)
}

func TestLocalesServiceUpsertUpdate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/locales/4aGeQYgByqQFJtToAOh2JJ")

		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("modified-name", payload["name"])
		assertions.Equal("modified-code", payload["code"])

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, string(readTestData("locale_1.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	locale, err := localeFromTestData("locale_1.json")
	assertions.Nil(err)

	locale.Name = "modified-name"
	locale.Code = "modified-code"

	err = cma.Locales.Upsert(spaceID, locale)
	assertions.Nil(err)
}

func TestLocalesServiceDelete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/locales/4aGeQYgByqQFJtToAOh2JJ")
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
	locale, err := localeFromTestData("locale_1.json")
	assertions.Nil(err)

	// delete locale
	err = cma.Locales.Delete(spaceID, locale)
	assertions.Nil(err)
}
