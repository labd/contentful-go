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
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/entries/hfM9RCJIk0wIm06WkEOQY/snapshots")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("snapshot-entry.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Snapshots.ListEntrySnapshots(spaceID, "hfM9RCJIk0wIm06WkEOQY").Next()
	assertions.Nil(err)
}

func TestEntriesServiceGetEntrySnapshot(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/entries/hfM9RCJIk0wIm06WkEOQY/snapshots/4FLrUHftHW3v2BLi9fzfjU")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("snapshot-entry_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Snapshots.GetEntrySnapshot(spaceID, "hfM9RCJIk0wIm06WkEOQY", "4FLrUHftHW3v2BLi9fzfjU")
	assertions.Nil(err)
}

func TestSnapshotsServiceListContentTypeSnapshots(t *testing.T) {
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

	_, err = cma.Snapshots.ListContentTypeSnapshots(spaceID, "hfM9RCJIk0wIm06WkEOQY").Next()
	assertions.Nil(err)
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

	_, err = cma.Snapshots.GetContentTypeSnapshots(spaceID, "hfM9RCJIk0wIm06WkEOQY", "4FLrUHftHW3v2BLi9fzfjU")
	assertions.Nil(err)
}
