package contentful

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnapshotsServiceListEntrySnapshots(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/entries/hfM9RCJIk0wIm06WkEOQY/snapshots")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("snapshot-entry.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Snapshots.ListEntrySnapshots(spaceID, "hfM9RCJIk0wIm06WkEOQY").Next()
	assert.Nil(err)
}

func TestEntriesServiceGetEntrySnapshot(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/entries/hfM9RCJIk0wIm06WkEOQY/snapshots/4FLrUHftHW3v2BLi9fzfjU")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("snapshot-entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Snapshots.GetEntrySnapshot(spaceID, "hfM9RCJIk0wIm06WkEOQY", "4FLrUHftHW3v2BLi9fzfjU")
	assert.Nil(err)
}

func TestSnapshotsServiceListContentTypeSnapshots(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/content_types/hfM9RCJIk0wIm06WkEOQY/snapshots")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("snapshot-content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Snapshots.ListContentTypeSnapshots(spaceID, "hfM9RCJIk0wIm06WkEOQY").Next()
	assert.Nil(err)
}

func TestEntriesServiceGetContentTypeSnapshot(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/content_types/hfM9RCJIk0wIm06WkEOQY/snapshots/4FLrUHftHW3v2BLi9fzfjU")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("snapshot-content_type_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Snapshots.GetContentTypeSnapshots(spaceID, "hfM9RCJIk0wIm06WkEOQY", "4FLrUHftHW3v2BLi9fzfjU")
	assert.Nil(err)
}
