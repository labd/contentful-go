package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntriesServiceList(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/entries")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("entry.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Entries.List(spaceID).Next()
	assert.Nil(err)
}

func TestEntriesServiceGet(t *testing.T) {
	var err error
	assert := assert.New(t)

	// Only tests master environment, as this is the only environment that always exists.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/entries/5KsDBWseXY6QegucYAoacS")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Entries.Get(spaceID, "5KsDBWseXY6QegucYAoacS")
	assert.Nil(err)
}

func TestEntriesServiceDelete(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/entries/4aGeQYgByqQFJtToAOh2JJ")
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
	entry, err := entryFromTestData("locale_1.json")
	assert.Nil(err)

	// delete locale
	err = cma.Entries.Delete(spaceID, entry.Sys.ID)
	assert.Nil(err)
}

func TestEntriesServiceUpsertCreate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "POST")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/entries")

		checkHeaders(r, assert)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)

		fields := payload["Fields"].(map[string]interface{})
		title := fields["title"].(map[string]interface{})
		body := fields["body"].(map[string]interface{})
		assert.Equal("Hello, World!", title["en-US"].(string))
		assert.Equal("Bacon is healthy!", body["en-US"].(string))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	entry := &Entry{
		locale: "en-US",
		Fields: map[string]interface{}{
			"title": map[string]interface{}{
				"en-US": "Hello, World!",
			},
			"body": map[string]interface{}{
				"en-US": "Bacon is healthy!",
			},
		},
	}

	err = cma.Entries.Upsert(spaceID, entry)
	assert.Nil(err)
}

func TestEntriessServiceUpsertUpdate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/entries/5KsDBWseXY6QegucYAoacS")

		checkHeaders(r, assert)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)

		fields := payload["Fields"].(map[string]interface{})
		title := fields["title"].(map[string]interface{})
		body := fields["body"].(map[string]interface{})
		assert.Equal("Hello, World!", title["en-US"].(string))
		assert.Equal("Edited text", body["en-US"].(string))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	entry, err := entryFromTestData("entry_1.json")
	assert.Nil(err)

	body := entry.Fields["body"].(map[string]interface{})
	body["en-US"] = "Edited text"

	err = cma.Entries.Upsert(spaceID, entry)
	assert.Nil(err)
}

func TestEntriesServicePublish(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/entries/5KsDBWseXY6QegucYAoacS/published")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	e, err := entryFromTestData("entry_1.json")
	assert.Nil(err)

	err = cma.Entries.Publish(spaceID, e)
	assert.Nil(err)
}

func TestEntriesServiceUnpublish(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/entries/5KsDBWseXY6QegucYAoacS/published")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	e, err := entryFromTestData("entry_1.json")
	assert.Nil(err)

	err = cma.Entries.Unpublish(spaceID, e)
	assert.Nil(err)
}

func TestEntriesServiceArchive(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/entries/5KsDBWseXY6QegucYAoacS/archived")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	e, err := entryFromTestData("entry_1.json")
	assert.Nil(err)

	err = cma.Entries.Archive(spaceID, e)
	assert.Nil(err)
}

func TestEntriesServiceUnarchive(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/entries/5KsDBWseXY6QegucYAoacS/archived")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	e, err := entryFromTestData("entry_1.json")
	assert.Nil(err)

	err = cma.Entries.Unarchive(spaceID, e)
	assert.Nil(err)
}
