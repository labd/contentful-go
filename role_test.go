package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRolesServiceList(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/roles")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("role.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Roles.List(spaceID).Next()
	assert.Nil(err)
}

func TestRolesServiceGet(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/spaces/"+spaceID+"/roles/0xvkNW6WdQ8JkWlWZ8BC4x")

		checkHeaders(r, assert)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("role_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Roles.Get(spaceID, "0xvkNW6WdQ8JkWlWZ8BC4x")
	assert.Nil(err)
}

func TestRolesServiceUpsertCreate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "POST")
		assert.Equal(r.RequestURI, "/spaces/"+spaceID+"/roles")

		checkHeaders(r, assert)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)

		name := payload["name"]
		description := payload["description"]
		assert.Equal("Author", name)
		assert.Equal("Describes the author", description)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("role_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	role := &Role{
		Name:        "Author",
		Description: "Describes the author",
		Policies: []Policies{
			Policies{
				Effect: "allow",
				Actions: []string{
					"create",
				},
				Constraint: Constraint{
					And: []ConstraintDetail{
						{
							Equals: DetailItem{
								Doc: map[string]interface{}{
									"doc": "sys.type",
								},
								ItemType: "Entry",
							},
						},
					},
				},
			},
		},
		Permissions: Permissions{
			ContentModel: []string{
				"read",
			},
			Settings:           "all",
			ContentDelivery:    "all",
			Environments:       "all",
			EnvironmentAliases: "all",
		},
	}

	err = cma.Roles.Upsert(spaceID, role)
	assert.Nil(err)
}
