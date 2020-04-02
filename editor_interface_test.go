package contentful

import (
	"encoding/json"
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
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/editor_interface")

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

	collection, err := cma.EditorInterfaces.List(spaceID).Next()
	assertions.Nil(err)

	interfaces := collection.ToEditorInterface()
	assertions.Equal(1, len(interfaces))
	assertions.Equal("name", interfaces[0].Controls[0].FieldID)
	assertions.Equal("extension", interfaces[0].SideBar[0].WidgetNameSpace)
}

func TestEditorInterfacesService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/content_types/hfM9RCJIk0wIm06WkEOQY/editor_interface")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("editor_interface_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	editorInterface, err := cma.EditorInterfaces.Get(spaceID, "hfM9RCJIk0wIm06WkEOQY")
	assertions.Nil(err)
	assertions.Equal("name", editorInterface.Controls[0].FieldID)
	assertions.Equal("extension", editorInterface.SideBar[0].WidgetNameSpace)
}

func TestEditorInterfacesService_Get_2(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/master/content_types/hfM9RCJIk0wIm06WkEOQY/editor_interface")

		checkHeaders(r, assertions)

		w.WriteHeader(400)
		_, _ = fmt.Fprintln(w, readTestData("editor_interface_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.EditorInterfaces.Get(spaceID, "hfM9RCJIk0wIm06WkEOQY")
	assertions.Nil(err)
}

func TestEditorInterfacesService_Update(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/master/content_types/hfM9RCJIk0wIm06WkEOQY/editor_interface")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)

		assertions.Nil(err)
		controls := payload["controls"].([]interface{})
		sidebar := payload["sidebar"].([]interface{})
		assertions.Equal("changed id", controls[0].(map[string]interface{})["widgetId"].(string))
		assertions.Equal("someuiextension", sidebar[0].(map[string]interface{})["widgetId"].(string))

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("editor_interface_updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	editorInterface, err := editorInterfaceFromTestFile("editor_interface_1.json")
	assertions.Nil(err)

	editorInterface.Controls[0].WidgetID = "changed id"

	err = cma.EditorInterfaces.Update(spaceID, "hfM9RCJIk0wIm06WkEOQY", editorInterface)
	assertions.Nil(err)
	assertions.Equal("changed id", editorInterface.Controls[0].WidgetID)
}
