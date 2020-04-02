package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessTokensServiceList(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/users/me/access_tokens")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("access_token.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.AccessTokens.List().Next()
	assertions.Nil(err)
	keys := collection.ToAccessToken()
	assertions.Equal(2, len(keys))
	assertions.Equal("hioj6879UYGIfyt654tyfFHG", keys[0].Sys.ID)
}

func TestAccessTokensServiceGet(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/users/me/access_tokens/hioj6879UYGIfyt654tyfFHG")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("access_token_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	key, err := cma.AccessTokens.Get("hioj6879UYGIfyt654tyfFHG")
	assertions.Nil(err)
	assertions.Equal("hioj6879UYGIfyt654tyfFHG", key.Sys.ID)
}

func TestEntriesServiceCreate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/users/me/access_tokens")

		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)

		name := payload["name"].(string)
		revokedAt := payload["revokedAt"]

		assertions.Equal("Example Access Token", name)
		assertions.Equal(nil, revokedAt)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("access_token_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	accessToken := &AccessToken{
		Name:      "Example Access Token",
		RevokedAt: "",
		Scopes: []string{
			"content_management_manage",
		},
	}

	err = cma.AccessTokens.Create(accessToken)
	assertions.Nil(err)
}

func TestAccessTokensService_Revoke(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/users/me/access_tokens/hioj6879UYGIfyt654tyfFHG/revoked")

		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("2020-03-25T14:40:24Z", payload["revokedAt"])

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("access_token_updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	accessToken, err := accessTokenFromTestFile("access_token_updated.json")
	assertions.Nil(err)

	accessToken.RevokedAt = "2020-03-25T14:40:24Z"

	err = cma.AccessTokens.Revoke(accessToken)
	assertions.Nil(err)
	assertions.Equal(2, accessToken.Sys.Version)
}
