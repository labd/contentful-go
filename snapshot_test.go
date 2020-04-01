package contentful

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnapshotsService_ListEntrySnapshots(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/entries/hfM9RCJIk0wIm06WkEOQY/snapshots")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("snapshot_entry.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.Snapshots.ListEntrySnapshots(spaceID, "hfM9RCJIk0wIm06WkEOQY").Next()
	assertions.Nil(err)
	entrySnapshot := collection.ToEntrySnapshot()
	assertions.Equal(1, len(entrySnapshot))
	assertions.Equal("Hello, World!", entrySnapshot[0].EntrySnapshotDetail.Fields["title"].(map[string]interface{})["en-US"])
}

func TestEntriesServiceGetEntrySnapshot(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/entries/hfM9RCJIk0wIm06WkEOQY/snapshots/4FLrUHftHW3v2BLi9fzfjU")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("snapshot_entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	entrySnapshot, err := cma.Snapshots.GetEntrySnapshot(spaceID, "hfM9RCJIk0wIm06WkEOQY", "4FLrUHftHW3v2BLi9fzfjU")
	assertions.Nil(err)
	assertions.Equal("Hello, World!", entrySnapshot.EntrySnapshotDetail.Fields["title"].(map[string]interface{})["en-US"])
}

func TestSnapshotsService_ListContentTypeSnapshots(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/content_types/hfM9RCJIk0wIm06WkEOQY/snapshots")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("snapshot-content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.Snapshots.ListContentTypeSnapshots(spaceID, "hfM9RCJIk0wIm06WkEOQY").Next()
	assertions.Nil(err)
	entrySnapshot := collection.ToContentTypeSnapshot()
	assertions.Equal(1, len(entrySnapshot))
	assertions.Equal("Blog Post", entrySnapshot[0].ContentTypeSnapshotDetail.Name)
}

func TestEntriesServiceGetContentTypeSnapshot(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/content_types/hfM9RCJIk0wIm06WkEOQY/snapshots/4FLrUHftHW3v2BLi9fzfjU")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("snapshot-content_type_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	entrySnapshot, err := cma.Snapshots.GetContentTypeSnapshots(spaceID, "hfM9RCJIk0wIm06WkEOQY", "4FLrUHftHW3v2BLi9fzfjU")
	assertions.Nil(err)
	assertions.Equal("Blog Post", entrySnapshot.ContentTypeSnapshotDetail.Name)

}
