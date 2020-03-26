package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScheduledActionsService_List(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/scheduled_actions?entity.sys.id=5KsDBWseXY6QegucYAoacS&environment.sys.id=master")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("scheduled_action.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.ScheduledActions.List(spaceID, "5KsDBWseXY6QegucYAoacS").Next()
	assertions.Nil(err)
	scheduledActions := collection.ToScheduledAction()
	assertions.Equal(1, len(scheduledActions))
	assertions.Equal("publish", scheduledActions[0].Action)
}

func TestScheduledActionsService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/scheduled_actions/3A13SXSDwO8c46NrjigFYT?entity.sys.id=5KsDBWseXY6QegucYAoacS&environment.sys.id=master")
		checkHeaders(r, assertions)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	scheduledAction, err := scheduledActionFromTestFile("scheduled_action_canceled.json")
	assertions.Nil(err)

	err = cma.ScheduledActions.Delete(spaceID, "5KsDBWseXY6QegucYAoacS", scheduledAction.Sys.ID)
	assertions.Nil(err)
	assertions.Equal("3A13SXSDwO8c46NrjigFYT", scheduledAction.Sys.ID)
}

func TestScheduledActionsService_Create(t *testing.T) {
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/scheduled_actions?entity.sys.id=5KsDBWseXY6QegucYAoacS&environment.sys.id=master")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("publish", payload["action"])

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, string(readTestData("scheduled_action_created.json")))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	scheduledAction := &ScheduledAction{
		Entity: Entity{
			Sys: Sys{
				Type:     "Link",
				LinkType: "Entry",
				ID:       "5KsDBWseXY6QegucYAoacS",
			},
		},
		Environment: EnvironmentLink{
			Sys: Sys{
				Type:     "Link",
				LinkType: "Environment",
				ID:       "master",
			},
		},
		ScheduledFor: map[string]string{
			"datetime": "2119-09-02T14:00:00.000Z",
		},
		Action: "publish",
	}

	err := cma.ScheduledActions.Create(spaceID, "5KsDBWseXY6QegucYAoacS", scheduledAction)
	assertions.Nil(err)

	assertions.Equal("publish", scheduledAction.Action)
}
