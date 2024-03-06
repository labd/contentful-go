package contentful

import (
	"fmt"
	"github.com/flaconi/contentful-go/pkgs/common"
	"net/http"
	"net/http/httptest"
	"reflect"
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
	_, ok := err.(common.NotFoundError)
	assertions.Equal(true, ok)
	notFoundError := err.(common.NotFoundError)
	assertions.Equal("the requested resource can not be found", notFoundError.Error())

	refVal := reflect.ValueOf(&notFoundError.APIError).Elem()

	assertions.Equal(int64(404), refVal.FieldByName("res").Elem().FieldByName("StatusCode").Int())
	assertions.Equal("request-id", notFoundError.APIError.Err.RequestID)
	assertions.Equal("The resource could not be found.", notFoundError.APIError.Err.Message)
	assertions.Equal("Error", notFoundError.APIError.Err.Sys.Type)
	assertions.Equal("NotFound", notFoundError.APIError.Err.Sys.ID)
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
	_, ok := err.(common.RateLimitExceededError)
	assertions.Equal(true, ok)
	rateLimitExceededError := err.(common.RateLimitExceededError)
	assertions.Equal("You are creating too many Spaces.", rateLimitExceededError.Error())
	refVal := reflect.ValueOf(&rateLimitExceededError.APIError).Elem()

	assertions.Equal(int64(403), refVal.FieldByName("res").Elem().FieldByName("StatusCode").Int())
	assertions.Equal("request-id", rateLimitExceededError.APIError.Err.RequestID)
	assertions.Equal("You are creating too many Spaces.", rateLimitExceededError.APIError.Err.Message)
	assertions.Equal("Error", rateLimitExceededError.APIError.Err.Sys.Type)
	assertions.Equal("RateLimitExceeded", rateLimitExceededError.APIError.Err.Sys.ID)
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
	_, ok := err.(common.AccessTokenInvalidError)
	assertions.Equal(true, ok)
	accessTokenInvalidError := err.(common.AccessTokenInvalidError)
	assertions.Equal("The access token you sent could not be found or is invalid.", accessTokenInvalidError.Error())
	refVal := reflect.ValueOf(&accessTokenInvalidError.APIError).Elem()

	assertions.Equal(int64(401), refVal.FieldByName("res").Elem().FieldByName("StatusCode").Int())
	assertions.Equal("64adff93598dff78d8494c9d520990", accessTokenInvalidError.APIError.Err.RequestID)
	assertions.Equal("The access token you sent could not be found or is invalid.", accessTokenInvalidError.APIError.Err.Message)
	assertions.Equal("Error", accessTokenInvalidError.APIError.Err.Sys.Type)
	assertions.Equal("AccessTokenInvalid", accessTokenInvalidError.APIError.Err.Sys.ID)
}

func TestInvalidEntry422_Error(t *testing.T) {
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
	_, ok := err.(common.InvalidEntryError)
	assertions.Equal(true, ok)
	invalidEntryError := err.(common.InvalidEntryError)
	assertions.Equal("Same field value present in other entry\n", invalidEntryError.Error())
	refVal := reflect.ValueOf(&invalidEntryError.APIError).Elem()

	assertions.Equal(int64(422), refVal.FieldByName("res").Elem().FieldByName("StatusCode").Int())
	assertions.Equal("23e4333f-8fea-4f56-ac4d-4adeb8159185", invalidEntryError.APIError.Err.RequestID)
	assertions.Equal("Validation error", invalidEntryError.APIError.Err.Message)
	assertions.Equal("Error", invalidEntryError.APIError.Err.Sys.Type)
	assertions.Equal("InvalidEntry", invalidEntryError.APIError.Err.Sys.ID)
}
