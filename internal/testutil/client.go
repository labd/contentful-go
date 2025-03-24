package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labd/contentful-go"
	client2 "github.com/labd/contentful-go/pkgs/client"
	"github.com/labd/contentful-go/pkgs/util"
	"github.com/labd/contentful-go/service/cma"
	"github.com/stretchr/testify/assert"
)

var (
	CMAToken       = "b4c0n73n7fu1"
	SpaceID        = "id1"
	OrganizationId = "org1"
)

// HTTPHandler type defines callback from doing a mock HTTP request
type HTTPHandler func(w http.ResponseWriter, r *http.Request)

type ValidateRequest func(r *http.Request)

func checkHeaders(req *http.Request, assert *assert.Assertions) {
	assert.Equal("Bearer "+CMAToken, req.Header.Get("Authorization"))
	assert.Equal("application/vnd.contentful.management.v1+json", req.Header.Get("Content-Type"))
}

type ResponseData struct {
	Path       string
	StatusCode int
}

func MockCMAClient(
	t *testing.T,
	assertions *assert.Assertions,
	fixture ResponseData,
	callback HTTPHandler, validation ValidateRequest) (cma.SpaceIdClientBuilder, *httptest.Server) {

	handler := func(w http.ResponseWriter, r *http.Request) {

		validation(r)

		checkHeaders(r, assertions)

		if callback != nil {
			callback(w, r)
		} else {
			w.WriteHeader(fixture.StatusCode)
			if fixture.Path != "" {
				_, _ = fmt.Fprintln(w, readTestData(fixture.Path))
			}
		}

	}

	ts := httptest.NewServer(http.HandlerFunc(handler))

	client, err := contentful.NewCMAV2(client2.ClientConfig{
		URL:       util.ToPointer(ts.URL),
		Debug:     false,
		UserAgent: util.ToPointer("testclient"),
		Token:     CMAToken,
	})

	if err != nil {
		t.Fatal(err)
	}

	return client, ts
}
