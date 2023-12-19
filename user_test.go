package contentful

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersService_Me(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/users/me")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("user_me.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	user, err := cmaClient.Users.Me()
	assertions.Nil(err)
	assertions.Equal("j.doe@labdigital.nl", user.Email)
}

func TestUsersService_Me_2(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/users/me")

		checkHeaders(r, assertions)

		w.WriteHeader(400)
		_, _ = fmt.Fprintln(w, readTestData("user_me.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	_, err = cmaClient.Users.Me()
	assertions.NotEmpty(err)
}
