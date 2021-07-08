package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntriesService_List(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("entry.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.Entries.List(env).Next()
	assertions.Nil(err)
	entry := collection.ToEntry()
	assertions.Equal("5KsDBWseXY6QegucYAoacS", entry[0].Sys.ID)
}

func TestEntriesService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries/5KsDBWseXY6QegucYAoacS")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	entry, err := cma.Entries.Get(env, "5KsDBWseXY6QegucYAoacS")
	assertions.Nil(err)
	assertions.Equal("5KsDBWseXY6QegucYAoacS", entry.Sys.ID)
}

func TestEntriesService_Get_2(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries/5KsDBWseXY6QegucYAoacS")

		checkHeaders(r, assertions)

		w.WriteHeader(400)
		_, _ = fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Entries.Get(env, "5KsDBWseXY6QegucYAoacS")
	assertions.Nil(err)
}

func TestEntriesService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries/4aGeQYgByqQFJtToAOh2JJ")
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
	entry, err := entryFromTestData("locale_1.json")
	assertions.Nil(err)

	// delete locale
	err = cma.Entries.Delete(env, entry.Sys.ID)
	assertions.Nil(err)
}

func TestEntriesService_Upsert_Create(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries")

		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)

		fields := payload["fields"].(map[string]interface{})
		title := fields["title"].(map[string]interface{})
		body := fields["body"].(map[string]interface{})
		assertions.Equal("Hello, World!", title["en-US"].(string))
		assertions.Equal("Bacon is healthy!", body["en-US"].(string))

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	entry := &Entry{
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

	err = cma.Entries.Upsert(env, "hfM9RCJIk0wIm06WkEOQY", entry)
	assertions.Nil(err)
}

func TestEntriesService_Upsert_Update(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries/5KsDBWseXY6QegucYAoacS")

		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)

		fields := payload["fields"].(map[string]interface{})
		title := fields["title"].(map[string]interface{})
		body := fields["body"].(map[string]interface{})
		assertions.Equal("Hello, World!", title["en-US"].(string))
		assertions.Equal("Edited text", body["en-US"].(string))

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	entry, err := entryFromTestData("entry_1.json")
	assertions.Nil(err)

	body := entry.Fields["body"].(map[string]interface{})
	body["en-US"] = "Edited text"

	err = cma.Entries.Upsert(env, "hfM9RCJIk0wIm06WkEOQY", entry)
	assertions.Nil(err)
}

func TestEntriesService_Publish(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries/5KsDBWseXY6QegucYAoacS/published")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	e, err := entryFromTestData("entry_1.json")
	assertions.Nil(err)

	err = cma.Entries.Publish(env, e)
	assertions.Nil(err)
}

func TestEntriesService_Unpublish(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries/5KsDBWseXY6QegucYAoacS/published")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	e, err := entryFromTestData("entry_1.json")
	assertions.Nil(err)

	err = cma.Entries.Unpublish(env, e)
	assertions.Nil(err)
}

func TestEntriesService_Archive(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries/5KsDBWseXY6QegucYAoacS/archived")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	e, err := entryFromTestData("entry_1.json")
	assertions.Nil(err)

	err = cma.Entries.Archive(env, e)
	assertions.Nil(err)
}

func TestEntriesService_Unarchive(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries/5KsDBWseXY6QegucYAoacS/archived")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	e, err := entryFromTestData("entry_1.json")
	assertions.Nil(err)

	err = cma.Entries.Unarchive(env, e)
	assertions.Nil(err)
}
