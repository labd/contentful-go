package contentful

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAppUploadService_Create(t *testing.T) {
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/organizations/organization_id/app_uploads")
		assertions.Equal("Bearer "+CMAToken, r.Header.Get("Authorization"))
		assertions.Equal("11187", r.Header.Get("Content-Length"))
		assertions.Equal("application/octet-stream", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("app_upload/upload.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.AppUpload.BaseURL = server.URL

	data, err := os.ReadFile("./testdata/resource_uploaded.png")

	assertions.Nil(err)

	uploadResult, err := cma.AppUpload.Create("organization_id", data)
	assertions.Nil(err)
	assertions.Equal("https://api.contentful.com", cma.BaseURL)
	assertions.Equal("<app_upload_id>", uploadResult.Sys.ID)
}
