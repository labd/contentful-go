package contentful

import (
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
