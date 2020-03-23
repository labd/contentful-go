package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMembershipsServiceList(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/space_memberships")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("membership.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Memberships.List(spaceID).Next()
	assert.Nil(err)
}

func TestMembershipsServiceGet(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/space_memberships/0xWanD4AZI2AR35wW9q51n")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("membership_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Memberships.Get(spaceID, "0xWanD4AZI2AR35wW9q51n")
	assert.Nil(err)
}

func TestMembershipsServiceUpsertCreate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "POST")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/space_memberships")

		checkHeaders(r, assert)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)

		email := payload["email"].(string)
		admin := payload["admin"].(bool)
		assert.Equal("johndoe@nonexistent.com", email)
		assert.Equal(true, admin)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("membership_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	membership := &Membership{
		Admin: true,
		Roles: []Roles{
			Roles{
				Sys: &Sys{
					ID:       "1ElgCn1mi1UHSBLTP2v4TD",
					Type:     "Link",
					LinkType: "Role",
				},
			},
		},
		Email: "johndoe@nonexistent.com",
	}

	err = cma.Memberships.Upsert(spaceID, membership)
	assert.Nil(err)
}

func TestMembershipsServiceUpsertUpdate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/space_memberships/0xWanD4AZI2AR35wW9q51n")

		checkHeaders(r, assert)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)

		email := payload["email"].(string)
		assert.Equal("editedmail@examplemail.com", email)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("membership_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	membership, err := membershipFromTestData("membership_1.json")
	assert.Nil(err)

	membership.Email = "editedmail@examplemail.com"

	err = cma.Memberships.Upsert(spaceID, membership)
	assert.Nil(err)
}

func TestMembershipsServiceDelete(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/space_memberships/0xWanD4AZI2AR35wW9q51n")
		checkHeaders(r, assert)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test role
	membership, err := membershipFromTestData("membership_1.json")
	assert.Nil(err)

	// delete role
	err = cma.Memberships.Delete(spaceID, membership.Sys.ID)
	assert.Nil(err)
}
