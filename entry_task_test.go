package contentful

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryTasksService_List(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/entries/5KsDBWseXY6QegucYAoacS/tasks")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("entry_task.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.EntryTasks.List(spaceID, "5KsDBWseXY6QegucYAoacS").Next()
	assertions.Nil(err)
	entryTasks := collection.ToEntryTask()
	assertions.Equal(1, len(entryTasks))
	assertions.Equal("Review translation", entryTasks[0].Body)
}
