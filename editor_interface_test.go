package contentful

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEditorInterfacesService_List(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/content_types/hfM9RCJIk0wIm06WkEOQY/editor_interface")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("editor_interface.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.EditorInterfaces.List(spaceID, "hfM9RCJIk0wIm06WkEOQY").Next()
	assertions.Nil(err)

	interfaces := collection.ToEditorInterface()
	assertions.Equal(1, len(interfaces))
	assertions.Equal("name", interfaces[0].Controls[0].FieldID)
	assertions.Equal("extension", interfaces[0].SideBar[0].WidgetNameSpace)
}
