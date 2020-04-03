package contentful

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResourcesService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/uploads/0xvkNW6WdQ8JkWlWZ8BC4x")

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("resource_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	urc = NewResourceClient(CMAToken)
	urc.BaseURL = server.URL

	resource, err := urc.Resources.Get(spaceID, "0xvkNW6WdQ8JkWlWZ8BC4x")
	assertions.Nil(err)
	assertions.Equal("2015-05-18T11:29:46.809Z", resource.Sys.CreatedAt)
	assertions.Equal("yadj1kx9rmg0", resource.Sys.Space.Sys.ID)
}

func TestResourcesService_Get_2(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/uploads/0xvkNW6WdQ8JkWlWZ8BC4x")

		w.WriteHeader(400)
		_, _ = fmt.Fprintln(w, readTestData("resource_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	urc = NewResourceClient(CMAToken)
	urc.BaseURL = server.URL

	_, err = urc.Resources.Get(spaceID, "0xvkNW6WdQ8JkWlWZ8BC4x")
	assertions.Nil(err)
}
