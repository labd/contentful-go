package testutil

import (
	"fmt"
	"github.com/flaconi/contentful-go"
	client2 "github.com/flaconi/contentful-go/pkgs/client"
	"github.com/flaconi/contentful-go/service/common"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	CMAToken = "b4c0n73n7fu1"
	SpaceID  = "id1"
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

func MockClient(
	t *testing.T,
	assertions *assert.Assertions,
	fixture ResponseData,
	callback HTTPHandler, validation ValidateRequest) (common.SpaceIdClientBuilder, *httptest.Server) {

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
		URL:       ts.URL,
		Debug:     false,
		UserAgent: "testclient",
		Token:     CMAToken,
	})

	if err != nil {
		t.Fatal(err)
	}

	return client, ts
}
