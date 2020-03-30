package contentful

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppDefinitionsService_List(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/organizations/organization_id/app_definitions")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("app_definition.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.AppDefinitions.List("organization_id").Next()
	assertions.Nil(err)

	definitions := collection.ToAppDefinition()
	assertions.Equal(1, len(definitions))
	assertions.Equal("app_definition_id", definitions[0].Sys.ID)
	assertions.Equal("https://example.com/app.html", definitions[0].SRC)
}
