package contentful

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundError_Error(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = fmt.Fprintln(w, readTestData("error_notfound.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	// test space
	_, err = cmaClient.Spaces.Get("unknown-space-id")
	assertions.NotNil(err)
	_, ok := err.(NotFoundError)
	assertions.Equal(true, ok)
	notFoundError := err.(NotFoundError)
	assertions.Equal("the requested resource can not be found", notFoundError.Error())
	assertions.Equal(404, notFoundError.APIError.res.StatusCode)
	assertions.Equal("request-id", notFoundError.APIError.err.RequestID)
	assertions.Equal("The resource could not be found.", notFoundError.APIError.err.Message)
	assertions.Equal("Error", notFoundError.APIError.err.Sys.Type)
	assertions.Equal("NotFound", notFoundError.APIError.err.Sys.ID)
}

func TestRateLimitExceededError_Error(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		_, _ = fmt.Fprintln(w, readTestData("error_ratelimit.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	// test space
	space := &Space{Name: "test-space"}
	err = cmaClient.Spaces.Upsert(space)
	assertions.NotNil(err)
	_, ok := err.(RateLimitExceededError)
	assertions.Equal(true, ok)
	rateLimitExceededError := err.(RateLimitExceededError)
	assertions.Equal("You are creating too many Spaces.", rateLimitExceededError.Error())
	assertions.Equal(403, rateLimitExceededError.APIError.res.StatusCode)
	assertions.Equal("request-id", rateLimitExceededError.APIError.err.RequestID)
	assertions.Equal("You are creating too many Spaces.", rateLimitExceededError.APIError.err.Message)
	assertions.Equal("Error", rateLimitExceededError.APIError.err.Sys.Type)
	assertions.Equal("RateLimitExceeded", rateLimitExceededError.APIError.err.Sys.ID)
}

func TestAccessTokenInvalidError_Error(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		_, _ = fmt.Fprintln(w, readTestData("error_accesstoken_invalid.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	// test space
	space := &Space{Name: "test-space"}
	err = cmaClient.Spaces.Upsert(space)
	assertions.NotNil(err)
	_, ok := err.(AccessTokenInvalidError)
	assertions.Equal(true, ok)
	accessTokenInvalidError := err.(AccessTokenInvalidError)
	assertions.Equal("The access token you sent could not be found or is invalid.", accessTokenInvalidError.Error())
	assertions.Equal(401, accessTokenInvalidError.APIError.res.StatusCode)
	assertions.Equal("64adff93598dff78d8494c9d520990", accessTokenInvalidError.APIError.err.RequestID)
	assertions.Equal("The access token you sent could not be found or is invalid.", accessTokenInvalidError.APIError.err.Message)
	assertions.Equal("Error", accessTokenInvalidError.APIError.err.Sys.Type)
	assertions.Equal("AccessTokenInvalid", accessTokenInvalidError.APIError.err.Sys.ID)
}

func TestAccessToken422_Error(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(422)
		_, _ = fmt.Fprintln(w, readTestData("error_unique_422.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	// test space
	space := &Space{Name: "test-space"}
	err = cmaClient.Spaces.Upsert(space)
	assertions.NotNil(err)
	_, ok := err.(InvalidEntryError)
	assertions.Equal(true, ok)
	accessTokenInvalidError := err.(InvalidEntryError)
	assertions.Equal("Same field value present in other entry\n", accessTokenInvalidError.Error())
	assertions.Equal(422, accessTokenInvalidError.APIError.res.StatusCode)
	assertions.Equal("23e4333f-8fea-4f56-ac4d-4adeb8159185", accessTokenInvalidError.APIError.err.RequestID)
	assertions.Equal("Validation error", accessTokenInvalidError.APIError.err.Message)
	assertions.Equal("Error", accessTokenInvalidError.APIError.err.Sys.Type)
	assertions.Equal("InvalidEntry", accessTokenInvalidError.APIError.err.Sys.ID)
}
