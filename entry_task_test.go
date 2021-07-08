package contentful

import (
	"encoding/json"
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
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries/5KsDBWseXY6QegucYAoacS/tasks")

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

	collection, err := cma.EntryTasks.List(env, "5KsDBWseXY6QegucYAoacS").Next()
	assertions.Nil(err)
	entryTasks := collection.ToEntryTask()
	assertions.Equal(1, len(entryTasks))
	assertions.Equal("Review translation", entryTasks[0].Body)
}

func TestEntryTasksService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries/5KsDBWseXY6QegucYAoacS/tasks/RHfHVRz3QkAgcMq4CGg2m5")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("entry_task_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	entryTask, err := cma.EntryTasks.Get(env, "5KsDBWseXY6QegucYAoacS", "RHfHVRz3QkAgcMq4CGg2m5")
	assertions.Nil(err)
	assertions.Equal("RHfHVRz3QkAgcMq4CGg2m5", entryTask.Sys.ID)
}

func TestEntryTasksService_Get_2(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries/5KsDBWseXY6QegucYAoacS/tasks/RHfHVRz3QkAgcMq4CGg2m5")

		checkHeaders(r, assertions)

		w.WriteHeader(400)
		_, _ = fmt.Fprintln(w, readTestData("entry_task_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.EntryTasks.Get(env, "5KsDBWseXY6QegucYAoacS", "RHfHVRz3QkAgcMq4CGg2m5")
	assertions.Nil(err)
}

func TestEntryTasksService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries/5KsDBWseXY6QegucYAoacS/tasks/RHfHVRz3QkAgcMq4CGg2m5")
		checkHeaders(r, assertions)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	entryTask, err := spaceFromTestData("entry_task_1.json")
	assertions.Nil(err)

	err = cma.EntryTasks.Delete(env, "5KsDBWseXY6QegucYAoacS", entryTask.Sys.ID)
	assertions.Nil(err)
}

func TestEntryTasksService_Upsert_Create(t *testing.T) {
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries/5KsDBWseXY6QegucYAoacS/tasks")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("new entry task", payload["body"])
		assertions.Equal("active", payload["status"])

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("entry_task_new.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	entryTask := &EntryTask{
		Body:   "new entry task",
		Status: "active",
		AssignedTo: AssignedTo{
			Sys: Sys{
				Type:     "Link",
				LinkType: "User",
				ID:       "7BslKh9TdKGOK41VmLDjFZ",
			},
		},
	}

	err := cma.EntryTasks.Upsert(env, "5KsDBWseXY6QegucYAoacS", entryTask)
	assertions.Nil(err)

	assertions.Equal("new entry task", entryTask.Body)
	assertions.Equal("active", entryTask.Status)
}

func TestEntryTasksService_Upsert_Update(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/"+environmentID+"/entries/5KsDBWseXY6QegucYAoacS/tasks/RHfHVRz3QkAgcMq4CGg2m5")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("Review translation", payload["body"])
		assertions.Equal("active", payload["status"])

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("entry_task_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	entryTask, err := entryTaskFromTestFile("entry_task_new.json")
	assertions.Nil(err)

	entryTask.Body = "Review translation"
	entryTask.Status = "active"

	err = cma.EntryTasks.Upsert(env, "5KsDBWseXY6QegucYAoacS", entryTask)
	assertions.Nil(err)
	assertions.Equal("Review translation", entryTask.Body)
	assertions.Equal("active", entryTask.Status)
}
